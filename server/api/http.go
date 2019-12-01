package api

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"miniboard.app/jwt"
)

func httpHandler(webHandler http.Handler, jwtService *jwt.Service) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/api/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}))
	mux.Handle("/logout", removeCookie())
	mux.Handle("/", homepageRedirect(webHandler, jwtService))

	handler := http.Handler(mux)
	handler = withAccessLogs(handler)
	handler = withCompression(handler)

	return handler
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func withAccessLogs(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wl := newLoggingResponseWriter(w)

		h.ServeHTTP(wl, r)

		log("access").WithFields(logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.String(),
			"ts":       start.Format(time.RFC3339),
			"duration": time.Since(start),
			"status":   wl.statusCode,
		}).Info()
	})
}
