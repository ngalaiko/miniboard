package backend

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"golang.org/x/sync/errgroup"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/crawler"
	"github.com/ngalaiko/miniboard/backend/db"
	"github.com/ngalaiko/miniboard/backend/httpx"
	"github.com/ngalaiko/miniboard/backend/items"
	"github.com/ngalaiko/miniboard/backend/logger"
	"github.com/ngalaiko/miniboard/backend/operations"
	"github.com/ngalaiko/miniboard/backend/subscriptions"
	"github.com/ngalaiko/miniboard/backend/tags"
	"github.com/ngalaiko/miniboard/backend/users"
)

// Config contains all server configuration.
type Config struct {
	Authorizations *authorizations.Config `yaml:"authorizations"`
	DB             *db.Config             `yaml:"db"`
	HTTP           *httpx.Config          `yaml:"http"`
	Operations     *operations.Config     `yaml:"operations"`
	Subscriptions  *subscriptions.Config  `yaml:"subscriptions"`
	Users          *users.Config          `yaml:"users"`
}

// Server is the main object.
type Server struct {
	log                   *logger.Logger
	db                    *sql.DB
	httpServer            *httpx.Server
	authorizationsService *authorizations.Service
	subscriptionsService  *subscriptions.Service
	operationsService     *operations.Service
}

const (
	contentTypeApplicationJSON = "application/json"
	contentTypeApplicationXML  = "application/xml"
	contentTypeTextXML         = "text/xml"
	contentTypeTextXOPML       = "text/x-opml"
)

// New returns a new initialized server object.
func New(log *logger.Logger, cfg *Config) (*Server, error) {
	db, err := db.New(cfg.DB, log)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}

	crawler := crawler.WithConcurrencyLimit(crawler.New(), 10)
	authorizationsService := authorizations.NewService(db, log)
	usersService := users.NewService(db, cfg.Users)
	operationsService := operations.NewService(log, db, cfg.Operations)
	tagsService := tags.NewService(db)
	itemsService := items.NewService(db, log)
	subscriptionsService := subscriptions.NewService(db, crawler, log, cfg.Subscriptions, itemsService)

	subscriptionsHandler := subscriptions.NewHandler(subscriptionsService, log, operationsService)
	authorizationsHandler := authorizations.NewHandler(usersService, authorizationsService, log, cfg.Authorizations)
	itemsHandler := items.NewHandler(itemsService, log)
	operationsHandler := operations.NewHandler(operationsService, log)
	tagsHandler := tags.NewHandler(tagsService, log)
	usersHandler := users.NewHandler(usersService, log)

	authMiddleware := authorizations.Middleware(authorizationsService, cfg.Authorizations, log)
	requireJSON := middleware.AllowContentType(contentTypeApplicationJSON)
	requireJSONorOPML := middleware.AllowContentType(
		contentTypeApplicationJSON,
		contentTypeApplicationXML,
		contentTypeTextXML,
		contentTypeTextXOPML,
	)

	r := chi.NewRouter()
	r.Use(logger.Middleware(log))
	r.Use(middleware.Recoverer)
	if cfg.HTTP.CORS != nil {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   cfg.HTTP.CORS.AllowedOrigins,
			AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Accept-Encoding"},
			AllowCredentials: true,
		}))
	}
	r.Route("/v1", func(r chi.Router) {
		r.Route("/authorizations", func(r chi.Router) {
			r.With(requireJSON).Post("/", authorizationsHandler.Create())
		})
		r.With(authMiddleware).Route("/subscriptions", func(r chi.Router) {
			r.With(requireJSONorOPML).Post("/", func(w http.ResponseWriter, r *http.Request) {
				switch r.Header.Get("content-type") {
				case contentTypeApplicationJSON:
					subscriptionsHandler.Create()(w, r)
				case contentTypeApplicationXML, contentTypeTextXML, contentTypeTextXOPML:
					subscriptionsHandler.Import()(w, r)
				default:
					w.WriteHeader(http.StatusUnsupportedMediaType)
				}
			})
			r.Get("/", subscriptionsHandler.List())
			r.Route("/{subscriptionId}", func(r chi.Router) {
				r.Route("/items", func(r chi.Router) {
					r.Get("/", func(w http.ResponseWriter, r *http.Request) {
						subscriptionID := chi.URLParam(r, "subscriptionId")
						itemsHandler.List(nil, &subscriptionID)(w, r)
					})
				})
			})
		})
		r.With(authMiddleware).Route("/items", func(r chi.Router) {
			r.Route("/{itemId}", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					itemsHandler.Get(chi.URLParam(r, "itemId"))(w, r)
				})
			})
			r.Get("/", itemsHandler.List(nil, nil))
		})
		r.With(authMiddleware).Route("/operations", func(r chi.Router) {
			r.Route("/{operationId}", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					operationsHandler.Get(chi.URLParam(r, "operationId"))(w, r)
				})
			})
		})
		r.Route("/users", func(r chi.Router) {
			r.With(requireJSON).Post("/", usersHandler.Create())
		})
		r.With(authMiddleware).Route("/tags", func(r chi.Router) {
			r.With(requireJSON).Post("/", tagsHandler.Create())
			r.Get("/", tagsHandler.List())
			r.Route("/{tagId}", func(r chi.Router) {
				r.Route("/items", func(r chi.Router) {
					r.Get("/", func(w http.ResponseWriter, r *http.Request) {
						tagID := chi.URLParam(r, "tagId")
						itemsHandler.List(&tagID, nil)(w, r)
					})
				})
			})
		})
	})

	httpServer, err := httpx.NewServer(cfg.HTTP, log, r)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize http server: %w", err)
	}

	return &Server{
		log:                   log,
		db:                    db,
		httpServer:            httpServer,
		authorizationsService: authorizationsService,
		subscriptionsService:  subscriptionsService,
		operationsService:     operationsService,
	}, nil
}

// Start starts all components of the server.
func (s *Server) Start(ctx context.Context) error {
	if err := db.Migrate(ctx, s.db, s.log); err != nil {
		return fmt.Errorf("failed to apply db migrations: %w", err)
	}
	if err := s.authorizationsService.Init(ctx); err != nil {
		return fmt.Errorf("failed to init jwt service: %w", err)
	}
	g, ctx := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := s.operationsService.Start(ctx); err != nil {
			return fmt.Errorf("failed to start operations service: %w", err)
		}
		return nil
	})
	g.Go(func() error {
		if err := s.subscriptionsService.Start(ctx); err != nil {
			return fmt.Errorf("failed to start subscroptions service: %w", err)
		}
		return nil
	})
	if err := s.httpServer.Start(); err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}
	return g.Wait()
}

// Shutdown gracefully stops all components of the server.
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop http server: %w", err)
	}
	if err := s.operationsService.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown operations: %w", err)
	}
	if err := s.subscriptionsService.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown subscriptions: %w", err)
	}
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close db: %w", err)
	}
	return nil
}
