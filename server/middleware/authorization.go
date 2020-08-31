package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ngalaiko/miniboard/server/api/actor"
	"github.com/ngalaiko/miniboard/server/jwt"
)

// AuthCookie is the name of authorization cookie
const AuthCookie = "auth"

// Authorized adds authorization check.
func Authorized(handler http.Handler, jwtService *jwt.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, cookie := range r.Cookies() {
			if cookie.Name != AuthCookie {
				continue
			}

			subject, err := jwtService.Validate(r.Context(), cookie.Value, "access_token")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"invalid auth token"}`))
				return
			}

			if r.URL.Path != "/api/v1/users/me" && !strings.HasPrefix(r.URL.Path, fmt.Sprintf("/api/v1/%s", subject)) {
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte(`{"error":"you are not allowed to access the resource"}`))
				return
			}

			ctx := actor.NewContext(r.Context(), subject)
			handler.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"auth cookie missing"}`))
	})
}
