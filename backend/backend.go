package backend

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/crawler"
	"github.com/ngalaiko/miniboard/backend/db"
	"github.com/ngalaiko/miniboard/backend/httpx"
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
	Users          *users.Config          `yaml:"users"`
}

// Server is the main object.
type Server struct {
	logger                *logger.Logger
	db                    *sql.DB
	httpServer            *httpx.Server
	authorizationsService *authorizations.Service
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

	authorizationsService := authorizations.NewService(db, logger)
	usersService := users.NewService(db, cfg.Users)
	operationsService := operations.NewService(logger, db, cfg.Operations)
	tagsService := tags.NewService(db)
	subscriptionsService := subscriptions.NewService(db, crawler.New(), logger)

	withAuth := authorizations.Authenticate(authorizationsService, cfg.Authorizations, logger)

	corsDomains := []string{}
	if cfg.Cors != nil {
		corsDomains = append(corsDomains, cfg.Cors.Domains...)
	}

	httpServer.Route("/v1/authorizations", httpx.Chain(
		authorizations.NewHandler(usersService, authorizationsService, logger, cfg.Authorizations),
		httpx.WithCors(corsDomains...),
	))
	httpServer.Route("/v1/subscriptions", httpx.Chain(
		subscriptions.NewHandler(subscriptionsService, logger, operationsService),
		httpx.WithCors(corsDomains...),
		withAuth,
	))
	httpServer.Route("/v1/operations", httpx.Chain(
		operations.NewHandler(operationsService, logger),
		httpx.WithCors(corsDomains...),
		withAuth,
	))
	httpServer.Route("/v1/tags", httpx.Chain(
		tags.NewHandler(tagsService, logger),
		httpx.WithCors(corsDomains...),
		withAuth,
	))
	httpServer.Route("/v1/users", httpx.Chain(
		users.NewHandler(usersService, logger),
		httpx.WithCors(corsDomains...),
	))

	return &Server{
		logger:                logger,
		db:                    db,
		httpServer:            httpServer,
		authorizationsService: authorizationsService,
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
	if err := s.operationsService.Start(ctx); err != nil {
		return fmt.Errorf("failed to start operations service: %w", err)
	}
	if err := s.httpServer.Start(); err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}
	return nil
}

// Shutdown gracefully stops all components of the server.
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop http server: %w", err)
	}
	s.operationsService.Shutdown(ctx)
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close db: %w", err)
	}
	return nil
}
