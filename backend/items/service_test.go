package items

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_service__Create(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &testLogger{})

	item, err := service.Create(ctx, "sid", "https://example.org", "title", time.Now(), "content")
	if err != nil {
		t.Fatalf("failed to create a item %s", err)
	}

	if item.Title != "title" {
		t.Errorf("title expected: %s, got %s", "Sample Item", item.Title)
	}

	if item.ID == "" {
		t.Errorf("id expected to not be empty")
	}

	if item.URL != "https://example.org" {
		t.Errorf("url should be https://example.org, got %s", item.URL)
	}

	if item.SubscriptionID != "sid" {
		t.Errorf("subscription id must be 'sid', got %s", item.SubscriptionID)
	}
}

func Test_service__Create_twice(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &testLogger{})

	if _, err := sqldb.Exec(`INSERT INTO users_subscriptions (user_id, subscription_id) VALUES ("user id", "sid")`); err != nil {
		t.Fatal(err)
	}

	_, err := service.Create(ctx, "sid", "https://example.org", "title", time.Now(), "content")
	if err != nil {
		t.Fatalf("failed to create a item %s", err)
	}

	_, secondErr := service.Create(ctx, "sid", "https://example.org", "title", time.Now(), "content")
	if secondErr != ErrAlreadyExists {
		t.Fatalf("expected %s, got %s", ErrAlreadyExists, secondErr)
	}
}

func Test_service__Get(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &testLogger{})

	if _, err := sqldb.Exec(`INSERT INTO users_subscriptions (user_id, subscription_id) VALUES ("user id", "sid")`); err != nil {
		t.Fatal(err)
	}

	item, err := service.Create(ctx, "sid", "https://example.org", "title", time.Now(), "content")
	if err != nil {
		t.Fatalf("failed to create a item: %s", err)
	}

	from, err := service.Get(ctx, item.ID, "user id")
	if err != nil {
		t.Fatalf("failed to get a item: %s", err)
	}

	if !cmp.Equal(*item, from.Item) {
		t.Error(cmp.Diff(*item, from.Item))
	}
}

func Test_service__Get_not_found(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &testLogger{})

	_, err := service.Get(ctx, "user id", "id")
	if err != errNotFound {
		t.Errorf("expected %s, got %s", errNotFound, err)
	}
}
