package items

import (
	"context"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test__Create(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &testLogger{})

	item, err := service.Create(ctx, "user id", "sid", mustParseURL("https://example.org"), "title")
	if err != nil {
		t.Fatalf("failed to create a item %s", err)
	}

	if item.UserID != "user id" {
		t.Errorf("user id expected: %s, got %s", "user id", item.UserID)
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

func Test__Create_twice(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &testLogger{})

	_, err := service.Create(ctx, "user id", "sid", mustParseURL("https://example.org"), "title")
	if err != nil {
		t.Fatalf("failed to create a item %s", err)
	}

	_, secondErr := service.Create(ctx, "user id", "sid", mustParseURL("https://example.org"), "title")
	if secondErr != errAlreadyExists {
		t.Fatalf("expected %s, got %s", errAlreadyExists, secondErr)
	}
}

func Test__Get(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &testLogger{})

	item, err := service.Create(ctx, "user id", "sid", mustParseURL("https://example.org"), "title")
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

func Test__Get_not_found(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &testLogger{})

	_, err := service.Get(ctx, "user id", "id")
	if err != errNotFound {
		t.Errorf("expected %s, got %s", errNotFound, err)
	}
}

func mustParseURL(raw string) *url.URL {
	url, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return url
}
