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
		tag, err := h.getOrCreateTag(ctx, userID, opmlTag.Title, req, respond)
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
				Target: fmt.Sprintf("#%s-children", tag.ID),
				Insert: beforeend,
			}
		}
	}
}

func (h *Handler) getOrCreateTag(ctx context.Context, userID string, title string, req *request, respond chan<- *response) (*tags.Tag, error) {
	tag, err := h.tagsService.GetByTitle(ctx, userID, title)
	switch err {
	case nil:
		return tag, nil
	case tags.ErrNotFound:
		newTag, err := h.tagsService.Create(ctx, userID, title)
		if err != nil {
			return newTag, err
		}
		html := &bytes.Buffer{}
		if err := templates.Tag(html, newTag); err != nil {
			h.logger.Error("failed to render tag: %s", err)
			return nil, errInternal
		}
		respond <- &response{
			ID:     req.ID,
			HTML:   fmt.Sprintf(`<div id="%s-children" class="tag-subscriptions" hidden></div>`, newTag.ID),
			Target: "#tags-list",
			Insert: afterbegin,
		}
		respond <- &response{
			ID:     req.ID,
			HTML:   html.String(),
			Target: "#tags-list",
			Insert: afterbegin,
		}
		return newTag, err
	default:
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
}
