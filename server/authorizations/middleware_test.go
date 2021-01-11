package authorizations

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_Middleware__no_token(t *testing.T) {
	middleware := Authenticate(&testJWValidator{}, &testErrorLogger{})

	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	r, _ := http.NewRequest("GET", "/", nil)

	rr := httptest.NewRecorder()

	middleware(testHandler).ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	if handlerCalled {
		t.Errorf("handler was called")
	}

	expected := fmt.Sprintf(`{"message":"%s"}`, errNoToken)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_Middleware__internal(t *testing.T) {
	middleware := Authenticate(&testJWValidator{
		Error: fmt.Errorf("internal"),
	}, &testErrorLogger{})

	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Add("Authorization", "bearer invalid")

	rr := httptest.NewRecorder()

	middleware(testHandler).ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusInternalServerError)
	}

	if handlerCalled {
		t.Errorf("handler was called")
	}

	expected := fmt.Sprintf(`{"message":"internal server error"}`)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_Middleware__invalid_token_type(t *testing.T) {
	middleware := Authenticate(&testJWValidator{}, &testErrorLogger{})

	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Add("Authorization", "smth invalid")

	rr := httptest.NewRecorder()

	middleware(testHandler).ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	if handlerCalled {
		t.Errorf("handler was called")
	}

	expected := fmt.Sprintf(`{"message":"%s"}`, errInvalidTokenFormat)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_Middleware__invalid_token(t *testing.T) {
	middleware := Authenticate(&testJWValidator{
		Error: errInvalidToken,
	}, &testErrorLogger{})

	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Add("Authorization", "bearer invalid")

	rr := httptest.NewRecorder()

	middleware(testHandler).ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	if handlerCalled {
		t.Errorf("handler was called")
	}

	expected := fmt.Sprintf(`{"message":"%s"}`, errInvalidToken)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func Test_Middleware__valid_token(t *testing.T) {
	middleware := Authenticate(&testJWValidator{
		Token: &Token{
			Token:     "token",
			UserID:    "id",
			ExpiresAt: time.Now(),
		},
	}, &testErrorLogger{})

	handlerCalled := false
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
	})

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Add("Authorization", "bearer token")

	rr := httptest.NewRecorder()

	middleware(testHandler).ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if !handlerCalled {
		t.Errorf("handler was not called")
	}
}

type testJWValidator struct {
	Token *Token
	Error error
}

func (tjv *testJWValidator) Verify(context.Context, string) (*Token, error) {
	return tjv.Token, tjv.Error
}

type testErrorLogger struct{}

func (tel *testErrorLogger) Error(string, ...interface{}) {}
