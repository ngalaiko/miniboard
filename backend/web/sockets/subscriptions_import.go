package sockets

import (
	"bytes"
	"context"
	"fmt"

	"github.com/ngalaiko/miniboard/backend/tags"
	"github.com/ngalaiko/miniboard/backend/web/sockets/opml"
	"github.com/ngalaiko/miniboard/backend/web/templates"
)

func (h *Handler) onSubscriptionsImport(ctx context.Context, userID string, req *request, respond chan<- *response) {
	file, ok := req.Params["file"]
	if !ok {
		respond <- errResponse(req, fmt.Errorf("'file' parameter is missing"))
		return
	}

	parsed, err := opml.Parse([]byte(file))
	if err != nil {
		respond <- errResponse(req, fmt.Errorf("failed to parse file: %w", err))
		return
	}

	for _, opmlTag := range parsed.Tags {
		tag, err := h.getOrCreateTag(ctx, userID, opmlTag.Title)
		if err != nil {
			h.logger.Error("failed to get or create tag: %s", err)
			respond <- errResponse(req, errInternal)
			continue
		}

		for _, feed := range opmlTag.Feeds {
			subscription, err := h.createSubscription(ctx, userID, feed.URL, []string{tag.ID})
			if err != nil {
				respond <- errResponse(req, err)
				continue
			}

			html := &bytes.Buffer{}
			if err := templates.Subscription(html, subscription); err != nil {
				h.logger.Error("failed to render subscription: %s", err)
				respond <- errResponse(req, errInternal)
				continue
			}

			respond <- &response{
				ID:     req.ID,
				HTML:   html.String(),
				Target: "#no-tags-list",
				Insert: afterbegin,
			}
		}
	}
}

func (h *Handler) getOrCreateTag(ctx context.Context, userID string, title string) (*tags.Tag, error) {
	tag, err := h.tagsService.GetByTitle(ctx, userID, title)
	switch err {
	case nil:
		return tag, nil
	case tags.ErrNotFound:
		return h.tagsService.Create(ctx, userID, title)
	default:
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
}
