package api

import (
	"net/http"
	"time"
)

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

// withAccessLogs adds access logging.
func withAccessLogs(logger logger, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wl := newLoggingResponseWriter(w)

		h.ServeHTTP(wl, r)

		logger.Info("%d %s %s %s", wl.statusCode, r.Method, r.URL.String(), time.Since(start))
	})
}
