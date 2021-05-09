package users

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/ngalaiko/miniboard/backend/db"
)

func Test_db__Create(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	user := &User{
		ID:      "test id",
		Hash:    []byte("hash"),
		Created: time.Now().Truncate(time.Nanosecond),
	}
	if err := db.Create(ctx, user); err != nil {
		t.Fatalf("failed to create a user: %s", err)
	}
}

func Test_db__Create_twice(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	user := &User{
		ID:      "test id",
		Hash:    []byte("hash"),
		Created: time.Now().Truncate(time.Nanosecond),
	}
	if err := db.Create(ctx, user); err != nil {
		t.Fatalf("failed to create a user: %s", err)
	}

	if err := db.Create(ctx, user); err == nil {
		t.Fatalf("second create shoud've failed")
	}
}

func Test_db__GetByID_not_found(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	user := &User{
		ID:      "test id",
		Hash:    []byte("hash"),
		Created: time.Now().Truncate(time.Nanosecond),
	}

	fromDB, err := db.GetByID(ctx, user.ID)
	if fromDB != nil {
		t.Fatalf("nothing should be returned, got %+v", fromDB)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected %s, got %s", sql.ErrNoRows, err)
	}
}

func Test_db__GetByID(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	user := &User{
		ID:      "test id",
		Hash:    []byte("hash"),
		Created: time.Now().UTC().Truncate(time.Nanosecond),
	}
	if err := db.Create(ctx, user); err != nil {
		t.Fatalf("failed to create a user: %s", err)
	}

	fromDB, err := db.GetByID(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to get user from the db: %s", err)
	}
	if !reflect.DeepEqual(user, fromDB) {
		t.Fatalf("expected %+v, got %+v", user, fromDB)
	}
}

func Test_db__GetByUsername_not_found(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	user := &User{
		ID:       "test id",
		Username: "username",
		Hash:     []byte("hash"),
		Created:  time.Now().Truncate(time.Nanosecond),
	}

	fromDB, err := db.GetByUsername(ctx, user.Username)
	if fromDB != nil {
		t.Fatalf("nothing should be returned, got %+v", fromDB)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected %s, got %s", sql.ErrNoRows, err)
	}
}

func Test_db__GetByUsername(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	user := &User{
		ID:       "test id",
		Username: "username",
		Hash:     []byte("hash"),
		Created:  time.Now().UTC().Truncate(time.Nanosecond),
	}
	if err := db.Create(ctx, user); err != nil {
		t.Fatalf("failed to create a user: %s", err)
	}

	fromDB, err := db.GetByUsername(ctx, user.Username)
	if err != nil {
		t.Fatalf("failed to get user from the db: %s", err)
	}
	if !reflect.DeepEqual(user, fromDB) {
		t.Fatalf("expected %+v, got %+v", user, fromDB)
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

func (tl *testLogger) Error(string, ...interface{}) {}
