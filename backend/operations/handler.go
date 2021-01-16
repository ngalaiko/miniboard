package operations

import (
	"errors"
	"net/http"
	"regexp"

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

// ServeHTTP implements http.Handler.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGet().ServeHTTP(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleGet() http.Handler {
	getOperation := regexp.MustCompile(`/(.*)$`)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := getOperation.FindStringSubmatch(r.URL.Path)
		if len(m) == 0 {
			http.NotFound(w, r)
			return
		}
		h.handleGetOperation(w, r, m[1])
	})
}

func (h *Handler) handleGetOperation(w http.ResponseWriter, r *http.Request, id string) {
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
