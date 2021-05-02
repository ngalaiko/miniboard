package keys

import (
	"context"
	"database/sql"
)

// Database is a persistent storage for jwt keys.
type Database struct {
	db *sql.DB
}

// NewDatabase crates new DB.
func NewDatabase(db *sql.DB) *Database {
	return &Database{
		db: db,
	}
}

// Create adds a key to the db.
func (db *Database) Create(ctx context.Context, key *Key) error {
	_, err := db.db.ExecContext(ctx, `
	INSERT INTO jwt_keys (
		id,
		public_der
	) VALUES (
		$1, $2
	)
	`, key.ID, key.PublicDER)
	return err
}

// Get returns a key by the given id.
func (db *Database) Get(ctx context.Context, id string) (*Key, error) {
	row := db.db.QueryRowContext(ctx, `
	SELECT
		id, public_der
	FROM
		jwt_keys
	WHERE
		id = $1
	`, id)
	return db.scanRow(row)
}

// Delete deletes key by id.
func (db *Database) Delete(ctx context.Context, id string) error {
	_, err := db.db.ExecContext(ctx, `
	DELETE
	FROM
		jwt_keys
	WHERE
		id = $1
	`, id)
	return err
}

// List returns all keys.
func (db *Database) List(ctx context.Context) ([]*Key, error) {
	rows, err := db.db.QueryContext(ctx, `
	SELECT
		id, public_der
	FROM
		jwt_keys
	ORDER BY
		id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	keys := []*Key{}
	for rows.Next() {
		key, err := db.scanRow(rows)
		if err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return keys, nil
}

type scannable interface {
	Scan(...interface{}) error
}

func (db *Database) scanRow(row scannable) (*Key, error) {
	key := &Key{}
	if err := row.Scan(&key.ID, &key.PublicDER); err != nil {
		return nil, err
	}
	return key, nil
}
