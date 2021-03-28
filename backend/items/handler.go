package items

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/httpx"
)

// Known errors.
var (
	errInvalidPageSize  = fmt.Errorf("failed to parse page_size")
	errInvalidCreatedLT = fmt.Errorf("failed to parse created_lt param")
)

// Handler handles http requests for user resource.
type Handler struct {
	service *Service
	logger  logger
}

// NewHandler creates a new handler for users resource.
func NewHandler(service *Service, logger logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Get returns item by id via http.
func (h *Handler) Get(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, auth := authorizations.FromContext(r.Context())
		if !auth {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		item, err := h.service.Get(r.Context(), id, token.UserID)
		switch {
		case err == nil:
			httpx.JSON(w, h.logger, item, http.StatusOK)
		case errors.Is(err, errNotFound):
			http.NotFound(w, r)
		default:
			h.logger.Error("failed to get operation: %s", err)
			httpx.InternalError(w, h.logger)
		}
	}
}

// List returns handler func that handles items list via http.
func (h *Handler) List(tagID *string, subscriptionID *string) http.HandlerFunc {
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

		items, err := h.service.List(r.Context(), token.UserID, pageSize, createdLT, subscriptionID, tagID)
		switch {
		case err == nil:
			type response struct {
				Items []*UserItem `json:"items"`
			}
			httpx.JSON(w, h.logger, &response{Items: items}, http.StatusOK)
		default:
			h.logger.Error("failed to list items: %s", err)
			httpx.InternalError(w, h.logger)
		}
	}
}
