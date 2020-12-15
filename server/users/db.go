package users

import (
	"context"
	"database/sql"
)

type database struct {
	db *sql.DB
}

func newDB(sqldb *sql.DB) *database {
	return &database{
		db: sqldb,
	}
}

// Create creates a user in the database.
func (d *database) Create(ctx context.Context, user *User) error {
	_, err := d.db.ExecContext(ctx, `
	INSERT INTO users (
		id
	) VALUES (
		$1
	)`, user.ID)

	return err
}

// Get returns a user from the database.
func (d *database) Get(ctx context.Context, id string) (*User, error) {
	row := d.db.QueryRowContext(ctx, `
	SELECT
		id
	FROM
		users
	WHERE
		id = $1
	`, id)

	return d.scanRow(row)
}

type scannable interface {
	Scan(...interface{}) error
}

func (d *database) scanRow(row scannable) (*User, error) {
	user := &User{}
	if err := row.Scan(&user.ID); err != nil {
		return nil, err
	}

	return user, nil
}
