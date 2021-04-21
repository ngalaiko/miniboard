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

func (h *Handler) readRequests(c *websocket.Conn) <-chan *request {
	requests := make(chan *request)
	go func() {
		buffer := make([]byte, 2048)
		for {
			n, err := c.Read(buffer)
			switch err {
			case nil:
			case io.EOF:
				if err := c.Close(); err != nil {
					h.logger.Error("failed to close ws connection: %s", err)
				}
				close(requests)
				return
			default:
				h.logger.Error("failed to read a frame: %s", err)
				h.onResponse(c, errResponse(&request{}, errInternal))
				continue
			}

			req := &request{}
			if err := json.Unmarshal(buffer[:n], req); err != nil {
				h.logger.Error("failed to unmarshal request: %s", err)
				h.onResponse(c, errResponse(&request{}, errInternal))
				continue
			}
			requests <- req
		}
	}()
	return requests
}

func (h *Handler) sendResponses(c *websocket.Conn, responses <-chan *response) {
	for resp := range responses {
		raw, err := json.Marshal(resp)
		if err != nil {
			h.logger.Error("failed to marshal response message: %s", err)
			h.onResponse(c, errResponse(&request{ID: resp.ID}, errInternal))
			return
		}
		if _, err := c.Write(raw); err != nil {
			h.logger.Error("failed to write response message: %s", err)
			return
		}
	}
}

func (h *Handler) handle(ctx context.Context, userID string) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		responses := make(chan *response)
		go h.sendResponses(c, responses)
		for req := range h.readRequests(c) {
			switch req.Event {
			case itemsSelect:
				responses <- h.onItemSelected(ctx, userID, req)
			case itemsLoad:
				rr := h.loadItems(ctx, userID, req)
				responses <- resetResponse(req, "#items-list")
				for _, r := range rr {
					responses <- r
				}
			case itemsLoadmore:
				rr := h.loadItems(ctx, userID, req)
				for _, r := range rr {
					responses <- r
				}
			default:
				h.onResponse(c, errResponse(req, fmt.Errorf("unknown event: '%s'", req.Event)))
			}
		}
	}
}

func (h *Handler) onResponse(c *websocket.Conn, resp *response) {
	data, err := json.Marshal(resp)
	if err != nil {
		h.logger.Error("failed to marshal response message: %s", err)
		h.onResponse(c, errResponse(&request{ID: resp.ID}, errInternal))
		return
	}
	if _, err := c.Write(data); err != nil {
		h.logger.Error("failed to write response message: %s", err)
		return
	}
}
