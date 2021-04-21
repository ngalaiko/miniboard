package sockets

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/ngalaiko/miniboard/backend/subscriptions"
	"github.com/ngalaiko/miniboard/backend/web/templates"
)

func (h *Handler) onSubscriptionsCreated(ctx context.Context, userID string, req *request) *response {
	urlParam, ok := req.Params["url"]
	if !ok {
		return errResponse(req, fmt.Errorf("'url' parameter is missing"))
	}
	url, err := url.ParseRequestURI(urlParam)
	if err != nil {
		return errResponse(req, fmt.Errorf("'url' is invalid: %s", err))
	}

	subscription, err := h.subscriptionsService.Create(ctx, userID, url, nil)
	switch {
	case err == nil:
	case errors.Is(err, subscriptions.ErrFailedToDownloadSubscription),
		errors.Is(err, subscriptions.ErrAlreadyExists),
		errors.Is(err, subscriptions.ErrFailedToParseSubscription):
		return errResponse(req, err)
	default:
		h.logger.Error("failed to create subscription: %s", err)
		return errResponse(req, fmt.Errorf("'url' is invalid: %s", err))
	}

	html := &bytes.Buffer{}
	if err := templates.Subscription(html, subscription); err != nil {
		h.logger.Error("failed to render subscription: %s", err)
		return errResponse(req, errInternal)
	}

	return &response{
		ID:     req.ID,
		HTML:   html.String(),
		Target: "#no-tags-list",
		Insert: afterbegin,
	}
}
