package db

import (
	"database/sql"
	"fmt"

	// postgres driver
	_ "github.com/lib/pq"

	"github.com/sirupsen/logrus"
)

// NewPostgres creates postgres backed database.
func NewPostgres(connectStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	log().Infof("connected to postgresql")

	return db, nil
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "postgres",
	})
}
