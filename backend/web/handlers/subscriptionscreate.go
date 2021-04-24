package handlers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/subscriptions"
	"github.com/ngalaiko/miniboard/backend/web/sockets"
	"github.com/ngalaiko/miniboard/backend/web/templates"
)

type subscriptionsService interface {
	Create(ctx context.Context, userID string, url *url.URL, tagIDs []string) (*subscriptions.UserSubscription, error)
}

func SubscriptionsCreate(logger logger, subscriptionsService subscriptionsService) sockets.Handler {
	return func(ctx context.Context, req *sockets.Request, respond sockets.Respond, _ sockets.Broadcast) {
		token, auth := authorizations.FromContext(ctx)
		if !auth {
			respond(&sockets.Response{
				ID:    req.ID,
				Error: "unauthorized",
			})
			return
		}
		urlParam, ok := req.Params["url"]
		if !ok {
			respond(sockets.Error(req, fmt.Errorf("'url' parameter is missing")))
			return
		}

		url, err := url.ParseRequestURI(urlParam)
		if err != nil {
			respond(sockets.Error(req, fmt.Errorf("'url' is invalid: %s", err)))
			return
		}

		subscription, err := subscriptionsService.Create(ctx, token.UserID, url, nil)
		switch {
		case err == nil:
		case errors.Is(err, subscriptions.ErrFailedToDownloadSubscription),
			errors.Is(err, subscriptions.ErrAlreadyExists),
			errors.Is(err, subscriptions.ErrFailedToParseSubscription):
			respond(sockets.Error(req, err))
			return
		default:
			logger.Error("failed to create subscription: %s", err)
			respond(sockets.Error(req, errInternal))
			return
		}

		html := &bytes.Buffer{}
		if err := templates.Subscription(html, subscription); err != nil {
			logger.Error("failed to render subscription: %s", err)
			sockets.Error(req, errInternal)
			return
		}
		respond(&sockets.Response{
			ID:     req.ID,
			HTML:   html.String(),
			Target: "#no-tags-list",
			Insert: sockets.Afterbegin,
		})
	}
}
