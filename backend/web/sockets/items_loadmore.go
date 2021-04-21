package sockets

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/ngalaiko/miniboard/backend/web/templates"
)

var errInvalidCreatedLT = fmt.Errorf("failed to parse createdLt param")

func (h *Handler) loadItems(ctx context.Context, userID string, req *request) ([]*response, error) {
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
			return nil, errInvalidCreatedLT
		}
		createdLT = &clt
	}

	items, err := h.itemsService.List(ctx, userID, 50, createdLT, subscriptionID, tagID)
	switch {
	case err == nil:
		rr := make([]*response, 0, len(items)+1)
		for _, item := range items {
			html := &bytes.Buffer{}
			if err := templates.Item(html, item); err != nil {
				h.logger.Error("failed to render reader: %s", err)
				return nil, errInternal
			}
			rr = append(rr, &response{
				ID:     req.ID,
				HTML:   html.String(),
				Target: "#items-list",
				Insert: beforeend,
			})
		}
		return rr, nil
	default:
		h.logger.Error("failed to get operation: %s", err)
		return nil, errInternal
	}
}
