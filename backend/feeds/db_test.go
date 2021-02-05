package feeds

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ngalaiko/miniboard/backend/db"
)

func Test_db__Create(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	feed := &Feed{
		ID:      "test id",
		UserID:  "user",
		URL:     "https://example.com",
		Title:   "title",
		Created: time.Now().Add(-1 * time.Hour),
	}
	feed.IconURL = new(string)
	*feed.IconURL = "https://icon.url"

	if err := db.Create(ctx, feed); err != nil {
		t.Fatalf("failed to create a feed: %s", err)
	}
}

func Test_db__Create_twice(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	feed := &Feed{
		ID:      "test id",
		UserID:  "user",
		URL:     "https://example.com",
		Title:   "title",
		Created: time.Now().Add(-1 * time.Hour),
	}
	feed.IconURL = new(string)
	*feed.IconURL = "https://icon.url"
	if err := db.Create(ctx, feed); err != nil {
		t.Fatalf("failed to create a feed: %s", err)
	}

	if err := db.Create(ctx, feed); err == nil {
		t.Fatalf("second create shoud've failed")
	}
}

func Test_db__Create_twice_for_different_users(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	feed1 := &Feed{
		ID:      "test id",
		UserID:  "user1",
		URL:     "https://example.com",
		Title:   "title",
		Created: time.Now().Add(-1 * time.Hour).Truncate(time.Millisecond),
		TagIDs:  []string{},
	}
	if err := db.Create(ctx, feed1); err != nil {
		t.Fatalf("failed to create a feed: %s", err)
	}
	fromDB1, err := db.Get(ctx, "user1", feed1.ID)
	if err != nil {
		t.Fatalf("failed to get feed from the db: %s", err)
	}
	if !cmp.Equal(feed1, fromDB1) {
		t.Error(cmp.Diff(feed1, fromDB1))
	}

	feed2 := &(*feed1)
	feed2.UserID = "user2"
	if err := db.Create(ctx, feed1); err != nil {
		t.Fatalf("failed to create a feed: %s", err)
	}
	fromDB2, err := db.Get(ctx, "user2", feed2.ID)
	if err != nil {
		t.Fatalf("failed to get feed from the db: %s", err)
	}
	if !cmp.Equal(feed1, fromDB2) {
		t.Error(cmp.Diff(feed1, fromDB1))
	}
}

func Test_db__Get_not_found(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	feed := &Feed{
		ID:      "test id",
		UserID:  "user",
		URL:     "https://example.com",
		Title:   "title",
		Created: time.Now().Add(-1 * time.Hour),
	}
	feed.IconURL = new(string)
	*feed.IconURL = "https://icon.url"

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
	db := newDB(createTestDB(ctx, t), &testLogger{})

	feed := &Feed{
		ID:      "test id",
		UserID:  "user",
		URL:     "https://example.com",
		Title:   "title",
		Created: time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond),
		TagIDs:  []string{"id1", "id2"},
	}
	feed.IconURL = new(string)
	*feed.IconURL = "https://icon.url"
	feed.Updated = new(time.Time)
	*feed.Updated = time.Now().Truncate(time.Nanosecond)

	if err := db.Create(ctx, feed); err != nil {
		t.Fatalf("failed to create a feed: %s", err)
	}

	fromDB, err := db.Get(ctx, feed.UserID, feed.ID)
	if err != nil {
		t.Fatalf("failed to get feed from the db: %s", err)
	}

	if !cmp.Equal(feed, fromDB) {
		t.Error(cmp.Diff(feed, fromDB))
	}
}

func Test_db__List_paginated_by_created(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	created := map[string]*Feed{}
	for i := 0; i < 100; i++ {
		feed := &Feed{
			ID:      fmt.Sprint(i),
			UserID:  "user",
			URL:     fmt.Sprintf("https://example%d.com", i),
			Title:   fmt.Sprintf("%d title", i),
			Created: time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond),
			TagIDs:  []string{},
		}
		feed.IconURL = new(string)
		*feed.IconURL = "https://icon.url"

		if err := db.Create(ctx, feed); err != nil {
			t.Fatal(err)
		}
		created[feed.ID] = feed
	}

	var createdLT *time.Time
	for i := 0; i < 20; i++ {
		feeds, err := db.List(ctx, "user", 5, createdLT)
		if err != nil {
			t.Fatal(err)
		}

		if len(feeds) != 5 {
			t.Errorf("expected 5 items, got %d", len(feeds))
		}

		for j, feed := range feeds {
			expectedID := fmt.Sprint(99 - i*5 - j)
			if feed.ID != expectedID {
				t.Fatalf("expected id %s, got %s", expectedID, feed.ID)
				break
			}
			if !cmp.Equal(feed, created[feed.ID]) {
				t.Fatal(cmp.Diff(feed, created[feed.ID]))
			}
			createdLT = &feed.Created
		}
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
