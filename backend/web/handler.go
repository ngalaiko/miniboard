package web

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/ngalaiko/miniboard/backend/items"
	"github.com/ngalaiko/miniboard/backend/subscriptions"
	"github.com/ngalaiko/miniboard/backend/tags"
	"github.com/ngalaiko/miniboard/backend/web/sockets"
	"github.com/ngalaiko/miniboard/backend/web/static"
	"github.com/ngalaiko/miniboard/backend/web/templates"
)

type tagsService interface {
	Create(ctx context.Context, userID string, title string) (*tags.Tag, error)
	GetByTitle(ctx context.Context, userID string, title string) (*tags.Tag, error)
	List(ctx context.Context, userID string, pageSize int, createdLT *time.Time) ([]*tags.Tag, error)
}

type subscriptionsService interface {
	Create(ctx context.Context, userID string, url *url.URL, tagIDs []string) (*subscriptions.UserSubscription, error)
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
	staticHandler := static.NewHandler(cfg.FS, log)
	templatesHandler := templates.NewHandler(log, itemsService, tagsService, subscriptionsService)
	socketsHandler := sockets.New(log, itemsService, tagsService, subscriptionsService)
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			switch r.URL.Path {
			default:
				http.NotFound(w, r)
			}
		case http.MethodGet:
			switch r.URL.Path {
			case "/api/ws/":
				socketsHandler.Receive().ServeHTTP(w, r)
			case "/users/":
				templatesHandler.ServeHTTP(w, r)
			default:
				staticHandler.ServeHTTP(w, r)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
