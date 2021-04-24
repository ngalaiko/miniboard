package handlers

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/items"
	"github.com/ngalaiko/miniboard/backend/web/sockets"
	"github.com/ngalaiko/miniboard/backend/web/templates"
)

var (
	errInvalidCreatedLT = fmt.Errorf("failed to parse createdLt param")
	errInternal         = fmt.Errorf("internal error")
)

type logger interface {
	Error(string, ...interface{})
}

type itemsService interface {
	Get(ctx context.Context, id string, userID string) (*items.UserItem, error)
	List(ctx context.Context, userID string, pageSize int, createdLT *time.Time, subscriptionID *string, tagID *string) ([]*items.UserItem, error)
}

func ItemsLoadmore(logger logger, itemsService itemsService) sockets.Handler {
	return func(ctx context.Context, req *sockets.Request, respond sockets.Respond, _ sockets.Broadcast) {
		token, auth := authorizations.FromContext(ctx)
		if !auth {
			respond(&sockets.Response{
				ID:    req.ID,
				Error: "unauthorized",
			})
			return
		}

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
				respond(sockets.Error(req, errInvalidCreatedLT))
				return
			}
			createdLT = &clt
		}

		items, err := itemsService.List(ctx, token.UserID, 50, createdLT, subscriptionID, tagID)
		if err != nil {
			logger.Error("failed to list items: %s", err)
			respond(sockets.Error(req, errInternal))
			return
		}

		for _, item := range items {
			html := &bytes.Buffer{}
			if err := templates.Item(html, item); err != nil {
				logger.Error("failed to render reader: %s", err)
				respond(sockets.Error(req, errInternal))
				return
			}
			respond(&sockets.Response{
				ID:     req.ID,
				HTML:   html.String(),
				Target: "#items-list",
				Insert: sockets.Beforeend,
			})
		}
	}
}
