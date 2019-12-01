package api

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/andybalholm/brotli"
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
	handler = withCompression(handler)
	handler = withAccessLogs(handler)

	return handler
}

type compressedResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w compressedResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func withCompression(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch true {
		case strings.Contains(r.Header.Get("Accept-Encoding"), "br"):
			w.Header().Set("Content-Encoding", "br")
			br := brotli.NewWriter(w)
			defer br.Close()
			crw := compressedResponseWriter{Writer: br, ResponseWriter: w}
			h.ServeHTTP(crw, r)
		case strings.Contains(r.Header.Get("Accept-Encoding"), "gzip"):
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			crw := compressedResponseWriter{Writer: gz, ResponseWriter: w}
			h.ServeHTTP(crw, r)
		default:
			h.ServeHTTP(w, r)
			return
		}
	})
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
