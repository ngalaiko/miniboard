package tags

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
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

	tag1, err := service.Create(ctx, "user id", "title")
	if err != nil {
		t.Fatalf("failed to create a tag %s", err)
	}

	tag2, secondErr := service.Create(ctx, "user id", "title")
	if secondErr != nil {
		t.Fatalf("failed to create a tag %s", secondErr)
	}

	if !cmp.Equal(tag1, tag2) {
		t.Error(cmp.Diff(tag1, tag2))
	}
}

func Test__Get(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb)

	tag, err := service.Create(ctx, "user id", "title")
	if err != nil {
		t.Fatalf("failed to create a tag: %s", err)
	}

	from, err := service.Get(ctx, "user id", tag.ID)
	if err != nil {
		t.Fatalf("failed to get a tag: %s", err)
	}

	if !cmp.Equal(tag, from) {
		t.Error(cmp.Diff(tag, from))
	}
}

func Test__Get_not_found(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb)

	_, err := service.Get(ctx, "user id", "id")
	if err != errNotFound {
		t.Errorf("expected %s, got %s", errNotFound, err)
	}
}
