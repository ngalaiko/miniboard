package operations

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ngalaiko/miniboard/backend/db"
)

func Test_db__Create(t *testing.T) {
	ctx := context.TODO()
	database := newDatabase(testDB(ctx, t))
	if err := database.Create(ctx, New("user")); err != nil {
		t.Fatalf("failed to create operation: %s", err)
	}
}

func Test_db__Create_twice(t *testing.T) {
	ctx := context.TODO()
	database := newDatabase(testDB(ctx, t))
	operation := New("user")
	if err := database.Create(ctx, operation); err != nil {
		t.Fatalf("failed to create operation: %s", err)
	}
	if err := database.Create(ctx, operation); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func Test_db__Get(t *testing.T) {
	ctx := context.TODO()
	database := newDatabase(testDB(ctx, t))
	operation := New("user")
	if err := database.Create(ctx, operation); err != nil {
		t.Fatalf("failed to create operation: %s", err)
	}

	fromDB, err := database.Get(ctx, operation.ID, operation.UserID)
	if err != nil {
		t.Fatalf("failed to get: %s", err)
	}

	if !reflect.DeepEqual(operation, fromDB) {
		t.Fatalf("expected %+v, got %+v", operation, fromDB)
	}
}

func Test_db__Get_not_exists(t *testing.T) {
	ctx := context.TODO()
	database := newDatabase(testDB(ctx, t))
	operation := New("user")

	_, err := database.Get(ctx, operation.ID, operation.UserID)
	if err != sql.ErrNoRows {
		t.Fatalf("expected %s, got %s", sql.ErrNoRows, err)
	}
}

func Test_db__Update_response(t *testing.T) {
	ctx := context.TODO()
	database := newDatabase(testDB(ctx, t))
	operation := New("user")
	if err := database.Create(ctx, operation); err != nil {
		t.Fatalf("failed to create operation: %s", err)
	}

	operation.Success(map[string]string{"key": "value"})

	if err := database.Update(ctx, operation); err != nil {
		t.Fatalf("expected nil, got %s", err)
	}

	fromDB, err := database.Get(ctx, operation.ID, operation.UserID)
	if err != nil {
		t.Fatalf("failed to get: %s", err)
	}

	operationMarshalled, err := json.Marshal(operation)
	if err != nil {
		t.Fatal("failed to marshal operation: %w", err)
	}

	fromDBMarshalled, err := json.Marshal(fromDB)
	if err != nil {
		t.Fatal("failed to marshal fromDB: %w", err)
	}

	if !bytes.Equal(operationMarshalled, fromDBMarshalled) {
		t.Fatalf("expected %+v, got %+v", string(operationMarshalled), string(fromDBMarshalled))
	}
}

func Test_db__Update_error(t *testing.T) {
	ctx := context.TODO()
	database := newDatabase(testDB(ctx, t))
	operation := New("user")
	if err := database.Create(ctx, operation); err != nil {
		t.Fatalf("failed to create operation: %s", err)
	}

	operation.Error("failed")

	if err := database.Update(ctx, operation); err != nil {
		t.Fatalf("expected nil, got %s", err)
	}

	fromDB, err := database.Get(ctx, operation.ID, operation.UserID)
	if err != nil {
		t.Fatalf("failed to get: %s", err)
	}

	if !reflect.DeepEqual(operation, fromDB) {
		t.Fatalf("expected %+v, got %+v", operation, fromDB)
	}
}

func Test_db__Update_not_exists(t *testing.T) {
	ctx := context.TODO()
	database := newDatabase(testDB(ctx, t))
	operation := New("user")

	if err := database.Update(ctx, operation); err != sql.ErrNoRows {
		t.Fatalf("expected %s, got %s", sql.ErrNoRows, err)
	}
}

func testDB(ctx context.Context, t *testing.T) *sql.DB {
	sqlite, err := db.New(&db.Config{
		Driver: "sqlite3",
		Addr:   "file::memory:",
	}, &testLogger{})
	if err != nil {
		t.Fatalf("failed to create database: %s", err)
	}

	if err := db.Migrate(ctx, sqlite, &testLogger{}); err != nil {
		t.Fatalf("failed to run migrations: %s", err)
	}

	return sqlite
}

//

type testLogger struct{}

func (tl *testLogger) Debug(string, ...interface{}) {}

func (tl *testLogger) Info(string, ...interface{}) {}

func (tl *testLogger) Error(string, ...interface{}) {}
