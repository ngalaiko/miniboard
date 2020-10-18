package api

import (
	"net/http"
)

// withRecover adds recovery.
func withRecover(logger logger, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("%s", r)
			}
		}()

		h.ServeHTTP(w, r)
	})
}
