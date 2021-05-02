package web

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/httpx"
	"github.com/ngalaiko/miniboard/backend/items"
	"github.com/ngalaiko/miniboard/backend/web/templates"
)

func usersHandler(log logger, itemsService itemsService, tagsService tagsService, subscriptionsService subscriptionsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, _ := authorizations.FromContext(r.Context())
		item, err := loadItem(r.Context(), token.UserID, itemsService, r.URL.Query().Get("item"))
		if err != nil {
			log.Error("failed to load item: %s", err)
			httpx.InternalError(w, log)
			return
		}
		items, err := itemsService.List(r.Context(), token.UserID, 50, nil, optionalURLParam(r.URL.Query(), "subscription"), optionalURLParam(r.URL.Query(), "tag"))
		if err != nil {
			log.Error("failed to load items: %s", err)
			httpx.InternalError(w, log)
			return
		}
		tags, err := tagsService.List(r.Context(), token.UserID, 10000, nil)
		if err != nil {
			log.Error("failed to load tags: %s", err)
			httpx.InternalError(w, log)
			return
		}
		subscriptions, err := subscriptionsService.List(r.Context(), token.UserID, 10000, nil)
		if err != nil {
			log.Error("failed to load subscriptions: %s", err)
			httpx.InternalError(w, log)
			return
		}
		if err := templates.UsersPage(w, item, items, tags, subscriptions); err != nil {
			log.Error("failed to render users page: %s", err)
			httpx.InternalError(w, log)
			return
		}
	}
}

func optionalURLParam(params url.Values, name string) *string {
	if pp, found := params[name]; found && len(pp) > 0 {
		return &pp[0]
	}
	return nil
}

func loadItem(ctx context.Context, userID string, itemsService itemsService, id string) (*items.UserItem, error) {
	if id == "" {
		return nil, nil
	}

	item, err := itemsService.Get(ctx, userID, id)
	switch {
	case err == nil:
		return item, nil
	case errors.Is(err, items.ErrNotFound):
		return nil, nil
	default:
		return nil, err
	}
}
