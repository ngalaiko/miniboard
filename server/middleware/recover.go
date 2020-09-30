package middleware

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// WithRecover adds recovery.
func WithRecover(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log("panic").WithFields(logrus.Fields{
					"message": r,
				}).Error()
			}
		}()

		h.ServeHTTP(w, r)
	})
}
