package authorizations

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ngalaiko/miniboard/backend/httpx"
)

type jwtValidator interface {
	Verify(context.Context, string) (*Token, error)
}

type errorLogger interface {
	Error(string, ...interface{})
}

// Known errors.
var (
	errNoToken            = fmt.Errorf("authorization token not found")
	errInvalidTokenFormat = fmt.Errorf("invalid auth token format")
)

// Authenticate is a http middleware that validates request Authorization token
// and populates a request context with the user id.
func Authenticate(jwtService jwtValidator, logger errorLogger) httpx.Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tok := getToken(r)
			if tok == "" {
				httpx.Error(w, logger, errNoToken, http.StatusUnauthorized)
				return
			}

			tokenParts := strings.Split(tok, " ")
			if len(tokenParts) != 2 {
				httpx.Error(w, logger, errInvalidTokenFormat, http.StatusUnauthorized)
				return
			}

			tokenType := tokenParts[0]
			if strings.ToLower(tokenType) != "bearer" {
				httpx.Error(w, logger, errInvalidTokenFormat, http.StatusUnauthorized)
				return
			}

			rawToken := tokenParts[1]
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

func getToken(r *http.Request) string {
	if token := r.Header.Get("Authorization"); token != "" {
		return token
	}

	if cookie, err := r.Cookie(cookieName); err == nil {
		return cookie.Value
	}

	return ""
}
