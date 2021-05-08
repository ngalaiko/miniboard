package items

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_service__Create(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &testLogger{})

	now := time.Now()
	summary := "summary"
	item, err := service.Create(ctx, "sid", "https://example.org", "title", &now, &summary)
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

func Test_service__Create_without_created(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &testLogger{})

	if _, err := sqldb.Exec(`INSERT INTO users_subscriptions (user_id, subscription_id) VALUES ("user id", "sid")`); err != nil {
		t.Fatal(err)
	}

	_, err := service.Create(ctx, "sid", "https://example.org", "title", nil, nil)
	if err != nil {
		t.Fatalf("failed to create: %s", err)
	}
}

func Test_service__Create_twice(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &testLogger{})

	if _, err := sqldb.Exec(`INSERT INTO users_subscriptions (user_id, subscription_id) VALUES ("user id", "sid")`); err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	_, err := service.Create(ctx, "sid", "https://example.org", "title", &now, nil)
	if err != nil {
		t.Fatalf("failed to create a item %s", err)
	}

	_, secondErr := service.Create(ctx, "sid", "https://example.org", "title", &now, nil)
	if !errors.Is(secondErr, ErrAlreadyExists) {
		t.Fatalf("expected %s, got %s", ErrAlreadyExists, secondErr)
	}
}

func Test_service__Get(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &testLogger{})

	if _, err := sqldb.Exec(`INSERT INTO users_subscriptions (user_id, subscription_id) VALUES ("user id", "sid")`); err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	summary := "summary"
	item, err := service.Create(ctx, "sid", "https://example.org", "title", &now, &summary)
	if err != nil {
		t.Fatalf("failed to create a item: %s", err)
	}

	from, err := service.Get(ctx, "user id", item.ID)
	if err != nil {
		t.Fatalf("failed to get a item: %s", err)
	}

	if !cmp.Equal(item, from) {
		t.Error(cmp.Diff(item, from))
	}
}

func Test_service__Get_not_found(t *testing.T) {
	t.Parallel()

	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &testLogger{})

	_, err := service.Get(ctx, "id", "user id")
	if !errors.Is(err, ErrNotFound) {
		t.Errorf("expected %s, got %s", ErrNotFound, err)
	}
}
