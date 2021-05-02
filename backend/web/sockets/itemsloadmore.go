package sockets

import (
	"bytes"
	"fmt"
	"time"

	"golang.org/x/net/websocket"

	"github.com/ngalaiko/miniboard/backend/web/templates"
)

var (
	errInvalidCreatedLT = fmt.Errorf("failed to parse createdLt param")
	errInternal         = fmt.Errorf("internal error")
)

func (s *Sockets) itemsLoadmore(ws *websocket.Conn, userID string, req *Request) {
	ctx := ws.Request().Context()

	var tagID, subscriptionID *string
	if id, ok := req.Params["tagId"]; ok {
		tagID = &id
	}

	if id, ok := req.Params["subscriptionId"]; ok {
		subscriptionID = &id
	}

	var createdLT *time.Time
	if cltRaw, ok := req.Params["createdLt"]; ok {
		clt, err := time.Parse(time.RFC3339, cltRaw)
		if err != nil {
			s.respond(ws, errResponse(req, errInvalidCreatedLT))
			return
		}
		createdLT = &clt
	}

	items, err := s.itemsService.List(ctx, userID, 50, createdLT, subscriptionID, tagID)
	if err != nil {
		s.logger.Error("failed to list items: %s", err)
		s.respond(ws, errResponse(req, errInternal))
		return
	}

	for _, item := range items {
		html := &bytes.Buffer{}
		if err := templates.Item(html, item); err != nil {
			s.logger.Error("failed to render reader: %s", err)
			s.respond(ws, errResponse(req, errInternal))
			return
		}
		s.respond(ws, &Response{
			ID:     req.ID,
			HTML:   html.String(),
			Target: "#items-list",
			Insert: Beforeend,
		})
	}
}
