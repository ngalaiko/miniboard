package users

import (
	"context"
	"database/sql"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/ngalaiko/miniboard/backend/db"
)

func Test_db__Create(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	user := &User{
		ID:   "test id",
		Hash: []byte("hash"),
	}
	if err := db.Create(ctx, user); err != nil {
		t.Fatalf("failed to create a user: %s", err)
	}
}

func Test_db__Create_twice(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	user := &User{
		ID:   "test id",
		Hash: []byte("hash"),
	}
	if err := db.Create(ctx, user); err != nil {
		t.Fatalf("failed to create a user: %s", err)
	}

	if err := db.Create(ctx, user); err == nil {
		t.Fatalf("second create shoud've failed")
	}
}

func Test_db__GetByID_not_found(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	user := &User{
		ID:   "test id",
		Hash: []byte("hash"),
	}

	fromDB, err := db.GetByID(ctx, user.ID)
	if fromDB != nil {
		t.Fatalf("nothing should be returned, got %+v", fromDB)
	}

	if err != sql.ErrNoRows {
		t.Fatalf("expected %s, got %s", sql.ErrNoRows, err)
	}
}

func Test_db__GetByID(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	user := &User{
		ID:   "test id",
		Hash: []byte("hash"),
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
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	user := &User{
		ID:       "test id",
		Username: "username",
		Hash:     []byte("hash"),
	}

	fromDB, err := db.GetByUsername(ctx, user.Username)
	if fromDB != nil {
		t.Fatalf("nothing should be returned, got %+v", fromDB)
	}

	if err != sql.ErrNoRows {
		t.Fatalf("expected %s, got %s", sql.ErrNoRows, err)
	}
}

func Test_db__GetByUsername(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	user := &User{
		ID:       "test id",
		Username: "username",
		Hash:     []byte("hash"),
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
	tmpFile, err := ioutil.TempFile(os.TempDir(), "testdb-")
	if err != nil {
		t.Fatalf("failed to create a temporary db file: %s", err)
	}

	t.Cleanup(func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Fatalf("failed to delete file for test db: %s", err)
		}
	})

	sqlite, err := db.New(&db.Config{
		Driver: "sqlite3",
		Addr:   tmpFile.Name(),
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

func (tl *testLogger) Info(string, ...interface{}) {}

func (tl *testLogger) Error(string, ...interface{}) {}
