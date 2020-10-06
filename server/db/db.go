package db

import (
	"context"
	"database/sql"
	"fmt"

	// supported drivers
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/sirupsen/logrus"
)

// New returns a database instance.
func New(ctx context.Context, driver string, address string) (*sql.DB, error) {
	db, err := sql.Open(driver, address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := migrate(ctx, db); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return db, nil
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "db",
	})
}
