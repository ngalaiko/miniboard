package users

import (
	"net/http"
)

// Handler handles http requests for user resource.
type Handler struct {
	service *Service
}

// NewHandler creates a new handler for users resource.
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
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
		h.handleCreateUser(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	// todo
}
