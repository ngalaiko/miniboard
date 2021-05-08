package web

import (
	"context"
	"net/http"
	"net/url"
	"time"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/httpx"
	"github.com/ngalaiko/miniboard/backend/items"
	"github.com/ngalaiko/miniboard/backend/subscriptions"
	"github.com/ngalaiko/miniboard/backend/tags"
	"github.com/ngalaiko/miniboard/backend/web/render"
	"github.com/ngalaiko/miniboard/backend/web/sockets"
	"github.com/ngalaiko/miniboard/backend/web/static"
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
	Get(ctx context.Context, userID string, id string) (*items.UserItem, error)
	List(ctx context.Context, userID string, pageSize int, createdLT *time.Time, subscriptionID *string, tagID *string) ([]*items.UserItem, error)
}

// Config contains web configuration.
type Config struct {
	FS bool `yaml:"fs"`
}

type logger interface {
	Debug(string, ...interface{})
	Info(string, ...interface{})
	Error(string, ...interface{})
}

//nolint: funlen
func NewHandler(
	cfg *Config,
	log logger,
	itemsService itemsService,
	tagsService tagsService,
	subscriptionsService subscriptionsService,
	usersService usersService,
	jwtService jwtService,
) http.Handler {
	if cfg.FS {
		log.Debug("serving files from fs")
	} else {
		log.Debug("serving files from memory")
	}

	staticHandler := static.NewHandler(cfg.FS)
	render := render.Load(cfg.FS)
	socketsHandler := sockets.New(log, itemsService, tagsService, subscriptionsService, render)
	signupHandler := signupHandler(log, usersService, jwtService, render)
	loginHandler := loginHandler(log, usersService, jwtService, render)

	r := chi.NewRouter()
	r.Use(Log(log))
	r.Use(middleware.Recoverer)
	r.Use(Authenticate(jwtService, log))

	r.Post("/login/", loginHandler)
	r.Post("/signup/", signupHandler)
	r.Get("/api/ws/", socketsHandler.Receive())

	r.With(requireAuth).Get("/users/", func(w http.ResponseWriter, r *http.Request) {
		usersHandler(log, render,
			withURL(r.URL),
			withAllSubscriptionsByTags(tagsService, subscriptionsService),
		).ServeHTTP(w, r)
	})
	r.With(requireAuth).Get("/users/{itemID}/", func(w http.ResponseWriter, r *http.Request) {
		usersHandler(log, render,
			withURL(r.URL),
			withAllSubscriptionsByTags(tagsService, subscriptionsService),
			withItem(itemsService, chi.URLParam(r, "itemID")),
		).ServeHTTP(w, r)
	})
	r.With(requireAuth).Get("/users/items/", usersHandler(log, render, withAllItems(itemsService)))
	r.With(requireAuth).Get("/users/items/{itemID}/", func(w http.ResponseWriter, r *http.Request) {
		usersHandler(log, render,
			withURL(r.URL),
			withAllSubscriptionsByTags(tagsService, subscriptionsService),
			withAllItems(itemsService),
			withItem(itemsService, chi.URLParam(r, "itemID")),
		).ServeHTTP(w, r)
	})
	r.With(requireAuth).Get("/users/tags/{tagID}/items/", func(w http.ResponseWriter, r *http.Request) {
		usersHandler(log, render,
			withURL(r.URL),
			withAllSubscriptionsByTags(tagsService, subscriptionsService),
			withItemsByTagID(itemsService, chi.URLParam(r, "tagID")),
		).ServeHTTP(w, r)
	})
	r.With(requireAuth).Get("/users/tags/{tagID}/items/{itemID}/", func(w http.ResponseWriter, r *http.Request) {
		usersHandler(log, render,
			withURL(r.URL),
			withAllSubscriptionsByTags(tagsService, subscriptionsService),
			withItemsByTagID(itemsService, chi.URLParam(r, "tagID")),
			withItem(itemsService, chi.URLParam(r, "itemID")),
		).ServeHTTP(w, r)
	})
	r.With(requireAuth).Get("/users/subscriptions/{subscriptionID}/items/", func(w http.ResponseWriter, r *http.Request) {
		usersHandler(log, render,
			withURL(r.URL),
			withAllSubscriptionsByTags(tagsService, subscriptionsService),
			withItemsBySubscriptionID(itemsService, chi.URLParam(r, "subscriptionID")),
		).ServeHTTP(w, r)
	})
	r.With(requireAuth).Get("/users/subscriptions/{subscriptionID}/items/{itemID}/", func(w http.ResponseWriter, r *http.Request) {
		usersHandler(log, render,
			withURL(r.URL),
			withAllSubscriptionsByTags(tagsService, subscriptionsService),
			withItemsBySubscriptionID(itemsService, chi.URLParam(r, "subscriptionID")),
			withItem(itemsService, chi.URLParam(r, "itemID")),
		).ServeHTTP(w, r)
	})

	r.Get("/logout/", func(w http.ResponseWriter, r *http.Request) {
		removeCookie(w, r.TLS != nil)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
	r.Get("/signup/", func(w http.ResponseWriter, r *http.Request) {
		if err := render.SignupPage(w, nil); err != nil {
			log.Error("failed to render login page: %s", err)
			httpx.InternalError(w, log)
		}
	})
	r.Get("/login/", func(w http.ResponseWriter, r *http.Request) {
		_, authorized := authorizations.FromContext(r.Context())
		if authorized {
			http.Redirect(w, r, "/users/", http.StatusSeeOther)
			return
		}
		if err := render.LoginPage(w, nil); err != nil {
			log.Error("failed to render login page: %s", err)
			httpx.InternalError(w, log)
		}
	})
	r.Get("/*", staticHandler)

	return r
}
