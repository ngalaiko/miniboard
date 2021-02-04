package feeds

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"

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
	}
	if err := db.Create(ctx, feed1); err != nil {
		t.Fatalf("failed to create a feed: %s", err)
	}
	fromDB1, err := db.Get(ctx, "user1", feed1.ID)
	if err != nil {
		t.Fatalf("failed to get feed from the db: %s", err)
	}
	if !reflect.DeepEqual(feed1, fromDB1) {
		t.Fatalf("expected %+v, got %+v", feed1, fromDB1)
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
	if !reflect.DeepEqual(feed1, fromDB2) {
		t.Fatalf("expected %+v, got %+v", feed1, fromDB2)
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

	if !reflect.DeepEqual(feed, fromDB) {
		t.Fatalf("expected %+v, got %+v", feed, fromDB)
	}
}

func Test_db__List_tag_id_empty(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	for i := 0; i < 10; i++ {
		feed := &Feed{
			ID:      fmt.Sprint(i),
			UserID:  "user",
			URL:     fmt.Sprintf("https://example%d.com", i),
			Title:   fmt.Sprintf("%d title", i),
			Created: time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond),
			TagIDs: []string{
				fmt.Sprintf("%d", i),
				fmt.Sprintf("%d", i+1),
			},
		}

		if err := db.Create(ctx, feed); err != nil {
			t.Fatal(err)
		}
	}

	feed := &Feed{
		ID:      "test id",
		UserID:  "user",
		URL:     "https://example.com",
		Title:   "title",
		Created: time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond),
	}

	if err := db.Create(ctx, feed); err != nil {
		t.Fatal(err)
	}

	tagID := new(string)
	*tagID = ""

	feeds, err := db.List(ctx, "user", 5, nil, tagID)
	if err != nil {
		t.Fatal(err)
	}

	if len(feeds) != 1 {
		t.Fatalf("expected 1 feed, got %d", len(feeds))
	}

	if (len(feeds[0].TagIDs)) != 0 {
		t.Errorf("no tags expected, got %v", feeds[0].TagIDs)
	}
}

func Test_db__List_tag_ids(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	for i := 0; i < 10; i++ {
		feed := &Feed{
			ID:      fmt.Sprint(i),
			UserID:  "user",
			URL:     fmt.Sprintf("https://example%d.com", i),
			Title:   fmt.Sprintf("%d title", i),
			Created: time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond),
			TagIDs: []string{
				fmt.Sprintf("%d", i),
				fmt.Sprintf("%d", i+1),
			},
		}

		if err := db.Create(ctx, feed); err != nil {
			t.Fatal(err)
		}
	}

	tagID := new(string)
	*tagID = "5"

	feeds, err := db.List(ctx, "user", 5, nil, tagID)
	if err != nil {
		t.Fatal(err)
	}

	if len(feeds) != 2 {
		t.Errorf("expected 2 feeds")
	}

	for _, feed := range feeds {
		found := false
		for _, tagID := range feed.TagIDs {
			if tagID == "5" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected '5', got %v", feed.TagIDs)
		}
	}
}

func Test_db__List_paginated_by_created(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	for i := 0; i < 100; i++ {
		feed := &Feed{
			ID:      fmt.Sprint(i),
			UserID:  "user",
			URL:     fmt.Sprintf("https://example%d.com", i),
			Title:   fmt.Sprintf("%d title", i),
			Created: time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond),
		}
		feed.IconURL = new(string)
		*feed.IconURL = "https://icon.url"

		if err := db.Create(ctx, feed); err != nil {
			t.Fatal(err)
		}
	}

	var createdLT *time.Time
	for i := 0; i < 20; i++ {
		feeds, err := db.List(ctx, "user", 5, createdLT, nil)
		if err != nil {
			t.Fatal(err)
		}

		if len(feeds) != 5 {
			t.Errorf("expected 5 items, got %d", len(feeds))
		}

		for j, feed := range feeds {
			expectedID := fmt.Sprint(99 - i*5 - j)
			if feed.ID != expectedID {
				t.Errorf("expected id %s, got %s", expectedID, feed.ID)
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
