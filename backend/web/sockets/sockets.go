package sockets

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/httpx"
	"golang.org/x/net/websocket"
)

type Respond func(*Response)

type Broadcast func(*Response)

type Handler func(context.Context, *Request, Respond, Broadcast)

type logger interface {
	Error(string, ...interface{})
}

type Sockets struct {
	logger logger

	openSocketsGuard *sync.RWMutex
	openSockets      map[string][]*websocket.Conn

	handlers map[string]Handler
}

func New(logger logger) *Sockets {
	return &Sockets{
		logger:           logger,
		openSocketsGuard: &sync.RWMutex{},
		openSockets:      make(map[string][]*websocket.Conn),
		handlers:         map[string]Handler{},
	}
}

func (s *Sockets) On(event string, handler Handler) *Sockets {
	s.handlers[event] = handler
	return s
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
	i := 0
	for _, w := range s.openSockets[userID] {
		if *w == *ws {
			break
		}
		i++
	}
	s.openSockets[userID] = append(s.openSockets[userID][:i], s.openSockets[userID][i+1:]...)
	s.openSocketsGuard.Unlock()
}

func (s *Sockets) handle(ctx context.Context, userID string) func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		s.register(userID, ws)

		for {
			req := &Request{}
			err := websocket.JSON.Receive(ws, req)
			switch {
			case err == nil:
				if handler, found := s.handlers[req.Event]; found {
					handler(ctx, req, s.respond(ws), s.broadcast(userID))
				} else {
					s.respond(ws)(&Response{
						ID:    req.ID,
						Error: fmt.Sprintf("unknown event '%s'", req.Event),
					})
				}
			case errors.Is(err, io.EOF):
				s.unregister(userID, ws)
			default:
				s.logger.Error("failed to read a request: %s", err)
			}
		}
	}
}

func (s *Sockets) respond(ws *websocket.Conn) func(*Response) {
	return func(resp *Response) {
		if err := websocket.JSON.Send(ws, resp); err != nil {
			s.logger.Error("failed to write response message: %s", err)
		}
	}
}

func (s *Sockets) broadcast(userID string) func(*Response) {
	return func(resp *Response) {
		s.openSocketsGuard.RLock()
		for _, ws := range s.openSockets[userID] {
			s.respond(ws)(resp)
		}
		s.openSocketsGuard.RUnlock()
	}
}
