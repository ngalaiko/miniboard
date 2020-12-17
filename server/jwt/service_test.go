package jwt

import (
	"context"
	"database/sql"
	"io/ioutil"
	"os"
	"testing"

	"github.com/ngalaiko/miniboard/server/db"
)

func Test_Init(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)

	service := NewService(sqldb, &testLogger{})

	if err := service.Init(ctx); err != nil {
		t.Fatalf("unexpected error")
	}
}

func Test_Init__twice(t *testing.T) {
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
