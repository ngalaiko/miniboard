package tags

import (
	"context"
	"testing"
)

func Test__Create(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb)

	tag, err := service.Create(ctx, "user id", "title")
	if err != nil {
		t.Fatalf("failed to create a tag %s", err)
	}

	if tag.UserID != "user id" {
		t.Errorf("user id expected: %s, got %s", "user id", tag.UserID)
	}

	if tag.Title != "title" {
		t.Errorf("title expected: %s, got %s", "title", tag.Title)
	}

	if tag.ID == "" {
		t.Errorf("id expected to not be empty")
	}
}

func Test__Create_twice(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb)

	_, err := service.Create(ctx, "user id", "title")
	if err != nil {
		t.Fatalf("failed to create a tag %s", err)
	}

	_, secondErr := service.Create(ctx, "user id", "title")
	if secondErr != ErrAlreadyExists {
		t.Fatalf("expected %s, got %s", ErrAlreadyExists, secondErr)
	}
}
