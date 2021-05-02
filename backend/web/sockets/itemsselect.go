package sockets

import (
	"bytes"
	"errors"
	"fmt"

	"golang.org/x/net/websocket"

	"github.com/ngalaiko/miniboard/backend/items"
	"github.com/ngalaiko/miniboard/backend/web/templates"
)

func (s *Sockets) itemsSelect(ws *websocket.Conn, userID string, req *Request) {
	ctx := ws.Request().Context()

	id, ok := req.Params["id"]
	if !ok {
		s.respond(ws, errResponse(req, fmt.Errorf("'id' parameter is missing")))
		return
	}
	item, err := s.itemsService.Get(ctx, userID, id)
	switch {
	case err == nil:
		html := &bytes.Buffer{}
		if err := templates.Reader(html, &item.Item); err != nil {
			s.logger.Error("failed to render reader: %s", err)
			s.respond(ws, errResponse(req, errInternal))
			return
		}
		s.respond(ws, &Response{
			ID:     req.ID,
			HTML:   html.String(),
			Reset:  true,
			Target: "#reader",
			Insert: Afterbegin,
		})
	case errors.Is(err, items.ErrNotFound):
		s.respond(ws, errResponse(req, fmt.Errorf("not foud")))
	default:
		s.logger.Error("failed to get operation: %s", err)
		s.respond(ws, errResponse(req, errInternal))
	}
}
