package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Migrate prepares database scheme.
func Migrate(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	if err := initMigrations(ctx, tx); err != nil {
		return fmt.Errorf("error while initializing migrations: %w", err)
	}

	applied, err := getAppliedMigrations(ctx, tx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	for _, m := range migrations() {
		if _, ok := applied[m.Name]; ok {
			continue
		}
		if _, err := tx.ExecContext(ctx, m.Query); err != nil {
			return fmt.Errorf("failed to apply %s: %w", m.Name, err)
		}
		if err := storeAppliedMigration(ctx, tx, m); err != nil {
			return fmt.Errorf("failed to store %s: %w", m.Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func initMigrations(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS migrations (
		name       TEXT      NOT NULL,
		PRIMARY KEY (name)
	)
	`); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	return nil
}

func getAppliedMigrations(ctx context.Context, tx *sql.Tx) (map[string]*migration, error) {
	rows, err := tx.QueryContext(ctx, `
	SELECT name
	FROM migrations
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mm := map[string]*migration{}
	for rows.Next() {
		m := new(migration)
		if err := rows.Scan(&m.Name); err != nil {
			return nil, err
		}
		mm[m.Name] = m
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return mm, nil
}

func storeAppliedMigration(ctx context.Context, tx *sql.Tx, m *migration) error {
	_, err := tx.ExecContext(ctx, `
	INSERT INTO migrations (
		name
	)
	VALUES (
		$1
	)
	`, m.Name)
	return err
}
