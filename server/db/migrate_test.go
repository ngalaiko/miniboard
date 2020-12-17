package db

import (
	"context"
	"database/sql"
	"io/ioutil"
	"os"
	"testing"
)

func Test_Migrate(t *testing.T) {
	ctx := context.Background()

	db := testDB(ctx, t)

	if err := Migrate(ctx, db, &testLogger{}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func Test_Migrate_twice(t *testing.T) {
	ctx := context.Background()
	db := testDB(ctx, t)

	if err := Migrate(ctx, db, &testLogger{}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if err := Migrate(ctx, db, &testLogger{}); err != nil {
		t.Fatalf("unexpected error during the second migration: %s", err)
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

	sqldb, err := New(&Config{
		Driver: "sqlite3",
		Addr:   tmpFile.Name(),
	}, &testLogger{})
	if err != nil {
		t.Fatalf("failed to create a db: %s", err)
	}

	return sqldb
}

type testLogger struct{}

func (l *testLogger) Info(string, ...interface{}) {}
