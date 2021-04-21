package sockets

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/ngalaiko/miniboard/backend/items"
	"github.com/ngalaiko/miniboard/backend/web/templates"
)

func (h *Handler) onItemSelected(ctx context.Context, userID string, req *request) *response {
	id, ok := req.Params["id"]
	if !ok {
		return errResponse(req, fmt.Errorf("'id' parameter is missing"))
	}
	item, err := h.itemsService.Get(ctx, userID, id)
	switch {
	case err == nil:
		html := &bytes.Buffer{}
		if err := templates.Reader(html, &item.Item); err != nil {
			h.logger.Error("failed to render reader: %s", err)
			return errResponse(req, errInternal)
		}
		return &response{
			ID:     req.ID,
			HTML:   html.String(),
			Reset:  true,
			Target: "#reader",
			Insert: afterbegin,
		}
	case errors.Is(err, items.ErrNotFound):
		return errResponse(req, errNotFound)
	default:
		h.logger.Error("failed to get operation: %s", err)
		return errResponse(req, errInternal)
	}
}
