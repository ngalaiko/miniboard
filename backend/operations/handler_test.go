package operations

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ngalaiko/miniboard/backend/authorizations"
)

func Test_Handler__Get_unauthorized(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	service := NewService(logger, testDB(ctx, t), nil)
	handler := NewHandler(service, logger)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.Get("")(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}

func Test_Handler__Get_not_found(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	db := testDB(ctx, t)
	service := NewService(logger, db, nil)
	go func(t *testing.T) {
		if err := service.Start(ctx); err != nil {
			t.Error(err)
		}
	}(t)
	handler := NewHandler(service, logger)

	req, err := http.NewRequest("GET", "/404", nil)
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(authorizations.NewContext(ctx, &authorizations.Token{
		UserID: "user",
	}))
	rr := httptest.NewRecorder()

	handler.Get("")(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func Test_Handler__Get_found(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	db := testDB(ctx, t)
	service := NewService(logger, db, nil)
	go func(t *testing.T) {
		if err := service.Start(ctx); err != nil {
			t.Error(err)
		}
	}(t)
	handler := NewHandler(service, logger)

	operation, err := service.Create(ctx, "user", func(context.Context, *Operation, chan<- *Operation) error {
		return nil
	})
	if err != nil {
		t.Fatalf("failed to create operation %s", err)
	}

	req, err := http.NewRequest("GET", "/"+operation.ID, nil)
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(authorizations.NewContext(ctx, &authorizations.Token{
		UserID: "user",
	}))
	rr := httptest.NewRecorder()

	handler.Get(operation.ID)(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	returned := &Operation{}
	if err := json.Unmarshal(rr.Body.Bytes(), &returned); err != nil {
		t.Fatalf("failed to unmarshal response body: %s", err)
	}

	if operation.ID != returned.ID {
		t.Errorf("expected id %s, got %s", operation.ID, returned.ID)
	}

	if operation.UserID != "user" {
		t.Errorf("expected user id %s, got %s", "user", returned.UserID)
	}
}
