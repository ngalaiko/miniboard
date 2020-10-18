package server

import (
	"context"
	"fmt"

	"github.com/ngalaiko/miniboard/server/api"
	"github.com/ngalaiko/miniboard/server/db"
	"github.com/ngalaiko/miniboard/server/email"
	"github.com/ngalaiko/miniboard/server/fetch"
	"github.com/ngalaiko/miniboard/server/jwt"
)

// Server is the api server.
type Server struct {
	httpAPI *api.HTTP
}

// Config contains server configuration.
type Config struct {
	HTTP *api.HTTPConfig
	DB   *db.Config
}

// New creates new api server.
func New(
	ctx context.Context,
	cfg *Config,
	emailClient email.Client,
) (*Server, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	sqldb, err := db.New(ctx, cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	fetcher := fetch.New()
	jwtService := jwt.NewService(ctx, sqldb)

	httpAPI, err := api.NewHTTP(ctx, cfg.HTTP, sqldb, fetcher, emailClient, jwtService)
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
