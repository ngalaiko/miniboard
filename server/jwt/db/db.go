package db

import (
	"context"
	"database/sql"
)

// PublicKey is a public key.
type PublicKey struct {
	ID        string
	DerBase64 string
}

// DB is a persistent storage for jwt keys.
type DB struct {
	db *sql.DB
}

// New crates new DB.
func New(db *sql.DB) *DB {
	return &DB{
		db: db,
	}
}

// Create adds a key to the db.
func (db *DB) Create(ctx context.Context, key *PublicKey) error {
	_, err := db.db.ExecContext(ctx, `
	INSERT INTO public_keys (
		id,
		der_base64
	) VALUES (
		$1, $2
	)
	`, key.ID, key.DerBase64)
	return err
}

// Get returns key by id.
func (db *DB) Get(ctx context.Context, id string) (*PublicKey, error) {
	row := db.db.QueryRowContext(ctx, `
	SELECT id, der_base64
	FROM public_keys
	WHERE id = $1
	`, id)
	key := &PublicKey{}
	err := row.Scan(&key.ID, &key.DerBase64)
	return key, err
}

// Delete deletes key by id.
func (db *DB) Delete(ctx context.Context, id string) error {
	_, err := db.db.ExecContext(ctx, `
	DELETE FROM public_keys
	WHERE id = $1
	`, id)
	return err
}

// List returns all keys.
func (db *DB) List(ctx context.Context) ([]*PublicKey, error) {
	rows, err := db.db.QueryContext(ctx, `
	SELECT id, der_base64
	FROM public_keys
	ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	keys := []*PublicKey{}
	for rows.Next() {
		key := &PublicKey{}
		if err := rows.Scan(&key.ID, &key.DerBase64); err != nil {
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
