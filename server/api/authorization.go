package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"miniboard.app/jwt"
)

var (
	errInvalidType   = []byte(`{"code":"16",error":"invalid Authorization type","message":"invalid Authorization type"}`)
	errMissingHeader = []byte(`{"code":"16","error":"authorization header missing","message":"authorization header missing"}`)
	errInvalidHeder  = []byte(`{"code":"16","error":"invalid Authorization header","message":"invalid Authorization header"}`)
	errForbidden     = []byte(`{"code":"7","error":"you don't have access to the resource","message":"you don't have access to the resource"}`)
)

func withAuthorization(h http.Handler, jwtService *jwt.Service) http.Handler {
	whitelist := map[string][]*regexp.Regexp{
		http.MethodPost: {
			regexp.MustCompile(`^\/api\/v1\/authorizations\/codes$`),
			regexp.MustCompile(`^\/api\/v1\/authorizations$`),
			regexp.MustCompile(`^\/api\/v1\/users$`),
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		list, ok := whitelist[r.Method]
		if ok {
			for _, whitelisted := range list {
				if whitelisted.MatchString(r.URL.Path) {
					h.ServeHTTP(w, r)
					return
				}
			}
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(errMissingHeader)
			return
		}

		parts := strings.Split(auth, " ")
		if len(parts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(errInvalidHeder)
			return
		}

		if strings.ToLower(parts[0]) != "bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write(errInvalidType)
			return
		}

		subject, err := jwtService.Validate(parts[1], "access")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf(`{"error":"invalid Authorization token","message":"%s"}`, err)))
			return
		}

		if !strings.HasPrefix(r.URL.Path, fmt.Sprintf("/api/v1/%s", subject)) {
			w.WriteHeader(http.StatusForbidden)
			w.Write(errForbidden)
			return
		}

		h.ServeHTTP(w, r)
	})
}
