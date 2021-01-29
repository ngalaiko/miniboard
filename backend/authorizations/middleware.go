package authorizations

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ngalaiko/miniboard/backend/httpx"
)

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
func Authenticate(jwtService jwtService, config *Config, logger errorLogger) httpx.Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tok, refreshable, err := getToken(r)
			if err != nil {
				httpx.Error(w, logger, err, http.StatusUnauthorized)
				return
			}

			rawToken := tok
			if rawToken == "" {
				httpx.Error(w, logger, errNoToken, http.StatusUnauthorized)
				return
			}

			token, err := jwtService.Verify(r.Context(), rawToken)
			switch {
			case err == nil:
				r = r.WithContext(NewContext(r.Context(), token))
			case errors.Is(err, errTokenExpired) && refreshable:
				token, err := jwtService.NewToken(r.Context(), token.UserID)
				if err != nil {
					logger.Error("failed to renew a new token: %s", err)
					httpx.InternalError(w, logger)
					return
				}

				setCookie(w, config, token)
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

func getToken(r *http.Request) (string, bool, error) {
	if cookie, err := r.Cookie(cookieName); err == nil {
		return cookie.Value, true, nil
	}

	token, err := getTokenFromHeader(r)
	if err != nil {
		return "", false, err
	}

	return token, false, nil
}

func getTokenFromHeader(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if token == "" {
		return "", errNoToken
	}

	tokenParts := strings.Split(token, " ")
	if len(tokenParts) != 2 {
		return "", errInvalidTokenFormat
	}

	tokenType := tokenParts[0]
	if strings.ToLower(tokenType) != "bearer" {
		return "", errInvalidTokenFormat
	}

	return tokenParts[1], nil
}