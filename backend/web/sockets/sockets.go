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

type itemsService interface {
	Get(ctx context.Context, id string, userID string) (*items.UserItem, error)
	List(ctx context.Context, userID string, pageSize int, createdLT *time.Time, subscriptionID *string, tagID *string) ([]*items.UserItem, error)
}

type subscriptionsService interface {
	Create(ctx context.Context, userID string, url *url.URL, tagIDs []string) (*subscriptions.UserSubscription, error)
	List(ctx context.Context, userID string, pageSize int, createdLT *time.Time) ([]*subscriptions.UserSubscription, error)
}

type tagsService interface {
	Create(ctx context.Context, userID string, title string) (*tags.Tag, error)
	GetByTitle(ctx context.Context, userID string, title string) (*tags.Tag, error)
	List(ctx context.Context, userID string, pageSize int, createdLT *time.Time) ([]*tags.Tag, error)
}

type Request struct {
	ID     uint              `json:"id"`
	Event  string            `json:"event"`
	Params map[string]string `json:"params"`
}

type Respond func(*Response)

type Broadcast func(*Response)

type Handler func(context.Context, *Request, Respond, Broadcast)

type logger interface {
	Error(string, ...interface{})
}

type Sockets struct {
	logger logger

	itemsService         itemsService
	tagsService          tagsService
	subscriptionsService subscriptionsService

	openSocketsGuard *sync.RWMutex
	openSockets      map[string][]*websocket.Conn
}

func New(
	logger logger,
	itemsService itemsService,
	tagsService tagsService,
	subscriptionsService subscriptionsService,
) *Sockets {
	return &Sockets{
		logger:               logger,
		itemsService:         itemsService,
		tagsService:          tagsService,
		subscriptionsService: subscriptionsService,
		openSocketsGuard:     &sync.RWMutex{},
		openSockets:          make(map[string][]*websocket.Conn),
	}
}

func (s *Sockets) Receive() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, auth := authorizations.FromContext(r.Context())
		if !auth {
			httpx.Error(w, s.logger, fmt.Errorf(""), http.StatusUnauthorized)
		} else {
			websocket.Handler(s.handle(r.Context(), token.UserID)).ServeHTTP(w, r)
		}
	}
}

func (s *Sockets) register(userID string, ws *websocket.Conn) {
	s.openSocketsGuard.Lock()
	s.openSockets[userID] = append(s.openSockets[userID], ws)
	s.openSocketsGuard.Unlock()
}

func (s *Sockets) unregister(userID string, ws *websocket.Conn) {
	s.openSocketsGuard.Lock()
	defer s.openSocketsGuard.Unlock()

	for i, w := range s.openSockets[userID] {
		if *w != *ws {
			continue
		}
		s.openSockets[userID] = append(s.openSockets[userID][:i], s.openSockets[userID][i+1:]...)
		return
	}
}

func (s *Sockets) handle(ctx context.Context, userID string) func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		s.register(userID, ws)
		defer s.unregister(userID, ws)

		for {
			req := &Request{}
			err := websocket.JSON.Receive(ws, req)
			switch {
			case err == nil:
				switch req.Event {
				case "items:select":
					s.itemsSelect(ws, userID, req)
				case "items:loadmore":
					s.itemsLoadmore(ws, userID, req)
				case "items:load":
					s.itemsLoad(ws, userID, req)
				case "subscriptions:create":
					s.subscriptionsCreate(ws, userID, req)
				case "subscriptions:import":
					s.subscriptionsImport(ws, userID, req)
				default:
					s.respond(ws, &Response{
						ID:    req.ID,
						Error: fmt.Sprintf("unknown event '%s'", req.Event),
					})
				}
			case errors.Is(err, io.EOF):
				return
			default:
				s.logger.Error("failed to read a request: %s", err)
			}
		}
	}
}

func (s *Sockets) respond(ws *websocket.Conn, resp *Response) {
	if err := websocket.JSON.Send(ws, resp); err != nil {
		s.logger.Error("failed to write response message: %s", err)
	}
}

func (s *Sockets) broadcast(userID string, resp *Response) {
	s.openSocketsGuard.RLock()
	for _, ws := range s.openSockets[userID] {
		s.respond(ws, resp)
	}
	s.openSocketsGuard.RUnlock()
}
