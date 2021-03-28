package tags

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/ngalaiko/miniboard/backend/db"
)

func Test_db__Create(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	tag := &Tag{
		ID:      "test id",
		UserID:  "user",
		Title:   "title",
		Created: time.Now().Add(-1 * time.Hour),
	}

	if err := db.Create(ctx, tag); err != nil {
		t.Fatalf("failed to create a tag: %s", err)
	}
}

func Test_db__Create_twice(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	tag := &Tag{
		ID:      "test id",
		UserID:  "user",
		Title:   "title",
		Created: time.Now().Add(-1 * time.Hour),
	}
	if err := db.Create(ctx, tag); err != nil {
		t.Fatalf("failed to create a tag: %s", err)
	}

	if err := db.Create(ctx, tag); err == nil {
		t.Fatalf("second create shoud've failed")
	}
}

func Test_db__GetByTitle_not_found(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	tag := &Tag{
		ID:      "test id",
		UserID:  "user",
		Title:   "title",
		Created: time.Now().Add(-1 * time.Hour),
	}

	fromDB, err := db.GetByTitle(ctx, tag.UserID, tag.Title)
	if fromDB != nil {
		t.Fatalf("nothing should be returned, got %+v", fromDB)
	}

	if err != sql.ErrNoRows {
		t.Fatalf("expected %s, got %s", sql.ErrNoRows, err)
	}
}

func Test_db__GetByID_not_found(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	tag := &Tag{
		ID:      "test id",
		UserID:  "user",
		Title:   "title",
		Created: time.Now().Add(-1 * time.Hour),
	}

	fromDB, err := db.GetByID(ctx, tag.UserID, tag.Title)
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

	tag := &Tag{
		ID:      "test id",
		UserID:  "user",
		Title:   "title",
		Created: time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond),
	}

	if err := db.Create(ctx, tag); err != nil {
		t.Fatalf("failed to create a tag: %s", err)
	}

	fromDB, err := db.GetByID(ctx, tag.UserID, tag.ID)
	if err != nil {
		t.Fatalf("failed to get tag from the db: %s", err)
	}

	if !reflect.DeepEqual(tag, fromDB) {
		t.Fatalf("expected %+v, got %+v", tag, fromDB)
	}
}

func Test_db__GetByTitle(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	tag := &Tag{
		ID:      "test id",
		UserID:  "user",
		Title:   "title",
		Created: time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond),
	}

	if err := db.Create(ctx, tag); err != nil {
		t.Fatalf("failed to create a tag: %s", err)
	}

	fromDB, err := db.GetByTitle(ctx, tag.UserID, tag.Title)
	if err != nil {
		t.Fatalf("failed to get tag from the db: %s", err)
	}

	if !reflect.DeepEqual(tag, fromDB) {
		t.Fatalf("expected %+v, got %+v", tag, fromDB)
	}
}

func Test_db__List_paginated_by_created(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t))

	for i := 0; i < 100; i++ {
		tag := &Tag{
			ID:      fmt.Sprint(i),
			UserID:  "user",
			Title:   fmt.Sprintf("%d title", i),
			Created: time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond),
		}

		if err := db.Create(ctx, tag); err != nil {
			t.Fatal(err)
		}
	}

	var createdLT *time.Time
	for i := 0; i < 20; i++ {
		tags, err := db.List(ctx, "user", 5, createdLT)
		if err != nil {
			t.Fatal(err)
		}

		if len(tags) != 5 {
			t.Errorf("expected 5 items, got %d", len(tags))
		}

		for j, tag := range tags {
			expectedID := fmt.Sprint(99 - i*5 - j)
			if tag.ID != expectedID {
				t.Errorf("expected id %s, got %s", expectedID, tag.ID)
			}
			createdLT = &tag.Created
		}
	}
}

func createTestDB(ctx context.Context, t *testing.T) *sql.DB {
	sqlite, err := db.New(&db.Config{
		Driver: "sqlite3",
		Addr:   "file::memory:",
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
