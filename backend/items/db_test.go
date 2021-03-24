package items

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

func testUserItem() *UserItem {
	return &UserItem{
		UserID: "user",
		Item: Item{
			ID:             "test id",
			URL:            "https://example.com",
			Title:          "title",
			SubscriptionID: "sid",
			Created:        time.Now().Add(-1 * time.Hour),
			Summary:        "example summary",
		},
	}
}

func Test_db__Create(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	if err := db.Create(ctx, &testUserItem().Item); err != nil {
		t.Fatalf("failed to create a item: %s", err)
	}
}

func Test_db__Create_twice(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	item := testUserItem()

	if err := db.Create(ctx, &item.Item); err != nil {
		t.Fatalf("failed to create a item: %s", err)
	}

	if err := db.Create(ctx, &item.Item); err == nil {
		t.Fatalf("second create shoud've failed")
	}
}

func Test_db__GetByURL_not_found(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	item := testUserItem()

	fromDB, err := db.Get(ctx, item.UserID, item.URL)
	if fromDB != nil {
		t.Fatalf("nothing should be returned, got %+v", fromDB)
	}

	if err != sql.ErrNoRows {
		t.Fatalf("expected %s, got %s", sql.ErrNoRows, err)
	}
}

func Test_db__GetByURL(t *testing.T) {
	ctx := context.TODO()
	sqldb := createTestDB(ctx, t)
	db := newDB(sqldb, &testLogger{})

	if _, err := sqldb.Exec(`INSERT INTO users_subscriptions (user_id, subscription_id) VALUES ("user", "sid")`); err != nil {
		t.Fatal(err)
	}

	item := testUserItem()

	if err := db.Create(ctx, &item.Item); err != nil {
		t.Fatalf("failed to create a item: %s", err)
	}

	fromDB, err := db.GetByURL(ctx, item.URL)
	if err != nil {
		t.Fatalf("failed to get item from the db: %s", err)
	}

	if !cmp.Equal(&item.Item, fromDB) {
		t.Error(cmp.Diff(&item.Item, fromDB))
	}
}

func Test_db__Get_not_found(t *testing.T) {
	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	item := testUserItem()

	fromDB, err := db.Get(ctx, item.UserID, item.ID)
	if fromDB != nil {
		t.Fatalf("nothing should be returned, got %+v", fromDB)
	}

	if err != sql.ErrNoRows {
		t.Fatalf("expected %s, got %s", sql.ErrNoRows, err)
	}
}

func Test_db__Get(t *testing.T) {
	ctx := context.TODO()
	sqldb := createTestDB(ctx, t)
	db := newDB(sqldb, &testLogger{})

	if _, err := sqldb.Exec(`INSERT INTO users_subscriptions (user_id, subscription_id) VALUES ("user", "sid")`); err != nil {
		t.Fatal(err)
	}

	item := testUserItem()

	if err := db.Create(ctx, &item.Item); err != nil {
		t.Fatalf("failed to create a item: %s", err)
	}

	fromDB, err := db.Get(ctx, item.UserID, item.ID)
	if err != nil {
		t.Fatalf("failed to get item from the db: %s", err)
	}

	if !cmp.Equal(item, fromDB) {
		t.Error(cmp.Diff(item, fromDB))
	}
}

func Test_db__List_paginated_by_created(t *testing.T) {
	ctx := context.TODO()
	sqldb := createTestDB(ctx, t)
	db := newDB(sqldb, &testLogger{})

	if _, err := sqldb.Exec(`INSERT INTO users_subscriptions (user_id, subscription_id) VALUES ("user", "sid")`); err != nil {
		t.Fatal(err)
	}

	created := map[string]*UserItem{}
	for i := 0; i < 100; i++ {
		item := &Item{
			ID:             fmt.Sprint(i),
			URL:            fmt.Sprintf("https://example%d.com", i),
			Title:          fmt.Sprintf("%d title", i),
			SubscriptionID: "sid",
			Created:        time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond),
		}

		if err := db.Create(ctx, item); err != nil {
			t.Fatal(err)
		}
		created[item.ID] = &UserItem{
			UserID: "user",
			Item:   *item,
		}
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
			createdLT = &item.Created
		}
	}
}

func Test_db__List_filtered_by_tag(t *testing.T) {
	ctx := context.TODO()
	sqldb := createTestDB(ctx, t)
	db := newDB(sqldb, &testLogger{})

	created := map[string]*UserItem{}
	for i := 0; i < 100; i++ {
		item := &UserItem{
			UserID: "user",
			Item: Item{
				ID:             fmt.Sprint(i),
				URL:            fmt.Sprintf("https://example%d.com", i),
				Title:          fmt.Sprintf("%d title", i),
				Created:        time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond),
				SubscriptionID: fmt.Sprint(i % 5),
			},
		}

		_, _ = sqldb.Exec(`INSERT INTO users_subscriptions (user_id, subscription_id) VALUES ("user", $1)`, item.SubscriptionID)
		_, _ = sqldb.Exec(`INSERT INTO tags_subscriptions (tag_id, subscription_id) VALUES ("tag", $1)`, item.SubscriptionID)

		if err := db.Create(ctx, &item.Item); err != nil {
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
	ctx := context.TODO()
	sqldb := createTestDB(ctx, t)
	db := newDB(sqldb, &testLogger{})

	created := map[string]*UserItem{}
	for i := 0; i < 100; i++ {
		item := &UserItem{
			UserID: "user",
			Item: Item{
				ID:             fmt.Sprint(i),
				URL:            fmt.Sprintf("https://example%d.com", i),
				Title:          fmt.Sprintf("%d title", i),
				Created:        time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond),
				SubscriptionID: fmt.Sprint(i % 5),
			},
		}

		_, _ = sqldb.Exec(`INSERT INTO users_subscriptions (user_id, subscription_id) VALUES ("user", $1)`, item.SubscriptionID)
		_, _ = sqldb.Exec(`INSERT INTO tags_subscriptions (tag_id, subscription_id) VALUES ("tag", $1)`, item.SubscriptionID)

		if err := db.Create(ctx, &item.Item); err != nil {
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
	ctx := context.TODO()
	sqldb := createTestDB(ctx, t)
	db := newDB(sqldb, &testLogger{})

	created := map[string]*UserItem{}
	for i := 0; i < 100; i++ {
		item := &UserItem{
			UserID: "user",
			Item: Item{
				ID:             fmt.Sprint(i),
				URL:            fmt.Sprintf("https://example%d.com", i),
				Title:          fmt.Sprintf("%d title", i),
				Created:        time.Now().Add(-1 * time.Hour).Truncate(time.Nanosecond),
				SubscriptionID: fmt.Sprint(i % 5),
			},
		}

		_, _ = sqldb.Exec(`INSERT INTO users_subscriptions (user_id, subscription_id) VALUES ("user", $1)`, item.SubscriptionID)

		if err := db.Create(ctx, &item.Item); err != nil {
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
