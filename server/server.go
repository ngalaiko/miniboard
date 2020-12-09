package server

import (
	"context"
	"fmt"

	"github.com/ngalaiko/miniboard/server/db"
	"github.com/ngalaiko/miniboard/server/logger"
)

// Config contains all server configuration.
type Config struct {
	DB *db.Config
}

// Server is the main object.
type Server struct{}

// New returns a new initialized server object.
func New(ctx context.Context, logger *logger.Logger, cfg *Config) (*Server, error) {
	_, err := db.New(ctx, cfg.DB, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize db: %w", err)
	}

	return &Server{}, nil
}
