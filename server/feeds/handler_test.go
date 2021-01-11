package feeds

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ngalaiko/miniboard/server/authorizations"
	"github.com/ngalaiko/miniboard/server/operations"
)

func Test_Handler__Get(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	db := createTestDB(ctx, t)
	crawler := &testCrawler{}
	service := NewService(db, crawler)
	handler := NewHandler(service, logger, &testOperationsService{})

	req, err := http.NewRequest("GET", "/", strings.NewReader(`{"url":"https://console.org"}`))
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
	db := createTestDB(ctx, t)
	crawler := &testCrawler{}
	service := NewService(db, crawler)
	handler := NewHandler(service, logger, &testOperationsService{})

	req, err := http.NewRequest("POST", "/404", strings.NewReader(`{"url":"https://console.org"}`))
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

func Test_Handler__Post_create_failed_create_operation(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	db := createTestDB(ctx, t)
	crawler := &testCrawler{}
	crawler = crawler.With("https://example.org", feedData)
	service := NewService(db, crawler)
	handler := NewHandler(service, logger, &testOperationsService{
		Error: fmt.Errorf("failed"),
	})

	req, err := http.NewRequest("POST", "/", strings.NewReader(`{"url":"https://console.org"}`))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(authorizations.NewContext(ctx, &authorizations.Token{
		UserID: "user",
	}))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	expected := fmt.Sprintf(`{"message":"internal server error"}`)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

func Test_Handler__Post_create_url_empty(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	db := createTestDB(ctx, t)
	crawler := &testCrawler{}
	service := NewService(db, crawler)
	handler := NewHandler(service, logger, &testOperationsService{})

	req, err := http.NewRequest("POST", "/", strings.NewReader(`{}`))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(authorizations.NewContext(ctx, &authorizations.Token{
		UserID: "user",
	}))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := fmt.Sprintf(`{"message":"%s"}`, errEmptyURL)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_Handler__Post_create_url_invalid(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	db := createTestDB(ctx, t)
	crawler := &testCrawler{}
	service := NewService(db, crawler)
	handler := NewHandler(service, logger, &testOperationsService{})

	req, err := http.NewRequest("POST", "/", strings.NewReader(`{"url": ":? ::///example.org"}`))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(authorizations.NewContext(ctx, &authorizations.Token{
		UserID: "user",
	}))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	expected := fmt.Sprintf(`{"message":"%s"}`, errInvalidURL)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_Handler__Post_create_url_not_found(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	db := createTestDB(ctx, t)
	crawler := &testCrawler{}
	service := NewService(db, crawler)
	handler := NewHandler(service, logger, &testOperationsService{})

	req, err := http.NewRequest("POST", "/", strings.NewReader(`{"url": "https://example.org"}`))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(authorizations.NewContext(ctx, &authorizations.Token{
		UserID: "user",
	}))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	returned := &operations.Operation{}
	if err := json.Unmarshal(rr.Body.Bytes(), &returned); err != nil {
		t.Fatalf("failed to unmarshal response body: %s", err)
	}

	if returned.ID == "" {
		t.Errorf("id must not be empty")
	}

	if returned.Result.Error == nil {
		t.Fatal("expected error")
	}

	if err := fmt.Sprint(returned.Result.Error.Message); err != errFailedToDownloadFeed.Error() {
		t.Fatalf("expected %s, got %s", errFailedToDownloadFeed, err)
	}
}

func Test_Handler__Post_create_success(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	db := createTestDB(ctx, t)
	crawler := &testCrawler{}
	crawler = crawler.With("https://example.org", feedData)
	service := NewService(db, crawler)
	handler := NewHandler(service, logger, &testOperationsService{})

	req, err := http.NewRequest("POST", "/", strings.NewReader(`{"url": "https://example.org"}`))
	if err != nil {
		t.Fatal(err)
	}

	req = req.WithContext(authorizations.NewContext(ctx, &authorizations.Token{
		UserID: "user",
	}))
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	returned := &operations.Operation{}
	if err := json.Unmarshal(rr.Body.Bytes(), &returned); err != nil {
		t.Fatalf("failed to unmarshal response body: %s", err)
	}

	if returned.ID == "" {
		t.Errorf("id must not be empty")
	}

	if returned.Result.Error != nil {
		t.Errorf("expected null error, got %+v", returned.Result.Error)
	}

	if returned.Result.Response == nil {
		t.Errorf("expected not null response")
	}

	title := fmt.Sprint(returned.Result.Response.(map[string]interface{})["title"])
	if title != "Sample Feed" {
		t.Errorf("expected title %s, got %s", "Sample Feed", title)
	}
}

type testOperationsService struct {
	Error error
}

func (tos *testOperationsService) Create(ctx context.Context, userID string, task operations.Task) (*operations.Operation, error) {
	if tos.Error != nil {
		return nil, tos.Error
	}

	operation := operations.New(userID)
	opChan := make(chan *operations.Operation, 1)
	if err := task(ctx, operation, opChan); err != nil {
		return nil, err
	}
	return <-opChan, nil
}