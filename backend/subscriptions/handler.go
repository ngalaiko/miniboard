package subscriptions

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
	Debug(string, ...interface{})
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

// List returns handler func that handles subscription list via http.
func (h *Handler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		subscriptions, err := h.service.List(r.Context(), token.UserID, pageSize, createdLT)
		switch {
		case err == nil:
			type response struct {
				Subscriptions []*UserSubscription `json:"subscriptions"`
			}
			httpx.JSON(w, h.logger, &response{Subscriptions: subscriptions}, http.StatusOK)
		default:
			h.logger.Error("failed to list subscriptions: %s", err)
			httpx.InternalError(w, h.logger)
		}
	}
}

// Import returns handler func that that handles subscriptions import from opml.
func (h *Handler) Import() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnsupportedMediaType)
	}
}

// Create returns handler func that handles subscription creation via http.
func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		tagIDs := []string{}
		if req.TagIDs != nil {
			tagIDs = req.TagIDs
		}

		operation, err := h.operationService.Create(r.Context(), token.UserID, h.createSubscription(token.UserID, url, tagIDs))
		switch {
		case err == nil:
			httpx.JSON(w, h.logger, operation, http.StatusOK)
		default:
			h.logger.Error("failed to create subscription: %s", err)
			httpx.InternalError(w, h.logger)
		}
	}
}

func (h *Handler) createSubscription(userID string, url *url.URL, tagIDs []string) operations.Task {
	return func(ctx context.Context, operation *operations.Operation, status chan<- *operations.Operation) error {
		subscription, err := h.service.Create(ctx, userID, url, tagIDs)
		switch {
		case err == nil:
			operation.Success(subscription)
			status <- operation
			return nil
		case errors.Is(err, errFailedToDownloadSubscription),
			errors.Is(err, errAlreadyExists),
			errors.Is(err, errFailedToParseSubscription):
			operation.Error(err.Error())
			status <- operation
			return nil
		default:
			h.logger.Error("failed to create subscription: %s", err)
			return fmt.Errorf("internal error")
		}
	}
}
