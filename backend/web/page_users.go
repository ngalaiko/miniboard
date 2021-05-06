package web

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/httpx"
	"github.com/ngalaiko/miniboard/backend/subscriptions"
	"github.com/ngalaiko/miniboard/backend/tags"
	"github.com/ngalaiko/miniboard/backend/web/render"
)

type usersDataProvider func(context.Context, *render.UsersData) error

func withAllSubscriptionsByTags(tagsService tagsService, subscriptionsService subscriptionsService) usersDataProvider {
	return func(ctx context.Context, data *render.UsersData) error {
		token := authorizations.MustFromContext(ctx)
		tt, err := tagsService.List(ctx, token.UserID, 10000, nil)
		if err != nil {
			return fmt.Errorf("failed to load tags for user '%s': %w", token.UserID, err)
		}
		ss, err := subscriptionsService.List(ctx, token.UserID, 10000, nil)
		if err != nil {
			return fmt.Errorf("failed to load subscriptions for user '%s': %w", token.UserID, err)
		}
		subscriptionsByTagID := map[string][]*subscriptions.UserSubscription{}
		noTagSubscriptions := []*subscriptions.UserSubscription{}
		for _, s := range ss {
			if len(s.TagIDs) == 0 {
				noTagSubscriptions = append(noTagSubscriptions, s)
			}
			for _, tagID := range s.TagIDs {
				subscriptionsByTagID[tagID] = append(subscriptionsByTagID[tagID], s)
			}
		}
		tagsByTagID := map[string]*tags.Tag{}
		for _, tag := range tt {
			tagsByTagID[tag.ID] = tag
		}
		tagsSubscriptions := []*render.UsersTag{}
		for tagID, ss := range subscriptionsByTagID {
			tagsSubscriptions = append(tagsSubscriptions, &render.UsersTag{
				Tag:           tagsByTagID[tagID],
				Subscriptions: ss,
			})
		}
		sort.Slice(tagsSubscriptions, func(i, j int) bool {
			ci := tagsSubscriptions[i].Tag.Created
			cj := tagsSubscriptions[j].Tag.Created
			return ci.Before(cj)
		})

		data.Tags = tagsSubscriptions
		data.Subscriptions = noTagSubscriptions

		return nil
	}
}

func withItemsBySubscriptionID(itemsService itemsService, subscriptionID string) usersDataProvider {
	return func(ctx context.Context, data *render.UsersData) error {
		token := authorizations.MustFromContext(ctx)
		items, err := itemsService.List(ctx, token.UserID, 50, nil, &subscriptionID, nil)
		if err != nil {
			return fmt.Errorf("failed to get items for user '%s': %w", token.UserID, err)
		}

		data.Items = items

		return nil
	}
}

func withItemsByTagID(itemsService itemsService, tagID string) usersDataProvider {
	return func(ctx context.Context, data *render.UsersData) error {
		token := authorizations.MustFromContext(ctx)
		items, err := itemsService.List(ctx, token.UserID, 50, nil, nil, &tagID)
		if err != nil {
			return fmt.Errorf("failed to get items for user '%s': %w", token.UserID, err)
		}

		data.Items = items

		return nil
	}
}

func withAllItems(itemsService itemsService) usersDataProvider {
	return func(ctx context.Context, data *render.UsersData) error {
		token := authorizations.MustFromContext(ctx)
		items, err := itemsService.List(ctx, token.UserID, 50, nil, nil, nil)
		if err != nil {
			return fmt.Errorf("failed to get items for user '%s': %w", token.UserID, err)
		}

		data.Items = items

		return nil
	}
}

func withItem(itemsService itemsService, itemID string) usersDataProvider {
	return func(ctx context.Context, data *render.UsersData) error {
		token := authorizations.MustFromContext(ctx)
		item, err := itemsService.Get(ctx, token.UserID, itemID)
		if err != nil {
			return fmt.Errorf("failed to get item '%s' for user '%s': %w", itemID, token.UserID, err)
		}

		data.Item = item

		return nil
	}
}

func usersHandler(log logger, rdr *render.Templates, dataProviders ...usersDataProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		usersData := &render.UsersData{}
		for _, provide := range dataProviders {
			if err := provide(r.Context(), usersData); err != nil {
				log.Error("failed to load data for users page: %s", err)
				httpx.InternalError(w, log)
				return
			}
		}
		if err := rdr.UsersPage(w, usersData); err != nil {
			log.Error("failed to render users page: %s", err)
			httpx.InternalError(w, log)
			return
		}
	}
}
