package users

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/ngalaiko/miniboard/backend/httpx"
)

type logger interface {
	Error(string, ...interface{})
}

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
	type request struct {
		Username string `json:"username"`
		Password string `json:"password"`
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

	user, err := h.service.Create(r.Context(), req.Username, []byte(req.Password))
	switch {
	case err == nil:
		httpx.JSON(w, h.logger, user, http.StatusOK)
	case errors.Is(err, ErrAlreadyExists):
		httpx.Error(w, h.logger, err, http.StatusBadRequest)
	case errors.Is(err, ErrUsernameEmpty),
		errors.Is(err, ErrPasswordEmpty):
		httpx.Error(w, h.logger, errors.Unwrap(err), http.StatusBadRequest)
	default:
		h.logger.Error("failed to create user: %s", err)
		httpx.InternalError(w, h.logger)
	}
}
