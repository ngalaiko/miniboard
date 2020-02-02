package api

import (
	"net/http"
	"regexp"
	"time"

	"github.com/sirupsen/logrus"
	"miniboard.app/images"
	"miniboard.app/jwt"
)

var imageRegExp = regexp.MustCompile("users/.+/articles/.+/images/.+")

func httpHandler(webHandler http.Handler, jwtService *jwt.Service, images *images.Service) http.Handler {
	imagesHandler := images.Handler()

	mux := http.NewServeMux()
	mux.Handle("/api/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}))
	mux.Handle("/logout", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     authCookie,
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
	}))
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if imageRegExp.MatchString(r.RequestURI) {
			imagesHandler.ServeHTTP(w, r)
			return
		}

		webHandler.ServeHTTP(w, r)
	}))

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
