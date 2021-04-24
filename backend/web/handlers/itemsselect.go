package handlers

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/items"
	"github.com/ngalaiko/miniboard/backend/web/sockets"
	"github.com/ngalaiko/miniboard/backend/web/templates"
)

var errNotFound = fmt.Errorf("not found")

func ItemsSelect(logger logger, itemsService itemsService) sockets.Handler {
	return func(ctx context.Context, req *sockets.Request, respond sockets.Respond, _ sockets.Broadcast) {
		token, auth := authorizations.FromContext(ctx)
		if !auth {
			respond(sockets.Error(req, fmt.Errorf("unauthorized")))
			return
		}
		id, ok := req.Params["id"]
		if !ok {
			respond(sockets.Error(req, fmt.Errorf("'id' parameter is missing")))
			return
		}
		item, err := itemsService.Get(ctx, token.UserID, id)
		switch {
		case err == nil:
			html := &bytes.Buffer{}
			if err := templates.Reader(html, &item.Item); err != nil {
				logger.Error("failed to render reader: %s", err)
				respond(sockets.Error(req, errInternal))
				return
			}
			respond(&sockets.Response{
				ID:     req.ID,
				HTML:   html.String(),
				Reset:  true,
				Target: "#reader",
				Insert: sockets.Afterbegin,
			})
		case errors.Is(err, items.ErrNotFound):
			respond(sockets.Error(req, errNotFound))
		default:
			logger.Error("failed to get operation: %s", err)
			respond(sockets.Error(req, errInternal))
		}
	}
}
