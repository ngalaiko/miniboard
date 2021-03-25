package items

import (
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

// List returns handler func that handles items list via http.
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

		var subscriptionID *string
		if subscriptionIDParam := r.URL.Query().Get("subscription_id_eq"); len(subscriptionIDParam) != 0 {
			subscriptionID = &subscriptionIDParam
		}

		var tagID *string
		if tagIDParam := r.URL.Query().Get("tag_id_eq"); len(tagIDParam) != 0 {
			tagID = &tagIDParam
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
