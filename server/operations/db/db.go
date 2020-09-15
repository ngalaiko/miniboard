package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/genproto/googleapis/longrunning"
	status "google.golang.org/genproto/googleapis/rpc/status"
)

// DB provides database access to operations.
type DB struct {
	db *sql.DB
}

// New returns new DB.
func New(sqldb *sql.DB) *DB {
	return &DB{
		db: sqldb,
	}
}

// Create adds operation to the database.
func (db *DB) Create(ctx context.Context, userID string, operation *longrunning.Operation) error {
	metadata, err := proto.Marshal(operation.Metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	_, createErr := db.db.ExecContext(ctx, `
	INSERT INTO operations (
		id, user_id, metadata
	) VALUES (
		$1, $2, $3
	)
	`,
		strings.Replace(operation.Name, "operations/", "", 1),
		userID,
		metadata,
	)
	return createErr
}

// Update updates operation in the database.
func (db *DB) Update(ctx context.Context, userID string, operation *longrunning.Operation) error {
	var operationError, operationResponse []byte

	if e := operation.GetError(); e != nil {
		var err error
		operationError, err = proto.Marshal(e)
		if err != nil {
			return fmt.Errorf("failed to marshal error: %w", err)
		}
	}

	if r := operation.GetResponse(); r != nil {
		var err error
		operationResponse, err = proto.Marshal(r)
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}
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
		operationError,
		operationResponse,
		strings.Replace(operation.Name, "operations/", "", 1),
		userID,
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
func (db *DB) Get(ctx context.Context, id string, userID string) (*longrunning.Operation, error) {
	row := db.db.QueryRowContext(ctx, `
	SELECT
		id,
		done,
		error,
		response,
		metadata
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

type dbOperation struct {
	Operation *longrunning.Operation
	ID        string
	Error     []byte
	Response  []byte
	Metadata  []byte
}

func (db *DB) scanRow(row scannable) (*longrunning.Operation, error) {
	operation := &dbOperation{
		Operation: &longrunning.Operation{},
	}
	err := row.Scan(
		&operation.ID,
		&operation.Operation.Done,
		&operation.Error,
		&operation.Response,
		&operation.Metadata,
	)

	if err != nil {
		return nil, err
	}

	operation.Operation.Name = fmt.Sprintf("operations/%s", operation.ID)

	if operation.Error != nil {
		status := &status.Status{}
		if err := proto.Unmarshal(operation.Error, status); err != nil {
			return nil, fmt.Errorf("failed to unmarshal error: %w", err)
		}
		operation.Operation.Result = &longrunning.Operation_Error{
			Error: status,
		}
	}

	if operation.Response != nil {
		operationResponse := &any.Any{}
		if err := proto.Unmarshal(operation.Response, operationResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
		operation.Operation.Result = &longrunning.Operation_Response{
			Response: operationResponse,
		}
	}

	if operation.Metadata != nil {
		operation.Operation.Metadata = &any.Any{}
		if err := proto.Unmarshal(operation.Metadata, operation.Operation.Metadata); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return operation.Operation, nil
}
