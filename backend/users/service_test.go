package users

import (
	"context"
	"reflect"
	"testing"
)

func Test_Service_Create(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &Config{BCryptCost: 10})

	user, err := service.Create(ctx, "username", []byte("password"))
	if err != nil {
		t.Fatalf("failed to create a user: %s", err)
	}

	if user.ID == "" {
		t.Fatalf("user id must not be empty")
	}

	if len(user.Hash) == 0 {
		t.Fatalf("user hash must not be empty")
	}
}

func Test_Service_Create__twice(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &Config{BCryptCost: 10})

	if _, err := service.Create(ctx, "username", []byte("password")); err != nil {
		t.Fatalf("failed to create a user: %s", err)
	}

	_, err := service.Create(ctx, "username", []byte("password"))
	if err != ErrAlreadyExists {
		t.Fatalf("expected %s, got %s", ErrAlreadyExists, err)
	}
}

func Test_Service_GetByID(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &Config{BCryptCost: 10})

	created, err := service.Create(ctx, "username", []byte("password"))
	if err != nil {
		t.Fatalf("failed to create a user: %s", err)
	}

	found, err := service.GetByID(ctx, created.ID)
	if err != nil {
		t.Fatalf("failed to get a user: %s", err)
	}

	if !reflect.DeepEqual(found, created) {
		t.Fatalf("expected %+v, got %+v", created, found)
	}
}

func Test_Service_GetByID__not_found(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &Config{BCryptCost: 10})

	_, err := service.GetByID(ctx, "id")
	if err != ErrNotFound {
		t.Fatalf("expected %s, got %s", ErrNotFound, err)
	}
}

func Test_Service_GetByUsername(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &Config{BCryptCost: 10})

	created, err := service.Create(ctx, "username", []byte("password"))
	if err != nil {
		t.Fatalf("failed to create a user: %s", err)
	}

	found, err := service.GetByUsername(ctx, created.Username)
	if err != nil {
		t.Fatalf("failed to get a user: %s", err)
	}

	if !reflect.DeepEqual(found, created) {
		t.Fatalf("expected %+v, got %+v", created, found)
	}
}

func Test_Service_GetByUsername__not_found(t *testing.T) {
	ctx := context.TODO()

	sqldb := createTestDB(ctx, t)
	service := NewService(sqldb, &Config{BCryptCost: 10})

	_, err := service.GetByUsername(ctx, "username")
	if err != ErrNotFound {
		t.Fatalf("expected %s, got %s", ErrNotFound, err)
	}
}
