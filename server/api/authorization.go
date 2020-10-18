package api

import (
	"net/http"

	"github.com/ngalaiko/miniboard/server/actor"
	"github.com/ngalaiko/miniboard/server/jwt"
)

// authCookie is the name of authorization cookie
const authCookie = "auth"

// authorized adds authorization check.
func authorized(handler http.Handler, jwtService *jwt.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, cookie := range r.Cookies() {
			if cookie.Name != authCookie {
				continue
			}

			subject, err := jwtService.Validate(r.Context(), cookie.Value, "access_token")
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = w.Write([]byte(`{"error":"invalid auth token"}`))
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
