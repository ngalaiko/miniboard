package db

import (
	"database/sql"
	"fmt"
	"net/url"

	// supported drivers
	_ "github.com/mattn/go-sqlite3"
)

type logger interface {
	Debug(string, ...interface{})
}

// Config contains database configuration.
type Config struct {
	Addr               string `yaml:"addr"`
	MaxOpenConnections int    `yaml:"max_open_connections"`
}

// New returns a database instance.
func New(cfg *Config, logger logger) (*sql.DB, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	if cfg.Addr == "" {
		return nil, fmt.Errorf("addr is not set")
	}

	u, err := url.Parse(cfg.Addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse addr: %w", err)
	}

	addr := fmt.Sprintf("%s%s?%s", u.Hostname(), u.Path, u.RawQuery)
	db, err := sql.Open(u.Scheme, addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	db.SetMaxOpenConns(cfg.MaxOpenConnections)

	return db, nil
}
