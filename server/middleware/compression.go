package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/andybalholm/brotli"
)

type compressedResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *compressedResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *compressedResponseWriter) Flush() {
	if fw, ok := w.Writer.(http.Flusher); ok {
		fw.Flush()
	}
}

// WithCompression adds compression.
func WithCompression(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.Header.Get("Accept-Encoding"), "br"):
			w.Header().Set("Content-Encoding", "br")
			br := brotli.NewWriter(w)
			defer br.Close()
			crw := &compressedResponseWriter{Writer: br, ResponseWriter: w}
			h.ServeHTTP(crw, r)
		case strings.Contains(r.Header.Get("Accept-Encoding"), "gzip"):
			w.Header().Set("Content-Encoding", "gzip")
			gz := gzip.NewWriter(w)
			defer gz.Close()
			crw := &compressedResponseWriter{Writer: gz, ResponseWriter: w}
			h.ServeHTTP(crw, r)
		default:
			h.ServeHTTP(w, r)
			return
		}
	})
}
