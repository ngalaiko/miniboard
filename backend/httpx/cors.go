package httpx

import (
	"net/http"
)

// WithCors adds cors to a a handler.
func WithCors(domains ...string) Middleware {
	allowed := map[string]bool{}
	for _, domain := range domains {
		allowed[domain] = true
	}
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if !allowed[origin] {
				return
			}

			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				return
			}

			handler.ServeHTTP(w, r)
		})
	}
}
