package web

import (
	"context"
	"net/http"
	"time"

	"github.com/ngalaiko/miniboard/backend/items"
	"github.com/ngalaiko/miniboard/backend/subscriptions"
	"github.com/ngalaiko/miniboard/backend/tags"
	"github.com/ngalaiko/miniboard/backend/web/static"
	"github.com/ngalaiko/miniboard/backend/web/templates"
)

type tagsService interface {
	List(ctx context.Context, userID string, pageSize int, createdLT *time.Time) ([]*tags.Tag, error)
}

type subscriptionsService interface {
	List(ctx context.Context, userID string, pageSize int, createdLT *time.Time) ([]*subscriptions.UserSubscription, error)
}

type itemsService interface {
	Get(ctx context.Context, id string, userID string) (*items.UserItem, error)
	List(ctx context.Context, userID string, pageSize int, createdLT *time.Time, subscriptionID *string, tagID *string) ([]*items.UserItem, error)
}

// Config contains web configuration.
type Config struct {
	FS bool `yaml:"fs"`
}

type logger interface {
	Debug(string, ...interface{})
	Error(string, ...interface{})
}

func NewHandler(cfg *Config, log logger, itemsService itemsService, tagsService tagsService, subscriptionsService subscriptionsService) http.HandlerFunc {
	useTemplates := map[string]bool{
		"/users":            true,
		"/users/":           true,
		"/users/index.html": true,
	}

	staticHandler := static.NewHandler(cfg.FS, log)
	templatesHandler := templates.NewHandler(log, itemsService, tagsService, subscriptionsService)
	return func(w http.ResponseWriter, r *http.Request) {
		if useTemplates[r.URL.Path] {
			templatesHandler.ServeHTTP(w, r)
		} else {
			staticHandler.ServeHTTP(w, r)
		}
	}
}
