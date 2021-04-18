package sockets

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/httpx"
	"github.com/ngalaiko/miniboard/backend/items"
)

type logger interface {
	Error(string, ...interface{})
}

type itemsService interface {
	Get(ctx context.Context, id string, userID string) (*items.UserItem, error)
	List(ctx context.Context, userID string, pageSize int, createdLT *time.Time, subscriptionID *string, tagID *string) ([]*items.UserItem, error)
}

type Handler struct {
	logger       logger
	itemsService itemsService
}

func NewHandler(logger logger, itemsService itemsService) *Handler {
	return &Handler{
		logger:       logger,
		itemsService: itemsService,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, auth := authorizations.FromContext(r.Context())
	if !auth {
		httpx.Error(w, h.logger, fmt.Errorf(""), http.StatusUnauthorized)
	}
	websocket.Handler(h.handle(r.Context(), token.UserID)).ServeHTTP(w, r)
}

func (h *Handler) handle(ctx context.Context, userID string) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		buffer := make([]byte, 2048)
		for {
			n, err := c.Read(buffer)
			switch err {
			case nil:
			case io.EOF:
				if err := c.Close(); err != nil {
					h.logger.Error("failed to close ws connection: %s", err)
				}
				return
			default:
				h.logger.Error("failed to read a frame: %s", err)
				h.onError(c, 0, errInternal)
				continue
			}

			req := &request{}
			if err := json.Unmarshal(buffer[:n], req); err != nil {
				h.logger.Error("failed to unmarshal request: %s", err)
				h.onError(c, 0, errInternal)
				continue
			}

			switch req.Event {
			case tagToggled:
			case itemSelected:
				response, err := h.onItemSelected(ctx, userID, req)
				if err != nil {
					h.onError(c, req.ID, err)
					continue
				}
				h.onResponse(c, response)
			case itemsLoad:
				response, err := h.loadItems(ctx, userID, req)
				if err != nil {
					h.onError(c, req.ID, err)
					continue
				}
				response.Reset = true
				h.onResponse(c, response)
			case itemsLoadmore:
				response, err := h.loadItems(ctx, userID, req)
				if err != nil {
					h.onError(c, req.ID, err)
					continue
				}
				h.onResponse(c, response)
			default:
				h.onError(c, req.ID, fmt.Errorf("unknown event: '%s'", req.Event))
			}
		}
	}
}

func (h *Handler) onResponse(c *websocket.Conn, response *response) {
	data, err := json.Marshal(response)
	if err != nil {
		h.logger.Error("failed to marshal response message: %s", err)
		h.onError(c, response.ID, errInternal)
		return
	}
	if _, err := c.Write(data); err != nil {
		h.logger.Error("failed to write response message: %s", err)
		return
	}
}

func (h *Handler) onError(c *websocket.Conn, id uint, err error) {
	raw, err := json.Marshal(&response{
		ID:    id,
		Error: err.Error(),
	})
	if err != nil {
		h.logger.Error("failed to marshal error message: %s", err)
		return
	}
	if _, err := c.Write(raw); err != nil {
		h.logger.Error("failed to write error message: %s", err)
		return
	}
}
