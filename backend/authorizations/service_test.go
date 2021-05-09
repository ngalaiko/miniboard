package authorizations

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/ngalaiko/miniboard/backend/db"
)

func Test_Init(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)

	service := NewService(sqldb, &testLogger{})

	if err := service.Init(ctx); err != nil {
		t.Fatalf("unexpected error")
	}
}

func Test_Init__twice(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)

	service := NewService(sqldb, &testLogger{})

	if err := service.Init(ctx); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if err := service.Init(ctx); err == nil {
		t.Fatalf("error expexted, got nil")
	}
}

func Test_NewToken(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	service := NewService(createTestDB(ctx, t), &testLogger{})

	token, err := service.NewToken(ctx, "user id")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if token.UserID != "user id" {
		t.Errorf("expected user id %s, got %s", "user id", token.UserID)
	}

	if token.Token == "" {
		t.Errorf("token is empty")
	}

	if token.ExpiresAt.String() == "" {
		t.Errorf("token expiry is empty")
	}
}

func Test_Verify__invalid(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	service := NewService(createTestDB(ctx, t), &testLogger{})

	if _, err := service.Verify(ctx, "invalid"); !errors.Is(err, ErrInvalidToken) {
		t.Errorf("exected %s, got %s", ErrInvalidToken, err)
	}
}

func Test_Verify(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	service := NewService(createTestDB(ctx, t), &testLogger{})

	token, err := service.NewToken(ctx, "user id")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	parsed, err := service.Verify(ctx, token.Token)
	if err != nil {
		t.Errorf("exected nil, got %s", err)
	}

	if !reflect.DeepEqual(parsed, token) {
		t.Errorf("tokens must be the equal")
	}
}

func createTestDB(ctx context.Context, t *testing.T) *sql.DB {
	sqlite, err := db.New(&db.Config{
		Addr: "sqlite3://file::memory::",
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
