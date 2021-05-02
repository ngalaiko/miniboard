package sockets

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"

	"golang.org/x/net/websocket"

	"github.com/ngalaiko/miniboard/backend/subscriptions"
	"github.com/ngalaiko/miniboard/backend/tags"
	"github.com/ngalaiko/miniboard/backend/web/render"
	"github.com/ngalaiko/miniboard/backend/web/sockets/opml"
)

func (s *Sockets) subscriptionsImport(ws *websocket.Conn, userID string, req *request) {
	ctx := ws.Request().Context()

	file, ok := req.Params["file"]
	if !ok {
		s.respond(ws, errResponse(req, fmt.Errorf("'file' parameter is missing")))
		return
	}

	parsed, err := opml.Parse([]byte(file))
	if err != nil {
		s.respond(ws, errResponse(req, fmt.Errorf("failed to parse file: %w", err)))
		return
	}

	for _, opmlTag := range parsed.Tags {
		tag, err := s.getOrCreateTag(ctx, userID, opmlTag.Title, req)
		if err != nil {
			s.logger.Error("failed to get or create tag: %s", err)
			s.respond(ws, errResponse(req, errInternal))
			continue
		}

		for _, feed := range opmlTag.Feeds {
			url, err := url.ParseRequestURI(feed.URL)
			if err != nil {
				s.respond(ws, errResponse(req, fmt.Errorf("failed to parse url: %w", err)))
				continue
			}
			subscription, err := s.subscriptionsService.Create(ctx, userID, url, []string{tag.ID})
			switch {
			case err == nil:
			case errors.Is(err, subscriptions.ErrFailedToDownloadSubscription),
				errors.Is(err, subscriptions.ErrAlreadyExists),
				errors.Is(err, subscriptions.ErrFailedToParseSubscription):
				s.respond(ws, errResponse(req, err))
				continue
			default:
				s.logger.Error("failed to create subscription: %s", err)
				s.respond(ws, errResponse(req, errInternal))
				continue
			}

			html := &bytes.Buffer{}
			if err := render.Subscription(html, subscription); err != nil {
				s.logger.Error("failed to render subscription: %s", err)
				s.respond(ws, errResponse(req, errInternal))
				continue
			}

			s.broadcast(userID, &response{
				ID:     req.ID,
				HTML:   html.String(),
				Target: fmt.Sprintf("#%s-children", tag.ID),
				Insert: beforeend,
			})
		}
	}
}

func (s *Sockets) getOrCreateTag(ctx context.Context, userID string, title string, req *request) (*tags.Tag, error) {
	tag, err := s.tagsService.GetByTitle(ctx, userID, title)
	switch {
	case err == nil:
		return tag, nil
	case errors.Is(err, tags.ErrNotFound):
		newTag, err := s.tagsService.Create(ctx, userID, title)
		if err != nil {
			return newTag, err
		}
		html := &bytes.Buffer{}
		if err := render.Tag(html, newTag); err != nil {
			s.logger.Error("failed to render tag: %s", err)
			return nil, errInternal
		}
		s.broadcast(userID, &response{
			ID:     req.ID,
			HTML:   fmt.Sprintf(`<div id="%s-children" class="tag-subscriptions" hidden></div>`, newTag.ID),
			Target: "#tags-list",
			Insert: afterbegin,
		})
		s.broadcast(userID, &response{
			ID:     req.ID,
			HTML:   html.String(),
			Target: "#tags-list",
			Insert: afterbegin,
		})
		return newTag, err
	default:
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}
}
