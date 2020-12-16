package db

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/ngalaiko/miniboard/server/db"
)

func Test_Create(t *testing.T) {
	ctx := context.Background()

	sqlite := testDB(ctx, t)

	database := New(sqlite)

	testKey := &PublicKey{
		ID:  "test",
		Der: []byte("some data"),
	}

	if err := database.Create(ctx, testKey); err != nil {
		t.Fatalf("faied to create a token: %s", err)
	}
}

func Test_Create_twice(t *testing.T) {
	ctx := context.Background()

	sqlite := testDB(ctx, t)

	database := New(sqlite)

	testKey := &PublicKey{
		ID:  "test",
		Der: []byte("some data"),
	}

	if err := database.Create(ctx, testKey); err != nil {
		t.Fatalf("faied to create a token: %s", err)
	}

	if err := database.Create(ctx, testKey); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func Test_Get(t *testing.T) {
	ctx := context.Background()

	sqlite := testDB(ctx, t)

	database := New(sqlite)

	testKey := &PublicKey{
		ID:  "test",
		Der: []byte("some data"),
	}

	if err := database.Create(ctx, testKey); err != nil {
		t.Fatalf("faied to create a token: %s", err)
	}

	key, err := database.Get(ctx, testKey.ID)
	if err != nil {
		t.Fatalf("failed to get key from the db: %s", err)
	}

	if !reflect.DeepEqual(testKey, key) {
		t.Fatalf("expected %+v, got %+v", testKey, key)
	}
}

func Test_Get_not_exists(t *testing.T) {
	ctx := context.Background()

	sqlite := testDB(ctx, t)

	database := New(sqlite)

	testKey := &PublicKey{
		ID:  "test",
		Der: []byte("some data"),
	}

	_, err := database.Get(ctx, testKey.ID)
	if err != sql.ErrNoRows {
		t.Fatalf("expected %s, got %s", sql.ErrNoRows, err)
	}
}

func Test_Delete(t *testing.T) {
	ctx := context.Background()

	sqlite := testDB(ctx, t)

	database := New(sqlite)

	testKey := &PublicKey{
		ID:  "test",
		Der: []byte("some data"),
	}

	if err := database.Create(ctx, testKey); err != nil {
		t.Fatalf("faied to create a token: %s", err)
	}

	if _, err := database.Get(ctx, testKey.ID); err != nil {
		t.Fatalf("failed to get key back: %s", err)
	}

	if err := database.Delete(ctx, testKey.ID); err != nil {
		t.Fatalf("failed to delete key: %s", err)
	}

	if _, err := database.Get(ctx, testKey.ID); err != sql.ErrNoRows {
		t.Fatalf("key was not deleted: %s", err)
	}
}

func Test_Delete_not_existing(t *testing.T) {
	ctx := context.Background()

	sqlite := testDB(ctx, t)

	database := New(sqlite)

	testKey := &PublicKey{
		ID:  "test",
		Der: []byte("some data"),
	}

	if err := database.Delete(ctx, testKey.ID); err != nil {
		t.Fatalf("failed to delete key: %s", err)
	}
}

func Test_List(t *testing.T) {
	ctx := context.Background()

	sqlite := testDB(ctx, t)

	database := New(sqlite)

	kk, err := database.List(ctx)
	if err != nil {
		t.Fatalf("failed to list keys: %s", err)
	}
	if len(kk) != 0 {
		t.Fatalf("expected 0 tokens, got %d", len(kk))
	}

	for i := 0; i < 10; i++ {
		testKey := &PublicKey{
			ID:  fmt.Sprintf("test-%d", i),
			Der: []byte("some data"),
		}

		if err := database.Create(ctx, testKey); err != nil {
			t.Fatalf("faied to create a token: %s", err)
		}
	}

	kk, listError := database.List(ctx)
	if listError != nil {
		t.Fatalf("failed to list tokens: %s", err)
	}
	if len(kk) != 10 {
		t.Fatalf("expected 10 tokens, got %d", len(kk))
	}
}

func testDB(ctx context.Context, t *testing.T) *sql.DB {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "testdb-")
	if err != nil {
		t.Fatalf("failed to create file for test db: %s", err)
	}
	t.Cleanup(func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Fatalf("failed to delete file for test db: %s", err)
		}
	})

	sqldb, err := db.New(&db.Config{
		Driver: "sqlite3",
		Addr:   tmpFile.Name(),
	}, &testLogger{})
	if err != nil {
		t.Fatalf("failed to create a db: %s", err)
	}

	if err := db.Migrate(ctx, sqldb, &testLogger{}); err != nil {
		t.Fatalf("failed to run migrations: %s", err)
	}

	return sqldb
}

type testLogger struct{}

func (l *testLogger) Info(string, ...interface{}) {}
