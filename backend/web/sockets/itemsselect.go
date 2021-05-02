package sockets

import (
	"bytes"
	"errors"
	"fmt"

	"golang.org/x/net/websocket"

	"github.com/ngalaiko/miniboard/backend/items"
	"github.com/ngalaiko/miniboard/backend/web/render"
)

func (s *Sockets) itemsSelect(ws *websocket.Conn, userID string, req *request) {
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
		if err := render.Reader(html, &item.Item); err != nil {
			s.logger.Error("failed to render reader: %s", err)
			s.respond(ws, errResponse(req, errInternal))
			return
		}
		s.respond(ws, &response{
			ID:     req.ID,
			HTML:   html.String(),
			Reset:  true,
			Target: "#reader",
			Insert: afterbegin,
		})
	case errors.Is(err, items.ErrNotFound):
		s.respond(ws, errResponse(req, fmt.Errorf("not foud")))
	default:
		s.logger.Error("failed to get operation: %s", err)
		s.respond(ws, errResponse(req, errInternal))
	}
}
