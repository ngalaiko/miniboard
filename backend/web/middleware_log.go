package web

import (
	"bufio"
	"net"
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

func (lrw *loggingResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return lrw.ResponseWriter.(http.Hijacker).Hijack()
}

// Log returns http logging middleware.
func Log(logger logger) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wl := newLoggingResponseWriter(w)

			handler.ServeHTTP(wl, r)

			logger.Info("%s %d %s %s", r.Method, wl.StatusCode, r.RequestURI, time.Since(start))
		})
	}
}
