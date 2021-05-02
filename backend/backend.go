package backend

import (
	"context"
	"database/sql"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/crawler"
	"github.com/ngalaiko/miniboard/backend/db"
	"github.com/ngalaiko/miniboard/backend/httpx"
	"github.com/ngalaiko/miniboard/backend/items"
	"github.com/ngalaiko/miniboard/backend/logger"
	"github.com/ngalaiko/miniboard/backend/subscriptions"
	"github.com/ngalaiko/miniboard/backend/tags"
	"github.com/ngalaiko/miniboard/backend/users"
	"github.com/ngalaiko/miniboard/backend/web"
)

// Config contains all server configuration.
type Config struct {
	DB            *db.Config            `yaml:"db"`
	HTTP          *httpx.Config         `yaml:"http"`
	Subscriptions *subscriptions.Config `yaml:"subscriptions"`
	Users         *users.Config         `yaml:"users"`
	Web           *web.Config           `yaml:"web"`
}

// Server is the main object.
type Server struct {
	log                   *logger.Logger
	db                    *sql.DB
	httpServer            *httpx.Server
	authorizationsService *authorizations.Service
	subscriptionsService  *subscriptions.Service
}

// New returns a new initialized server object.
func New(log *logger.Logger, cfg *Config) (*Server, error) {
	db, err := db.New(cfg.DB, log)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}

	crawler := crawler.WithConcurrencyLimit(crawler.New(), 10)
	authorizationsService := authorizations.NewService(db, log)
	usersService := users.NewService(db, cfg.Users)
	tagsService := tags.NewService(db)
	itemsService := items.NewService(db, log)
	subscriptionsService := subscriptions.NewService(db, crawler, log, cfg.Subscriptions, itemsService)

	webHandler := web.NewHandler(cfg.Web, log, itemsService, tagsService, subscriptionsService, usersService, authorizationsService)

	httpServer, err := httpx.NewServer(cfg.HTTP, log, webHandler)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize http server: %w", err)
	}

	return &Server{
		log:                   log,
		db:                    db,
		httpServer:            httpServer,
		authorizationsService: authorizationsService,
		subscriptionsService:  subscriptionsService,
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
	if err := s.subscriptionsService.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown subscriptions: %w", err)
	}
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close db: %w", err)
	}
	return nil
}
