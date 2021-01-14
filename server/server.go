package server

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ngalaiko/miniboard/server/authorizations"
	"github.com/ngalaiko/miniboard/server/crawler"
	"github.com/ngalaiko/miniboard/server/db"
	"github.com/ngalaiko/miniboard/server/feeds"
	"github.com/ngalaiko/miniboard/server/httpx"
	"github.com/ngalaiko/miniboard/server/logger"
	"github.com/ngalaiko/miniboard/server/operations"
	"github.com/ngalaiko/miniboard/server/users"
)

// Config contains all server configuration.
type Config struct {
	DB         *db.Config         `yaml:"db"`
	HTTP       *httpx.Config      `yaml:"http"`
	Operations *operations.Config `yaml:"operations"`
	Users      *users.Config      `yaml:"users"`
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
	feedsService := feeds.NewService(db, crawler.New())

	withAuth := authorizations.Authenticate(authorizationsService, logger)

	httpServer.Route("/v1/authorizations", httpx.Chain(
		authorizations.NewHandler(usersService, authorizationsService, logger),
		httpx.WithCors(),
	))
	httpServer.Route("/v1/feeds", httpx.Chain(
		feeds.NewHandler(feedsService, logger, operationsService),
		httpx.WithCors(),
		withAuth,
	))
	httpServer.Route("/v1/operations", httpx.Chain(
		operations.NewHandler(operationsService, logger),
		httpx.WithCors(),
		withAuth,
	))
	httpServer.Route("/v1/users", httpx.Chain(
		users.NewHandler(usersService, logger),
		httpx.WithCors(),
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
