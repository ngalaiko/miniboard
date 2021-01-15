package feeds

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/ngalaiko/miniboard/server/authorizations"
	"github.com/ngalaiko/miniboard/server/httpx"
	"github.com/ngalaiko/miniboard/server/operations"
)

// Known errors.
var (
	errInvalidURL = fmt.Errorf("got invalid url")
	errEmptyURL   = fmt.Errorf("got empty url")
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

func (h *Handler) handleCreateFeed(w http.ResponseWriter, r *http.Request) {
	token, auth := authorizations.FromContext(r.Context())
	if !auth {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type request struct {
		URL string `json:"url"`
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

	operation, err := h.operationService.Create(r.Context(), token.UserID, h.createFeed(token.UserID, url))
	switch {
	case err == nil:
		httpx.JSON(w, h.logger, operation)
	default:
		h.logger.Error("failed to create feed: %s", err)
		httpx.InternalError(w, h.logger)
	}
}

func (h *Handler) createFeed(userID string, url *url.URL) operations.Task {
	return func(ctx context.Context, operation *operations.Operation, status chan<- *operations.Operation) error {
		feed, err := h.service.Create(ctx, userID, url)
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
			return err
		}
	}
}
