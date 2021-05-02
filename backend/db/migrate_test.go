package db

import (
	"context"
	"database/sql"
	"testing"
)

func Test_Migrate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	db := testDB(t)

	if err := Migrate(ctx, db, &testLogger{}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
}

func Test_Migrate_twice(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	db := testDB(t)

	if err := Migrate(ctx, db, &testLogger{}); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if err := Migrate(ctx, db, &testLogger{}); err != nil {
		t.Fatalf("unexpected error during the second migration: %s", err)
	}
}

func testDB(t *testing.T) *sql.DB {
	sqldb, err := New(&Config{
		Driver: "sqlite3",
		Addr:   "file::memory:",
	}, &testLogger{})
	if err != nil {
		t.Fatalf("failed to create a db: %s", err)
	}

	return sqldb
}

type testLogger struct{}

func (l *testLogger) Debug(string, ...interface{}) {}
