package operations

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

type database struct {
	db *sql.DB
}

func newDatabase(sqldb *sql.DB) *database {
	return &database{
		db: sqldb,
	}
}

// Create adds operation to the database.
func (db *database) Create(ctx context.Context, operation *Operation) error {
	_, createErr := db.db.ExecContext(ctx, `
	INSERT INTO operations (
		id, user_id, done
	) VALUES (
		$1, $2, $3
	)
	`,
		operation.ID,
		operation.UserID,
		operation.Done,
	)
	return createErr
}

// Update updates operation in the database.
func (db *database) Update(ctx context.Context, operation *Operation) error {
	var errorMessage *string
	var response *[]byte
	if operation.Result != nil {
		bytes, err := json.Marshal(operation.Result.Response)
		if err != nil {
			return fmt.Errorf("failed to marshal response to json: %w", err)
		}
		response = &bytes
	}

	if operation.Result != nil && operation.Result.Error != nil {
		errorMessage = &operation.Result.Error.Message
	}

	result, updateErr := db.db.ExecContext(ctx, `
	UPDATE operations
	SET
		done = $1,
		error = $2,
		response = $3
	WHERE
		id = $4
		AND user_id = $5
	`,
		operation.Done,
		errorMessage,
		response,
		operation.ID,
		operation.UserID,
	)
	if updateErr != nil {
		return updateErr
	}

	n, err := result.RowsAffected()
	switch {
	case err != nil:
		return fmt.Errorf("failed to get rows affected: %w", err)
	case n == 0:
		return sql.ErrNoRows
	default:
		return nil
	}
}

// Get returns operation by id.
func (db *database) Get(ctx context.Context, id string, userID string) (*Operation, error) {
	row := db.db.QueryRowContext(ctx, `
	SELECT
		id,
		user_id,
		done,
		error,
		response
	FROM
		operations
	WHERE
		id = $1
		AND user_id = $2
	`, id, userID)

	return db.scanRow(row)
}

type scannable interface {
	Scan(...interface{}) error
}

func (db *database) scanRow(row scannable) (*Operation, error) {
	var errorMessage *string
	var response *[]byte

	operation := &Operation{}

	err := row.Scan(
		&operation.ID,
		&operation.UserID,
		&operation.Done,
		&errorMessage,
		&response,
	)

	if err != nil {
		return nil, err
	}

	if response != nil {
		operation.Result = &Result{
			Response: map[string]interface{}{},
		}
		if err := json.Unmarshal(*response, &operation.Result.Response); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response (%s): %w", string(*response), err)
		}
	}

	if errorMessage != nil {
		operation.Result = &Result{
			Error: &Error{
				Message: *errorMessage,
			},
		}
	}

	return operation, nil
}
