package server

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ngalaiko/miniboard/server/db"
	"github.com/ngalaiko/miniboard/server/http"
	"github.com/ngalaiko/miniboard/server/jwt"
	"github.com/ngalaiko/miniboard/server/logger"
	"github.com/ngalaiko/miniboard/server/users"
)

// Config contains all server configuration.
type Config struct {
	DB   *db.Config
	HTTP *http.Config
}

// Server is the main object.
type Server struct {
	logger     *logger.Logger
	db         *sql.DB
	httpServer *http.Server
	jwtService *jwt.Service
}

// New returns a new initialized server object.
func New(logger *logger.Logger, cfg *Config) (*Server, error) {
	db, err := db.New(cfg.DB, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}

	_ = users.NewService(db)

	jwtService := jwt.NewService(db, logger)

	httpServer, err := http.NewServer(cfg.HTTP, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize http server: %w", err)
	}

	return &Server{
		logger:     logger,
		db:         db,
		httpServer: httpServer,
		jwtService: jwtService,
	}, nil
}

// Start starts all components of the server.
func (s *Server) Start(ctx context.Context) error {
	if err := db.Migrate(ctx, s.db, s.logger); err != nil {
		return fmt.Errorf("failed to apply db migrations: %w", err)
	}
	if err := s.jwtService.Init(ctx); err != nil {
		return fmt.Errorf("failed to init jwt service: %w", err)
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
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("failed to close db: %w", err)
	}
	return nil
}
