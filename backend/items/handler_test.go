package items

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ngalaiko/miniboard/backend/authorizations"
)

func Test_Handler__Get(t *testing.T) {
	ctx := context.Background()
	logger := &testLogger{}
	db := createTestDB(ctx, t)
	service := NewService(db, &testLogger{})
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
