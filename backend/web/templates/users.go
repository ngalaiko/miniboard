package templates

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/ngalaiko/miniboard/backend/items"
	"github.com/ngalaiko/miniboard/backend/subscriptions"
	"github.com/ngalaiko/miniboard/backend/tags"
)

type itemsService interface {
	Get(ctx context.Context, id string, userID string) (*items.UserItem, error)
	List(ctx context.Context, userID string, pageSize int, createdLT *time.Time, subscriptionID *string, tagID *string) ([]*items.UserItem, error)
}

type subscriptionsService interface {
	List(ctx context.Context, userID string, pageSize int, createdLT *time.Time) ([]*subscriptions.UserSubscription, error)
}

type tagsService interface {
	List(ctx context.Context, userID string, pageSize int, createdLT *time.Time) ([]*tags.Tag, error)
}

type usersData struct {
	Item                 *items.UserItem
	Items                []*items.UserItem
	Tags                 []*tags.Tag
	SubscriptionsByTagID map[string][]*subscriptions.UserSubscription
}

func loadUsersData(
	r *http.Request,
	userID string,
	itemsService itemsService,
	tagsService tagsService,
	subscriptionsService subscriptionsService,
) (*usersData, error) {
	item, err := loadItem(r.Context(), userID, itemsService, r.URL.Query().Get("item"))
	if err != nil {
		return nil, err
	}
	items, err := loadItems(r.Context(), userID, itemsService, optionalURLParam(r.URL.Query(), "tag"), optionalURLParam(r.URL.Query(), "subscription"))
	if err != nil {
		return nil, err
	}
	tags, err := loadTags(r.Context(), userID, tagsService)
	if err != nil {
		return nil, err
	}
	subscriptionsByTagID, err := loadSubscriptionsByTagID(r.Context(), userID, subscriptionsService)
	if err != nil {
		return nil, err
	}
	return &usersData{
		Item:                 item,
		Items:                items,
		Tags:                 tags,
		SubscriptionsByTagID: subscriptionsByTagID,
	}, nil
}

func optionalURLParam(params url.Values, name string) *string {
	if pp, found := params[name]; found && len(pp) > 0 {
		return &pp[0]
	}
	return nil
}

func loadTags(ctx context.Context, userID string, tagsService tagsService) ([]*tags.Tag, error) {
	return tagsService.List(ctx, userID, 10000, nil)
}

func loadSubscriptionsByTagID(ctx context.Context, userID string, subscriptionsService subscriptionsService) (map[string][]*subscriptions.UserSubscription, error) {
	ss, err := subscriptionsService.List(ctx, userID, 10000, nil)
	if err != nil {
		return nil, err
	}
	subscriptionsByTagID := map[string][]*subscriptions.UserSubscription{}
	for _, s := range ss {
		if len(s.TagIDs) == 0 {
			subscriptionsByTagID[""] = append(subscriptionsByTagID[""], s)
		}
		for _, tagID := range s.TagIDs {
			subscriptionsByTagID[tagID] = append(subscriptionsByTagID[tagID], s)
		}
	}
	return subscriptionsByTagID, err
}

func loadItems(ctx context.Context, userID string, itemsService itemsService, tagID *string, subscriptionID *string) ([]*items.UserItem, error) {
	return itemsService.List(ctx, userID, 50, nil, subscriptionID, tagID)
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
