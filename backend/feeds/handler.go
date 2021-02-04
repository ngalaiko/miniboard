package feeds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/httpx"
	"github.com/ngalaiko/miniboard/backend/operations"
)

// Known errors.
var (
	errInvalidURL       = fmt.Errorf("got invalid url")
	errEmptyURL         = fmt.Errorf("got empty url")
	errInvalidPageSize  = fmt.Errorf("failed to parse page_size")
	errInvalidCreatedLT = fmt.Errorf("failed to parse created_lt param")
)

type logger interface {
	Error(string, ...interface{})
}

type operationService interface {
	Create(context.Context, string, operations.Task) (*operations.Operation, error)
}

// Handler handles http requests for user resource.
type Handler struct {
	service          *Service
	logger           logger
	operationService operationService
}

// NewHandler creates a new handler for users resource.
func NewHandler(service *Service, logger logger, operationService operationService) *Handler {
	return &Handler{
		service:          service,
		logger:           logger,
		operationService: operationService,
	}
}

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodGet:
		h.handleGet(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		h.handleCreateFeed(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) handleGet(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		h.handleListFeeds(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) handleListFeeds(w http.ResponseWriter, r *http.Request) {
	token, auth := authorizations.FromContext(r.Context())
	if !auth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var pageSize int = 100
	if pageSizeRaw := r.URL.Query().Get("page_size"); len(pageSizeRaw) != 0 {
		pageSizeParsed, err := strconv.Atoi(pageSizeRaw)
		if err != nil {
			httpx.Error(w, h.logger, errInvalidPageSize, http.StatusBadRequest)
			return
		}
		pageSize = pageSizeParsed
	}

	var createdLT *time.Time
	if createdLTParam := r.URL.Query().Get("created_lt"); len(createdLTParam) != 0 {
		createdLTParsed, err := time.Parse(time.RFC3339, createdLTParam)
		if err != nil {
			httpx.Error(w, h.logger, errInvalidCreatedLT, http.StatusBadRequest)
			return
		}
		createdLT = &createdLTParsed
	}

	var tagID *string
	if tagIDEqParam, ok := r.URL.Query()["tag_id_eq"]; ok {
		if len(tagIDEqParam[0]) > 0 {
			tagID = &tagIDEqParam[0]
		} else {
			tagID = new(string)
			*tagID = ""
		}
	}

	feeds, err := h.service.List(r.Context(), token.UserID, pageSize, createdLT, tagID)
	switch {
	case err == nil:
		type response struct {
			Feeds []*Feed `json:"feeds"`
		}
		httpx.JSON(w, h.logger, &response{Feeds: feeds}, http.StatusOK)
	default:
		h.logger.Error("failed to list feeds: %s", err)
		httpx.InternalError(w, h.logger)
	}
}

func (h *Handler) handleCreateFeed(w http.ResponseWriter, r *http.Request) {
	token, auth := authorizations.FromContext(r.Context())
	if !auth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type request struct {
		URL    string   `json:"url"`
		TagIDs []string `json:"tag_ids"`
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed to read request body: %s", err)
		httpx.InternalError(w, h.logger)
		return
	}

	req := &request{}
	if len(body) > 0 {
		if err := json.Unmarshal(body, req); err != nil {
			h.logger.Error("failed unmarshal request: %s", err)
			httpx.InternalError(w, h.logger)
			return
		}
	}

	if req.URL == "" {
		httpx.Error(w, h.logger, errEmptyURL, http.StatusBadRequest)
		return
	}

	url, err := url.ParseRequestURI(req.URL)
	if err != nil {
		httpx.Error(w, h.logger, errInvalidURL, http.StatusBadRequest)
		return
	}

	operation, err := h.operationService.Create(r.Context(), token.UserID, h.createFeed(token.UserID, url, req.TagIDs))
	switch {
	case err == nil:
		httpx.JSON(w, h.logger, operation, http.StatusOK)
	default:
		h.logger.Error("failed to create feed: %s", err)
		httpx.InternalError(w, h.logger)
	}
}

func (h *Handler) createFeed(userID string, url *url.URL, tagIDs []string) operations.Task {
	return func(ctx context.Context, operation *operations.Operation, status chan<- *operations.Operation) error {
		feed, err := h.service.Create(ctx, userID, url, tagIDs)
		switch {
		case err == nil:
			operation.Success(feed)
			status <- operation
			return nil
		case errors.Is(err, errFailedToDownloadFeed),
			errors.Is(err, errAlreadyExists),
			errors.Is(err, errFailedToParseFeed):
			operation.Error(err.Error())
			status <- operation
			return nil
		default:
			h.logger.Error("failed to create feed: %s", err)
			return fmt.Errorf("internal error")
		}
	}
}
