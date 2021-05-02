package web

import (
	"errors"
	"net/http"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/httpx"
)

type errorLogger interface {
	Error(string, ...interface{})
}

func Authenticate(jwtService jwtService, log errorLogger) func(http.Handler) http.Handler {
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
				r = r.WithContext(authorizations.NewContext(r.Context(), token))
			case errors.Is(err, authorizations.ErrInvalidToken):
			case errors.Is(err, authorizations.ErrTokenExpired):
				token, err := jwtService.NewToken(r.Context(), token.UserID)
				if err != nil {
					log.Error("failed to create a new token: %s", err)
					httpx.InternalError(w, log)
					return
				}

				setCookie(w, r.TLS != nil, token)
			default:
				log.Error("failed to validate token: %s", err)
				httpx.InternalError(w, log)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
