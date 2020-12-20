package jwt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ngalaiko/miniboard/server/users"
)

func Test_Handler__Get(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	handler := NewHandler(&testUsersService{}, NewService(createTestDB(ctx, t), logger), logger)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusMethodNotAllowed)
	}
}

func Test_Handler__Post_404(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	handler := NewHandler(&testUsersService{}, NewService(createTestDB(ctx, t), logger), logger)

	req, err := http.NewRequest("POST", "/404", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusNotFound)
	}
}

func Test_Handler__Create_unknown_user(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	tus := &testUsersService{Error: users.ErrNotFound}
	handler := NewHandler(tus, NewService(createTestDB(ctx, t), logger), logger)

	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(`
		{"username": "404"}
	`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := fmt.Sprintf(`{"error":"%s"}`, users.ErrNotFound)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_Handler__Create_invalid_password(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	tus := &testUsersService{User: &users.User{
		Username: "username",
		Hash:     []byte("$2a$14$.3JXNBhZzqfEKwqRZg8WV.kpelsYPEgs4wft/NgU9nRM1ZxomzXem"),
	}}
	handler := NewHandler(tus, NewService(createTestDB(ctx, t), logger), logger)

	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(`
		{"username": "username", "password": "invalid"}
	`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := fmt.Sprintf(`{"error":"%s"}`, users.ErrInvalidPassword)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_Handler__Create(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	tus := &testUsersService{User: &users.User{
		ID:       "id",
		Username: "username",
		Hash:     []byte("$2a$14$.3JXNBhZzqfEKwqRZg8WV.kpelsYPEgs4wft/NgU9nRM1ZxomzXem"),
	}}
	handler := NewHandler(tus, NewService(createTestDB(ctx, t), logger), logger)

	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(`
		{"username": "username", "password": "password"}
	`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	token := &Token{}
	if err := json.Unmarshal(rr.Body.Bytes(), token); err != nil {
		t.Fatalf("failed to unmarshal response: %s", err)
	}

	if token.Token == "" {
		t.Errorf("token is empty")
	}

	if token.UserID != "id" {
		t.Errorf("unexpected error id: got %s, expexted %s", token.UserID, "id")
	}
}

type testUsersService struct {
	User  *users.User
	Error error
}

func (tus *testUsersService) GetByUsername(_ context.Context, userID string) (*users.User, error) {
	return tus.User, tus.Error
}
