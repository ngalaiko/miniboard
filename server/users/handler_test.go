package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Handler__Get(t *testing.T) {
	ctx := context.Background()
	handler := NewHandler(NewService(createTestDB(ctx, t)), &testLogger{})

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
	handler := NewHandler(NewService(createTestDB(ctx, t)), &testLogger{})

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

func Test_Handler__Post_create_username_empty(t *testing.T) {
	ctx := context.Background()
	handler := NewHandler(NewService(createTestDB(ctx, t)), &testLogger{})

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(``)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := fmt.Sprintf(`{"error":"%s"}`, ErrUsernameEmpty)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_Handler__Post_create_password_empty(t *testing.T) {
	ctx := context.Background()
	handler := NewHandler(NewService(createTestDB(ctx, t)), &testLogger{})

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(`
	{"username":"test"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := fmt.Sprintf(`{"error":"%s"}`, ErrPasswordEmpty)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_Handler__Post_already_exists(t *testing.T) {
	ctx := context.Background()
	handler := NewHandler(NewService(createTestDB(ctx, t)), &testLogger{})

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(`
	{"username":"test","password":"1234"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	req, err = http.NewRequest("POST", "/", bytes.NewBuffer([]byte(`
	{"username":"test","password":"1234"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := fmt.Sprintf(`{"error":"%s"}`, ErrAlreadyExists)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_Handler__Post_create_success(t *testing.T) {
	ctx := context.Background()
	handler := NewHandler(NewService(createTestDB(ctx, t)), &testLogger{})

	req, err := http.NewRequest("POST", "/", bytes.NewBuffer([]byte(`
	{"username":"test","password":"1234"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	user := &User{}
	if err := json.Unmarshal(rr.Body.Bytes(), &user); err != nil {
		t.Fatalf("failed to unmarshal response body: %s", err)
	}

	if user.Username != "test" {
		t.Errorf("expected username %s, got %s", "test", user.Username)
	}

	if user.ID == "" {
		t.Errorf("id is empty")
	}
}
