package tags

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ngalaiko/miniboard/backend/authorizations"
)

func Test_Handler__Get(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	db := createTestDB(ctx, t)
	service := NewService(db)
	handler := NewHandler(service, logger)

	tc := []struct {
		URL        string
		Error      error
		StatusCode int
	}{
		{
			URL:        "/?page_size=invalid",
			Error:      errInvalidPageSize,
			StatusCode: http.StatusBadRequest,
		},
		{
			URL:        "/?created_lt=invalid",
			Error:      errInvalidCreatedLT,
			StatusCode: http.StatusBadRequest,
		},
	}

	for _, testCase := range tc {
		t.Run(testCase.URL, func(t *testing.T) {
			req, err := http.NewRequest("GET", testCase.URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			req = req.WithContext(authorizations.NewContext(ctx, &authorizations.Token{
				UserID: "user",
			}))
			rr := httptest.NewRecorder()
			handler.List()(rr, req)

			expectedCode := testCase.StatusCode
			if status := rr.Code; status != expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, expectedCode)
			}

			expected := fmt.Sprintf(`{"message":"%s"}`, testCase.Error)
			if rr.Body.String() != expected {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), expected)
			}
		})
	}
}

func Test_Handler__Post_create_already_exists(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	db := createTestDB(ctx, t)
	service := NewService(db)
	handler := NewHandler(service, logger)

	var rr *httptest.ResponseRecorder
	for i := 0; i < 2; i++ {
		req, err := http.NewRequest("POST", "/", strings.NewReader(`{"title":"title"}`))
		if err != nil {
			t.Fatal(err)
		}

		req = req.WithContext(authorizations.NewContext(ctx, &authorizations.Token{
			UserID: "user",
		}))
		rr = httptest.NewRecorder()

		handler.Create()(rr, req)
	}

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := fmt.Sprintf(`{"message":"%s"}`, errAlreadyExists)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_Handler__Post_create_title_empty(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	db := createTestDB(ctx, t)
	service := NewService(db)
	handler := NewHandler(service, logger)

	req, err := http.NewRequest("POST", "/", strings.NewReader(`{}`))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(authorizations.NewContext(ctx, &authorizations.Token{
		UserID: "user",
	}))
	rr := httptest.NewRecorder()

	handler.Create()(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := fmt.Sprintf(`{"message":"%s"}`, errEmptyTitle)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_Handler__Post_create_success(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	db := createTestDB(ctx, t)
	service := NewService(db)
	handler := NewHandler(service, logger)

	req, err := http.NewRequest("POST", "/", strings.NewReader(`{"title": "title"}`))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(authorizations.NewContext(ctx, &authorizations.Token{
		UserID: "user",
	}))
	rr := httptest.NewRecorder()

	handler.Create()(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	returned := &Tag{}
	if err := json.Unmarshal(rr.Body.Bytes(), &returned); err != nil {
		t.Fatalf("failed to unmarshal response body: %s", err)
	}

	if returned.ID == "" {
		t.Errorf("id must not be empty")
	}

	if returned.Title != "title" {
		t.Errorf("expected title, got %+v", returned.Title)
	}
}
