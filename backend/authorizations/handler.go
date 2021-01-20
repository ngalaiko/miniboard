package authorizations

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/ngalaiko/miniboard/backend/httpx"
	"github.com/ngalaiko/miniboard/backend/users"
)

// Config contains authorizations handler config.
type Config struct {
	Domain *string `yaml:"domain"`
	Secure bool    `yaml:"secure"`
}

type usersService interface {
	GetByUsername(context.Context, string) (*users.User, error)
}

type jwtService interface {
	NewToken(context.Context, string) (*Token, error)
}

// Handler handlers authorization http requests.
type Handler struct {
	logger       logger
	usersService usersService
	jwtService   jwtService
	config       *Config
}

// NewHandler initializes a new handler.
func NewHandler(usersService usersService, jwtService jwtService, logger logger, cfg *Config) *Handler {
	if cfg == nil {
		cfg = &Config{}
	}
	return &Handler{
		logger:       logger,
		usersService: usersService,
		jwtService:   jwtService,
		config:       cfg,
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
		h.handleCreateAuthorization(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) handleCreateAuthorization(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.usersService.GetByUsername(r.Context(), req.Username)
	switch {
	case err == nil:
	case errors.Is(users.ErrNotFound, err):
		httpx.Error(w, h.logger, users.ErrNotFound, http.StatusBadRequest)
		return
	default:
		h.logger.Error("failed to get user: %s", err)
		httpx.InternalError(w, h.logger)
		return
	}

	validatePasswordErr := user.ValidatePassword([]byte(req.Password))
	switch {
	case validatePasswordErr == nil:
	case errors.Is(validatePasswordErr, users.ErrInvalidPassword):
		httpx.Error(w, h.logger, users.ErrInvalidPassword, http.StatusBadRequest)
		return
	default:
		h.logger.Error("failed to validate password: %s", err)
		httpx.InternalError(w, h.logger)
		return
	}

	token, err := h.jwtService.NewToken(r.Context(), user.ID)
	if err != nil {
		h.logger.Error("failed to create a new token: %s", err)
		httpx.InternalError(w, h.logger)
		return
	}

	setCookie(w, h.config, token)

	httpx.JSON(w, h.logger, token, http.StatusOK)
}
