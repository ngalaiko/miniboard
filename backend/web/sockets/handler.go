package sockets

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/websocket"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/httpx"
	"github.com/ngalaiko/miniboard/backend/items"
	"github.com/ngalaiko/miniboard/backend/subscriptions"
	"github.com/ngalaiko/miniboard/backend/tags"
)

type logger interface {
	Error(string, ...interface{})
}

type itemsService interface {
	Get(ctx context.Context, id string, userID string) (*items.UserItem, error)
	List(ctx context.Context, userID string, pageSize int, createdLT *time.Time, subscriptionID *string, tagID *string) ([]*items.UserItem, error)
}

type subscriptionsService interface {
	Create(ctx context.Context, userID string, url *url.URL, tagIDs []string) (*subscriptions.UserSubscription, error)
}

type tagsService interface {
	Create(ctx context.Context, userID string, title string) (*tags.Tag, error)
	GetByTitle(ctx context.Context, userID string, title string) (*tags.Tag, error)
}

type Handler struct {
	logger               logger
	itemsService         itemsService
	subscriptionsService subscriptionsService
	tagsService          tagsService

	openSocketsGuard *sync.RWMutex
	openSockets      map[string][]*websocket.Conn
}

func NewHandler(
	logger logger,
	itemsService itemsService,
	subscriptionsService subscriptionsService,
	tagsService tagsService,
) *Handler {
	return &Handler{
		logger:               logger,
		itemsService:         itemsService,
		subscriptionsService: subscriptionsService,
		tagsService:          tagsService,
		openSocketsGuard:     &sync.RWMutex{},
		openSockets:          make(map[string][]*websocket.Conn),
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, auth := authorizations.FromContext(r.Context())
	if !auth {
		httpx.Error(w, h.logger, fmt.Errorf(""), http.StatusUnauthorized)
	}
	websocket.Handler(h.handle(r.Context(), token.UserID)).ServeHTTP(w, r)
}

func (h *Handler) readRequests(ws *websocket.Conn) <-chan *request {
	requests := make(chan *request)
	go func() {
		for {
			req := &request{}
			err := websocket.JSON.Receive(ws, req)
			switch {
			case err == nil:
				requests <- req
			case errors.Is(err, io.EOF):
				if err := ws.Close(); err != nil {
					h.logger.Error("failed to close ws connection: %s", err)
				}
				close(requests)
				return
			default:
				h.logger.Error("failed to read a request: %s", err)
			}
		}
	}()
	return requests
}

func (h *Handler) broadcastResponses(ctx context.Context, userID string, responses <-chan *response) {
	for {
		select {
		case <-ctx.Done():
			return
		case resp := <-responses:
			h.openSocketsGuard.RLock()
			for _, ws := range h.openSockets[userID] {
				if err := websocket.JSON.Send(ws, resp); err != nil {
					h.logger.Error("failed to write response message: %s", err)
				}
			}
			h.openSocketsGuard.RUnlock()
		}
	}
}

func (h *Handler) openSocket(userID string, ws *websocket.Conn) {
	h.openSocketsGuard.Lock()
	h.openSockets[userID] = append(h.openSockets[userID], ws)
	h.openSocketsGuard.Unlock()
}

func (h *Handler) closeSocket(userID string, ws *websocket.Conn) {
	h.openSocketsGuard.Lock()
	i := 0
	for _, s := range h.openSockets[userID] {
		if s == ws {
			break
		}
		i++
	}
	h.openSockets[userID] = append(h.openSockets[userID][:i], h.openSockets[userID][i+1:]...)
	h.openSocketsGuard.Unlock()
}

func (h *Handler) handle(ctx context.Context, userID string) func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		h.openSocket(userID, ws)
		defer h.closeSocket(userID, ws)

		responses := make(chan *response)
		go h.broadcastResponses(ctx, userID, responses)

		for req := range h.readRequests(ws) {
			switch req.Event {
			case subscriptionsImport:
				h.onSubscriptionsImport(context.Background(), userID, req, responses)
			case subscriptionsCreated:
				responses <- h.onSubscriptionsCreated(context.Background(), userID, req)
			case itemsSelect:
				h.respond(ws, h.onItemSelected(ctx, userID, req))
			case itemsLoad:
				rr := h.loadItems(ctx, userID, req)
				h.respond(ws, resetResponse(req, "#items-list"))
				for _, r := range rr {
					h.respond(ws, r)
				}
			case itemsLoadmore:
				rr := h.loadItems(ctx, userID, req)
				for _, r := range rr {
					h.respond(ws, r)
				}
			default:
				h.respond(ws, errResponse(req, fmt.Errorf("unknown event: '%s'", req.Event)))
			}
		}
	}
}

func (h *Handler) respond(ws *websocket.Conn, resp *response) {
	if err := websocket.JSON.Send(ws, resp); err != nil {
		h.logger.Error("failed to write response message: %s", err)
	}
}
