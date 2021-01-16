package feeds

import (
	"context"
	"database/sql"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/ngalaiko/miniboard/backend/db"
)

func Test_db__Create(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	feed := &Feed{
		ID:      "test id",
		UserID:  "user",
		URL:     "https://example.com",
		Title:   "title",
		Created: time.Now().Add(-1 * time.Hour),
	}

	if err := db.Create(ctx, feed); err != nil {
		t.Fatalf("failed to create a feed: %s", err)
	}
}

func Test_db__Create_twice(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	feed := &Feed{
		ID:      "test id",
		UserID:  "user",
		URL:     "https://example.com",
		Title:   "title",
		Created: time.Now().Add(-1 * time.Hour),
	}
	if err := db.Create(ctx, feed); err != nil {
		t.Fatalf("failed to create a feed: %s", err)
	}

	if err := db.Create(ctx, feed); err == nil {
		t.Fatalf("second create shoud've failed")
	}
}

func Test_db__Get_not_found(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	feed := &Feed{
		ID:      "test id",
		UserID:  "user",
		URL:     "https://example.com",
		Title:   "title",
		Created: time.Now().Add(-1 * time.Hour),
	}

	fromDB, err := db.Get(ctx, feed.UserID, feed.ID)
	if fromDB != nil {
		t.Fatalf("nothing should be returned, got %+v", fromDB)
	}

	if err != sql.ErrNoRows {
		t.Fatalf("expected %s, got %s", sql.ErrNoRows, err)
	}
}

func Test_db__Get(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	feed := &Feed{
		ID:      "test id",
		UserID:  "user",
		URL:     "https://example.com",
		Title:   "title",
		Created: time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond),
	}
	feed.Updated = new(time.Time)
	*feed.Updated = time.Now().Truncate(time.Nanosecond)

	if err := db.Create(ctx, feed); err != nil {
		t.Fatalf("failed to create a feed: %s", err)
	}

	fromDB, err := db.Get(ctx, feed.UserID, feed.ID)
	if err != nil {
		t.Fatalf("failed to get feed from the db: %s", err)
	}

	if !reflect.DeepEqual(feed, fromDB) {
		t.Fatalf("expected %+v, got %+v", feed, fromDB)
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
