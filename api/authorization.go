package api // import "miniboard.app/api"

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"miniboard.app/jwt"
)

func withAuthorization(h http.Handler, jwtService *jwt.Service) http.Handler {
	whitelist := map[string][]*regexp.Regexp{
		http.MethodPost: {
			regexp.MustCompile("/api/v1/users/*/authorizations"),
			regexp.MustCompile("/api/v1/users"),
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		list, ok := whitelist[r.Method]
		if ok {
			for _, whitelisted := range list {
				if whitelisted.MatchString(r.URL.Path) {
					h.ServeHTTP(w, r)
				}
			}
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if strings.ToLower(parts[0]) != "bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		subject, err := jwtService.Validate(parts[1])
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !strings.HasPrefix(r.URL.Path, fmt.Sprintf("/api/v1/%s", subject)) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	})
}
