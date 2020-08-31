package db

import (
	"database/sql"

	// sqllite driver
	_ "github.com/mattn/go-sqlite3"
)

// NewSQLite creates sqlite backed database.
func NewSQLite(name string) (*sql.DB, error) {
	return sql.Open("sqlite3", name)
}
