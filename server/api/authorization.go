package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ngalaiko/miniboard/server/api/actor"
	"github.com/ngalaiko/miniboard/server/jwt"
)

const authCookie = "auth"

func authorize(handler http.Handler, jwtService *jwt.Service) http.Handler {
	whitelisted := map[string]bool{
		"/api/v1/codes":  true,
		"/api/v1/tokens": true,
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if whitelisted[r.URL.Path] {
			handler.ServeHTTP(w, r)
			return
		}

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
