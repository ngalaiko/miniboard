package tags

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/httpx"
)

// Known errors.
var (
	errEmptyTitle       = fmt.Errorf("got empty title")
	errInvalidPageSize  = fmt.Errorf("failed to parse page_size")
	errInvalidCreatedLT = fmt.Errorf("failed to parse created_lt param")
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

// List returns http handler that lists tag.
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

		tags, err := h.service.List(r.Context(), token.UserID, pageSize, createdLT)
		switch {
		case err == nil:
			type response struct {
				Tags []*Tag `json:"tags"`
			}
			httpx.JSON(w, h.logger, &response{Tags: tags}, http.StatusOK)
		default:
			h.logger.Error("failed to list tags: %s", err)
			httpx.InternalError(w, h.logger)
		}
	}
}

// Create returns http handler that creates tag.
func (h *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, auth := authorizations.FromContext(r.Context())
		if !auth {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type request struct {
			Title string `json:"title"`
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

		if req.Title == "" {
			httpx.Error(w, h.logger, errEmptyTitle, http.StatusBadRequest)
			return
		}

		tag, err := h.service.Create(r.Context(), token.UserID, req.Title)
		switch {
		case errors.Is(err, errAlreadyExists):
			httpx.Error(w, h.logger, err, http.StatusBadRequest)
		case err == nil:
			httpx.JSON(w, h.logger, tag, http.StatusOK)
		default:
			h.logger.Error("failed to create tag: %s", err)
			httpx.InternalError(w, h.logger)
		}
	}
}
