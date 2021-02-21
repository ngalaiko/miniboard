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
	"github.com/ngalaiko/miniboard/backend/operations"
	"github.com/ngalaiko/miniboard/backend/subscriptions"
	"github.com/ngalaiko/miniboard/backend/tags"
	"github.com/ngalaiko/miniboard/backend/users"
)

type corsConfig struct {
	Domains []string `yaml:"domains"`
}

// Config contains all server configuration.
type Config struct {
	Authorizations *authorizations.Config `yaml:"authorizations"`
	DB             *db.Config             `yaml:"db"`
	HTTP           *httpx.Config          `yaml:"http"`
	Cors           *corsConfig            `yaml:"cors"`
	Operations     *operations.Config     `yaml:"operations"`
	Subscriptions  *subscriptions.Config  `yaml:"subscriptions"`
	Users          *users.Config          `yaml:"users"`
}

// Server is the main object.
type Server struct {
	logger                *logger.Logger
	db                    *sql.DB
	httpServer            *httpx.Server
	authorizationsService *authorizations.Service
	subscriptionsService  *subscriptions.Service
	operationsService     *operations.Service
}

// New returns a new initialized server object.
func New(logger *logger.Logger, cfg *Config) (*Server, error) {
	db, err := db.New(cfg.DB, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}

	httpServer, err := httpx.NewServer(cfg.HTTP, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize http server: %w", err)
	}

	crawler := crawler.WithConcurrencyLimit(crawler.New(), 10)
	authorizationsService := authorizations.NewService(db, logger)
	usersService := users.NewService(db, cfg.Users)
	operationsService := operations.NewService(logger, db, cfg.Operations)
	tagsService := tags.NewService(db)
	itemsService := items.NewService(db, logger)
	subscriptionsService := subscriptions.NewService(db, crawler, logger, cfg.Subscriptions, itemsService)

	withAuth := authorizations.Authenticate(authorizationsService, cfg.Authorizations, logger)

	corsDomains := []string{}
	if cfg.Cors != nil {
		corsDomains = append(corsDomains, cfg.Cors.Domains...)
	}
	withCORS := httpx.WithCors(corsDomains...)

	httpServer.Route("/v1/authorizations", httpx.Chain(
		authorizations.NewHandler(usersService, authorizationsService, logger, cfg.Authorizations),
		withCORS,
	))
	httpServer.Route("/v1/subscriptions", httpx.Chain(
		subscriptions.NewHandler(subscriptionsService, logger, operationsService),
		withCORS,
		withAuth,
	))
	httpServer.Route("/v1/operations", httpx.Chain(
		operations.NewHandler(operationsService, logger),
		withCORS,
		withAuth,
	))
	httpServer.Route("/v1/items", httpx.Chain(
		items.NewHandler(itemsService, logger),
		withCORS,
		withAuth,
	))
	httpServer.Route("/v1/tags", httpx.Chain(
		tags.NewHandler(tagsService, logger),
		withCORS,
		withAuth,
	))
	httpServer.Route("/v1/users", httpx.Chain(
		users.NewHandler(usersService, logger),
		withCORS,
	))

	return &Server{
		logger:                logger,
		db:                    db,
		httpServer:            httpServer,
		authorizationsService: authorizationsService,
		subscriptionsService:  subscriptionsService,
		operationsService:     operationsService,
	}, nil
}

// Start starts all components of the server.
func (s *Server) Start(ctx context.Context) error {
	if err := db.Migrate(ctx, s.db, s.logger); err != nil {
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
	g.Go(func() error {
		if err := s.httpServer.Start(); err != nil {
			return fmt.Errorf("failed to start http server: %w", err)
		}
		return nil
	})
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
