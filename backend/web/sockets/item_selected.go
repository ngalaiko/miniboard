package sockets

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/ngalaiko/miniboard/backend/items"
	"github.com/ngalaiko/miniboard/backend/web/components"
)

func (h *Handler) onItemSelected(ctx context.Context, userID string, req *request) (*response, error) {
	id, ok := req.Params["id"]
	if !ok {
		return nil, fmt.Errorf("'id' parameter is missing")
	}
	item, err := h.itemsService.Get(ctx, userID, id)
	switch {
	case err == nil:
		html := &bytes.Buffer{}
		if err := components.Reader(html, &item.Item); err != nil {
			h.logger.Error("failed to render reader: %s", err)
			return nil, errInternal
		}
		return &response{
			ID:     req.ID,
			HTML:   html.String(),
			Reset:  true,
			Target: "#reader",
			Insert: afterbegin,
		}, nil
	case errors.Is(err, items.ErrNotFound):
		return nil, errNotFound
	default:
		h.logger.Error("failed to get operation: %s", err)
		return nil, errInternal
	}
}
