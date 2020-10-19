package db

import (
	"context"
	"database/sql"
	"fmt"

	// supported drivers
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type logger interface {
	Info(string, ...interface{})
}

// Config contains database configuration.
type Config struct {
	Driver string
	Addr   string
}

// New returns a database instance.
func New(ctx context.Context, cfg *Config, logger logger) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	logger.Info("connected to %s", cfg.Driver)

	if err := migrate(ctx, db, logger); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return db, nil
}