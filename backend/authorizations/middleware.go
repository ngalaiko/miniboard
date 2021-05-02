package authorizations

import (
	"errors"
	"fmt"
	"net/http"

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

func Authenticate(jwtService jwtService, logger errorLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(cookieName)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			token, err := jwtService.Verify(r.Context(), cookie.Value)
			switch {
			case err == nil:
				r = r.WithContext(NewContext(r.Context(), token))
			case errors.Is(err, errInvalidToken):
			case errors.Is(err, errTokenExpired):
			default:
				logger.Error("failed to validate token: %s", err)
				httpx.InternalError(w, logger)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
