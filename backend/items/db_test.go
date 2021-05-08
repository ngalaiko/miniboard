package items

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ngalaiko/miniboard/backend/db"
)

func testItem() *Item {
	return &Item{
		ID:             "test id",
		URL:            "https://example.com",
		Title:          "title",
		SubscriptionID: "sid",
		Created:        nil,
		Summary:        nil,
	}
}

func Test_db__Create(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	if err := db.Create(ctx, testItem()); err != nil {
		t.Fatalf("failed to create a item: %s", err)
	}
}

func Test_db__Create_twice(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	item := testItem()

	if err := db.Create(ctx, item); err != nil {
		t.Fatalf("failed to create a item: %s", err)
	}

	if err := db.Create(ctx, item); err == nil {
		t.Fatalf("second create shoud've failed")
	}
}

func Test_db__GetByURL_not_found(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	item := testItem()

	fromDB, err := db.Get(ctx, "user", item.URL)
	if fromDB != nil {
		t.Fatalf("nothing should be returned, got %+v", fromDB)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected %s, got %s", sql.ErrNoRows, err)
	}
}

func Test_db__GetByURL(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	sqldb := createTestDB(ctx, t)
	db := newDB(sqldb, &testLogger{})

	if _, err := sqldb.Exec(`INSERT INTO users_subscriptions (user_id, subscription_id) VALUES ("user", "sid")`); err != nil {
		t.Fatal(err)
	}

	item := testItem()

	if err := db.Create(ctx, item); err != nil {
		t.Fatalf("failed to create a item: %s", err)
	}

	fromDB, err := db.GetByURL(ctx, item.URL)
	if err != nil {
		t.Fatalf("failed to get item from the db: %s", err)
	}

	if !cmp.Equal(item, fromDB) {
		t.Error(cmp.Diff(item, fromDB))
	}
}

func Test_db__Get_not_found(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	item := testItem()

	fromDB, err := db.Get(ctx, "user", item.ID)
	if fromDB != nil {
		t.Fatalf("nothing should be returned, got %+v", fromDB)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected %s, got %s", sql.ErrNoRows, err)
	}
}

func Test_db__Get(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	sqldb := createTestDB(ctx, t)
	db := newDB(sqldb, &testLogger{})

	if _, err := sqldb.Exec(`INSERT INTO users_subscriptions (user_id, subscription_id) VALUES ("user", "sid")`); err != nil {
		t.Fatal(err)
	}

	item := testItem()

	if err := db.Create(ctx, item); err != nil {
		t.Fatalf("failed to create a item: %s", err)
	}

	fromDB, err := db.Get(ctx, "user", item.ID)
	if err != nil {
		t.Fatalf("failed to get item from the db: %s", err)
	}

	if !cmp.Equal(item, fromDB) {
		t.Error(cmp.Diff(item, fromDB))
	}
}

func Test_db__List_paginated_by_created(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	sqldb := createTestDB(ctx, t)
	db := newDB(sqldb, &testLogger{})

	if _, err := sqldb.Exec(`INSERT INTO users_subscriptions (user_id, subscription_id) VALUES ("user", "sid")`); err != nil {
		t.Fatal(err)
	}

	created := map[string]*Item{}
	for i := 0; i < 100; i++ {
		createdTS := time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond)
		item := &Item{
			ID:             fmt.Sprint(i),
			URL:            fmt.Sprintf("https://example%d.com", i),
			Title:          fmt.Sprintf("%d title", i),
			SubscriptionID: "sid",
			Created:        &createdTS,
		}

		if err := db.Create(ctx, item); err != nil {
			t.Fatal(err)
		}
		created[item.ID] = item
	}

	var createdLT *time.Time
	for i := 0; i < 20; i++ {
		items, err := db.List(ctx, "user", 5, createdLT, nil, nil)
		if err != nil {
			t.Fatal(err)
		}

		if len(items) != 5 {
			t.Errorf("expected 5 items, got %d", len(items))
		}

		for j, item := range items {
			expectedID := fmt.Sprint(99 - i*5 - j)
			if item.ID != expectedID {
				t.Fatalf("expected id %s, got %s", expectedID, item.ID)
				break
			}
			if !cmp.Equal(item, created[item.ID]) {
				t.Fatal(cmp.Diff(item, created[item.ID]))
			}
			createdLT = item.Created
		}
	}
}

