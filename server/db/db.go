package db

import (
	"context"
	"database/sql"
	"fmt"

	// supported drivers
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// Config contains database configuration.
type Config struct {
	Driver string
	Addr   string
}

// New returns a database instance.
func New(ctx context.Context, cfg *Config) (*sql.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := migrate(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return db, nil
}
