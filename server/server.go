package server

import (
	"context"
	"fmt"

	"github.com/ngalaiko/miniboard/server/api"
	"github.com/ngalaiko/miniboard/server/db"
	"github.com/ngalaiko/miniboard/server/email"
	"github.com/ngalaiko/miniboard/server/fetch"
	"github.com/ngalaiko/miniboard/server/jwt"
	"github.com/ngalaiko/miniboard/server/logger"
)

// Server is the api server.
type Server struct {
	httpAPI *api.HTTP
}

// Config contains server configuration.
type Config struct {
	HTTP  *api.HTTPConfig
	DB    *db.Config
	Email *email.Config
}

// New creates new api server.
func New(
	ctx context.Context,
	cfg *Config,
) (*Server, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	logger := logger.New()

	sqldb, err := db.New(ctx, cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	emailClient, err := email.New(cfg.Email, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create email client: %w", err)
	}

	fetcher := fetch.New()

	jwtService, err := jwt.NewService(ctx, logger, sqldb)
	if err != nil {
		return nil, fmt.Errorf("faied to create jwt service: %w", err)
	}

	httpAPI, err := api.NewHTTP(ctx, cfg.HTTP, logger, sqldb, fetcher, emailClient, jwtService)
	if err != nil {
		return nil, fmt.Errorf("failed to create http api: %w", err)
	}

	return &Server{
		httpAPI: httpAPI,
	}, nil
}

// Serve starts the server.
func (s *Server) Serve(ctx context.Context) error {
	if err := s.httpAPI.ListenAndServe(ctx); err != nil {
		return fmt.Errorf("failed to start http api: %w", err)
	}
	return nil
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.httpAPI.Shutdown(ctx); err != nil {
		return fmt.Errorf("error stopping http api: %w", err)
	}

	return nil
}
