package httpx

import (
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter

	StatusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func withAccessLogs(logger logger) Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wl := newLoggingResponseWriter(w)

			handler.ServeHTTP(wl, r)

			logger.Info("%s %d %s %s", r.Method, wl.StatusCode, r.RequestURI, time.Since(start))
		})
	}
}
