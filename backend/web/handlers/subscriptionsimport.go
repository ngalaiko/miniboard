package handlers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/subscriptions"
	"github.com/ngalaiko/miniboard/backend/tags"
	"github.com/ngalaiko/miniboard/backend/web/handlers/opml"
	"github.com/ngalaiko/miniboard/backend/web/sockets"
	"github.com/ngalaiko/miniboard/backend/web/templates"
)

type tagsService interface {
	Create(ctx context.Context, userID string, title string) (*tags.Tag, error)
	GetByTitle(ctx context.Context, userID string, title string) (*tags.Tag, error)
}

//nolint: funlen,gocognit
func SubscriptionsImport(logger logger, subscriptionsService subscriptionsService, tagsService tagsService) sockets.Handler {
	getOrCreateTag := func(ctx context.Context, userID string, title string, req *sockets.Request, respond func(*sockets.Response)) (*tags.Tag, error) {
		tag, err := tagsService.GetByTitle(ctx, userID, title)
		switch err {
		case nil:
			return tag, nil
		case tags.ErrNotFound:
			newTag, err := tagsService.Create(ctx, userID, title)
			if err != nil {
				return newTag, err
			}
			html := &bytes.Buffer{}
			if err := templates.Tag(html, newTag); err != nil {
				logger.Error("failed to render tag: %s", err)
				return nil, errInternal
			}
			respond(&sockets.Response{
				ID:     req.ID,
				HTML:   fmt.Sprintf(`<div id="%s-children" class="tag-subscriptions" hidden></div>`, newTag.ID),
				Target: "#tags-list",
				Insert: sockets.Afterbegin,
			})
			respond(&sockets.Response{
				ID:     req.ID,
				HTML:   html.String(),
				Target: "#tags-list",
				Insert: sockets.Afterbegin,
			})
			return newTag, err
		default:
			return nil, fmt.Errorf("failed to get tag: %w", err)
		}
	}
	return func(ctx context.Context, req *sockets.Request, respond sockets.Respond, broadcast sockets.Broadcast) {
		token, auth := authorizations.FromContext(ctx)
		if !auth {
			respond(&sockets.Response{
				ID:    req.ID,
				Error: "unauthorized",
			})
			return
		}
		file, ok := req.Params["file"]
		if !ok {
			respond(sockets.Error(req, fmt.Errorf("'file' parameter is missing")))
			return
		}

		parsed, err := opml.Parse([]byte(file))
		if err != nil {
			respond(sockets.Error(req, fmt.Errorf("failed to parse file: %w", err)))
			return
		}

		for _, opmlTag := range parsed.Tags {
			tag, err := getOrCreateTag(ctx, token.UserID, opmlTag.Title, req, broadcast)
			if err != nil {
				logger.Error("failed to get or create tag: %s", err)
				respond(sockets.Error(req, errInternal))
				continue
			}

			for _, feed := range opmlTag.Feeds {
				url, err := url.ParseRequestURI(feed.URL)
				if err != nil {
					respond(sockets.Error(req, fmt.Errorf("failed to parse url: %s", err)))
					continue
				}
				subscription, err := subscriptionsService.Create(ctx, token.UserID, url, []string{tag.ID})
				switch {
				case err == nil:
				case errors.Is(err, subscriptions.ErrFailedToDownloadSubscription),
					errors.Is(err, subscriptions.ErrAlreadyExists),
					errors.Is(err, subscriptions.ErrFailedToParseSubscription):
					respond(sockets.Error(req, err))
					continue
				default:
					logger.Error("failed to create subscription: %s", err)
					respond(sockets.Error(req, errInternal))
					continue
				}

				html := &bytes.Buffer{}
				if err := templates.Subscription(html, subscription); err != nil {
					logger.Error("failed to render subscription: %s", err)
					respond(sockets.Error(req, errInternal))
					continue
				}

				broadcast(&sockets.Response{
					ID:     req.ID,
					HTML:   html.String(),
					Target: fmt.Sprintf("#%s-children", tag.ID),
					Insert: sockets.Beforeend,
				})
			}
		}
	}
}