func Test_db__List_filtered_by_tag(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	sqldb := createTestDB(ctx, t)
	db := newDB(sqldb, &testLogger{})

	created := map[string]*Item{}
	for i := 0; i < 100; i++ {
		createdTS := time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond)
		item := &Item{
			ID:             fmt.Sprint(i),
			URL:            fmt.Sprintf("https://example%d.com", i),
			Title:          fmt.Sprintf("%d title", i),
			Created:        &createdTS,
			SubscriptionID: fmt.Sprint(i % 5),
		}

		_, _ = sqldb.Exec(`INSERT INTO users_subscriptions (user_id, subscription_id) VALUES ("user", $1)`, item.SubscriptionID)
		_, _ = sqldb.Exec(`INSERT INTO tags_subscriptions (tag_id, subscription_id) VALUES ("tag", $1)`, item.SubscriptionID)

		if err := db.Create(ctx, item); err != nil {
			t.Fatal(err)
		}
		created[item.ID] = item
	}

	tagID := "tag"
	items, err := db.List(ctx, "user", 100, nil, nil, &tagID)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 100 {
		t.Fatalf("expected 100 items, got %d", len(items))
	}

	for _, item := range items {
		if !cmp.Equal(item, created[item.ID]) {
			t.Error(cmp.Diff(item, created[item.ID]))
		}
	}
}

func Test_db__List_filtered_by_subscription_and_tag(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	sqldb := createTestDB(ctx, t)
	db := newDB(sqldb, &testLogger{})

	created := map[string]*Item{}
	for i := 0; i < 100; i++ {
		createdTS := time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond)
		item := &Item{
			ID:             fmt.Sprint(i),
			URL:            fmt.Sprintf("https://example%d.com", i),
			Title:          fmt.Sprintf("%d title", i),
			Created:        &createdTS,
			SubscriptionID: fmt.Sprint(i % 5),
		}

		_, _ = sqldb.Exec(`INSERT INTO users_subscriptions (user_id, subscription_id) VALUES ("user", $1)`, item.SubscriptionID)
		_, _ = sqldb.Exec(`INSERT INTO tags_subscriptions (tag_id, subscription_id) VALUES ("tag", $1)`, item.SubscriptionID)

		if err := db.Create(ctx, item); err != nil {
			t.Fatal(err)
		}
		created[item.ID] = item
	}

	sID := "2"
	tagID := "tag"
	items, err := db.List(ctx, "user", 100, nil, &sID, &tagID)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 20 {
		t.Fatalf("expected 20 items, got %d", len(items))
	}

	for _, item := range items {
		if !cmp.Equal(item, created[item.ID]) {
			t.Error(cmp.Diff(item, created[item.ID]))
		}
	}
}

func Test_db__List_filtered_by_subscription(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	sqldb := createTestDB(ctx, t)
	db := newDB(sqldb, &testLogger{})

	created := map[string]*Item{}
	for i := 0; i < 100; i++ {
		createdTS := time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond)
		item := &Item{
			ID:             fmt.Sprint(i),
			URL:            fmt.Sprintf("https://example%d.com", i),
			Title:          fmt.Sprintf("%d title", i),
			Created:        &createdTS,
			SubscriptionID: fmt.Sprint(i % 5),
		}

		_, _ = sqldb.Exec(`INSERT INTO users_subscriptions (user_id, subscription_id) VALUES ("user", $1)`, item.SubscriptionID)

		if err := db.Create(ctx, item); err != nil {
			t.Fatal(err)
		}
		created[item.ID] = item
	}

	sID := "2"
	items, err := db.List(ctx, "user", 100, nil, &sID, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 20 {
		t.Fatalf("expected 20 items, got %d", len(items))
	}

	for _, item := range items {
		if !cmp.Equal(item, created[item.ID]) {
			t.Error(cmp.Diff(item, created[item.ID]))
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

func (tl *testLogger) Debug(string, ...interface{}) {}

func (tl *testLogger) Error(string, ...interface{}) {}
