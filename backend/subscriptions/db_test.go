package subscriptions

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/ngalaiko/miniboard/backend/db"
	"github.com/ngalaiko/miniboard/backend/tags"
)

func testSubscription() *Subscription {
	subscription := &Subscription{}
	subscription.ID = "test id"
	subscription.URL = "https://example.com"
	subscription.Title = "title"
	subscription.Created = time.Now().Add(-1 * time.Hour)
	subscription.IconURL = new(string)
	subscription.TagIDs = []string{}
	*subscription.IconURL = "https://icon.url"
	return subscription
}

func Test_db__Create(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	subscription := testSubscription()

	if err := db.Create(ctx, "user", subscription); err != nil {
		t.Fatalf("failed to create a subscription: %s", err)
	}
}

func Test_db__Create_twice(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	subscription := testSubscription()

	if err := db.Create(ctx, "user", subscription); err != nil {
		t.Fatalf("failed to create a subscription: %s", err)
	}

	if err := db.Create(ctx, "user", subscription); err == nil {
		t.Fatalf("second create shoud've failed")
	}
}

func Test_db__Create_twice_for_different_users(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	subscription1 := testSubscription()

	if err := db.Create(ctx, "user1", subscription1); err != nil {
		t.Fatalf("failed to create a subscription: %s", err)
	}
	fromDB1, err := db.Get(ctx, "user1", subscription1.ID)
	if err != nil {
		t.Fatalf("failed to get subscription from the db: %s", err)
	}
	if !cmp.Equal(subscription1, fromDB1) {
		t.Error(cmp.Diff(subscription1, fromDB1))
	}

	subscription2 := subscription1
	if err := db.Create(ctx, "user2", subscription1); err != nil {
		t.Fatalf("failed to create a subscription: %s", err)
	}
	fromDB2, err := db.Get(ctx, "user2", subscription2.ID)
	if err != nil {
		t.Fatalf("failed to get subscription from the db: %s", err)
	}
	if !cmp.Equal(subscription1, fromDB2) {
		t.Error(cmp.Diff(subscription1, fromDB1))
	}
}

func Test_db__Get_not_found(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	subscription := testSubscription()

	fromDB, err := db.Get(ctx, "user", subscription.ID)
	if fromDB != nil {
		t.Fatalf("nothing should be returned, got %+v", fromDB)
	}

	if !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("expected %s, got %s", sql.ErrNoRows, err)
	}
}
func Test_db__Get_without_tags(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	sqldb := createTestDB(ctx, t)
	db := newDB(sqldb, &testLogger{})

	subscription := testSubscription()
	subscription.Updated = new(time.Time)
	*subscription.Updated = time.Now().Truncate(time.Nanosecond)

	if err := db.Create(ctx, "user", subscription); err != nil {
		t.Fatalf("failed to create a subscription: %s", err)
	}

	fromDB, err := db.Get(ctx, "user", subscription.ID)
	if err != nil {
		t.Fatalf("failed to get subscription from the db: %s", err)
	}

	if !cmp.Equal(subscription, fromDB) {
		t.Error(cmp.Diff(subscription, fromDB))
	}
}

func Test_db__Get_with_tags(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	sqldb := createTestDB(ctx, t)
	db := newDB(sqldb, &testLogger{})

	tagService := tags.NewService(sqldb)
	tag1, err := tagService.Create(ctx, "user", "id1")
	if err != nil {
		t.Fatal(err)
	}

	tag2, err := tagService.Create(ctx, "user", "id2")
	if err != nil {
		t.Fatal(err)
	}

	subscription := testSubscription()
	subscription.TagIDs = []string{tag1.ID, tag2.ID}
	subscription.Updated = new(time.Time)
	*subscription.Updated = time.Now().Truncate(time.Nanosecond)

	if err := db.Create(ctx, "user", subscription); err != nil {
		t.Fatalf("failed to create a subscription: %s", err)
	}

	fromDB, err := db.Get(ctx, "user", subscription.ID)
	if err != nil {
		t.Fatalf("failed to get subscription from the db: %s", err)
	}

	if !cmp.Equal(subscription, fromDB) {
		t.Error(cmp.Diff(subscription, fromDB))
	}
}

func Test_db__List_paginated_by_created(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()
	db := newDB(createTestDB(ctx, t), &testLogger{})

	created := map[string]*Subscription{}
	for i := 0; i < 100; i++ {
		subscription := testSubscription()
		subscription.ID = fmt.Sprint(i)
		subscription.URL = fmt.Sprintf("https://example%d.com", i)
		subscription.Title = fmt.Sprintf("%d title", i)

		if err := db.Create(ctx, "user", subscription); err != nil {
			t.Fatal(err)
		}
		created[subscription.ID] = subscription
	}

	var createdLT *time.Time
	for i := 0; i < 20; i++ {
		subscriptions, err := db.List(ctx, "user", 5, createdLT)
		if err != nil {
			t.Fatal(err)
		}

		if len(subscriptions) != 5 {
			t.Errorf("expected 5 items, got %d", len(subscriptions))
		}

		for j, subscription := range subscriptions {
			expectedID := fmt.Sprint(99 - i*5 - j)
			if subscription.ID != expectedID {
				t.Fatalf("expected id %s, got %s", expectedID, subscription.ID)
				break
			}
			if !cmp.Equal(subscription, created[subscription.ID]) {
				t.Fatal(cmp.Diff(subscription, created[subscription.ID]))
			}
			createdLT = &subscription.Created
		}
	}
}

func createTestDB(ctx context.Context, t *testing.T) *sql.DB {
	sqlite, err := db.New(&db.Config{
		Addr: "sqlite3://file::memory::",
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

func (tl *testLogger) Info(string, ...interface{}) {}

func (tl *testLogger) Error(string, ...interface{}) {}
