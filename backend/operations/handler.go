package operations

import (
	"errors"
	"net/http"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/httpx"
)

// Handler handles http requests for user resource.
type Handler struct {
	service *Service
	logger  logger
}

// NewHandler creates a new handler for operations resource.
func NewHandler(service *Service, logger logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

// Get returns operations by id via http.
func (h *Handler) Get(id string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, auth := authorizations.FromContext(r.Context())
		if !auth {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		operation, err := h.service.Get(r.Context(), id, token.UserID)
		switch {
		case err == nil:
			httpx.JSON(w, h.logger, operation, http.StatusOK)
		case errors.Is(err, errNotFound):
			http.NotFound(w, r)
		default:
			h.logger.Error("failed to get operation: %s", err)
			httpx.InternalError(w, h.logger)
		}
	}
}
