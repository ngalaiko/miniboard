package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/ngalaiko/miniboard/server/httpx"
)

type jwtValidator interface {
	Verify(context.Context, string) (*Token, error)
}

type errorLogger interface {
	Error(string, ...interface{})
}

// Known errors.
var (
	errNoToken = fmt.Errorf("authorization token not found")
)

// Authenticate is a http middleware that validates request Authorization token
// and populates a request context with the user id.
func Authenticate(jwtService jwtValidator, logger errorLogger) httpx.Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawToken := r.Header.Get("Authorization")
			if rawToken == "" {
				httpx.Error(w, logger, errNoToken, http.StatusUnauthorized)
				return
			}

			token, err := jwtService.Verify(r.Context(), rawToken)
			switch {
			case err == nil:
				r = r.WithContext(NewContext(r.Context(), token))
			case errors.Is(err, errInvalidToken):
				httpx.Error(w, logger, errInvalidToken, http.StatusUnauthorized)
				return
			default:
				logger.Error("failed to validate token: %s", err)
				httpx.InternalError(w, logger)
				return
			}

			handler.ServeHTTP(w, r)
		})
	}
}
