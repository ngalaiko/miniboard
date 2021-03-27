package db

import (
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
	Driver             string `yaml:"driver"`
	Addr               string `yaml:"addr"`
	MaxOpenConnections int    `yaml:"max_open_connections"`
}

// New returns a database instance.
func New(cfg *Config, logger logger) (*sql.DB, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	if cfg.Driver == "" {
		return nil, fmt.Errorf("driver is not set")
	}

	if _, supported := map[string]bool{
		"sqlite3":  true,
		"postgres": true,
	}[cfg.Driver]; !supported {
		return nil, fmt.Errorf("'%s' is not supported", cfg.Driver)
	}

	if cfg.Addr == "" {
		return nil, fmt.Errorf("addr is not set")
	}

	db, err := sql.Open(cfg.Driver, cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	db.SetMaxOpenConns(cfg.MaxOpenConnections)

	return db, nil
}
