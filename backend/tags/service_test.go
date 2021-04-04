package tags

import (
	"context"
	"reflect"
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

	if !reflect.DeepEqual(tag, from) {
		t.Errorf("unexpected response, expected %+v, got %+v", tag, from)
	}
}

func Test__Get_not_found(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb)

	_, err := service.Get(ctx, "user id", "id")
	if err != ErrNotFound {
		t.Errorf("expected %s, got %s", ErrNotFound, err)
	}
}
