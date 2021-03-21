package users

import (
	"context"
	"database/sql"
	"time"
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
		id,
		username,
		hash,
		created_epoch
	) VALUES (
		$1, $2, $3, $4
	)`, user.ID, user.Username, user.Hash, user.Created.UTC().UnixNano())

	return err
}

// GetByID returns a user from the database with the given id.
func (d *database) GetByID(ctx context.Context, id string) (*User, error) {
	row := d.db.QueryRowContext(ctx, `
	SELECT
		id,
		username,
		hash,
		created_epoch
	FROM
		users
	WHERE
		id = $1
	`, id)

	return d.scanRow(row)
}

// GetByUsername returns a user from the database with the given username.
func (d *database) GetByUsername(ctx context.Context, username string) (*User, error) {
	row := d.db.QueryRowContext(ctx, `
	SELECT
		id,
		username,
		hash,
		created_epoch
	FROM
		users
	WHERE
		username = $1
	`, username)

	return d.scanRow(row)
}

type scannable interface {
	Scan(...interface{}) error
}

func (d *database) scanRow(row scannable) (*User, error) {
	user := &User{}
	var createdEpoch int64
	if err := row.Scan(&user.ID, &user.Username, &user.Hash, &createdEpoch); err != nil {
		return nil, err
	}

	user.Created = time.Unix(0, createdEpoch).UTC().Round(time.Nanosecond)

	return user, nil
}
