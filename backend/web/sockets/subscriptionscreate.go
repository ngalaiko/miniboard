package sockets

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"

	"golang.org/x/net/websocket"

	"github.com/ngalaiko/miniboard/backend/subscriptions"
	"github.com/ngalaiko/miniboard/backend/web/render"
)

func (s *Sockets) subscriptionsCreate(ws *websocket.Conn, userID string, req *request) {
	ctx := ws.Request().Context()

	urlParam, ok := req.Params["url"]
	if !ok {
		s.respond(ws, errResponse(req, fmt.Errorf("'url' parameter is missing")))
		return
	}

	url, err := url.ParseRequestURI(urlParam)
	if err != nil {
		s.respond(ws, errResponse(req, fmt.Errorf("'url' is invalid: %s", err)))
		return
	}

	subscription, err := s.subscriptionsService.Create(ctx, userID, url, nil)
	switch {
	case err == nil:
	case errors.Is(err, subscriptions.ErrFailedToDownloadSubscription),
		errors.Is(err, subscriptions.ErrAlreadyExists),
		errors.Is(err, subscriptions.ErrFailedToParseSubscription):
		s.respond(ws, errResponse(req, err))
		return
	default:
		s.logger.Error("failed to create subscription: %s", err)
		s.respond(ws, errResponse(req, errInternal))
		return
	}

	html := &bytes.Buffer{}
	if err := render.Subscription(html, subscription); err != nil {
		s.logger.Error("failed to render subscription: %s", err)
		s.respond(ws, errResponse(req, errInternal))
		return
	}

	s.broadcast(userID, &response{
		ID:     req.ID,
		HTML:   html.String(),
		Target: "#no-tags-list",
		Insert: afterbegin,
	})
}
