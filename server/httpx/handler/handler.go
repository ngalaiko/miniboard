package handler

import (
	"net/http"
	"strings"
)

type route struct {
	Prefix  string
	Handler http.Handler
}

// Handler routes and handles http requests.
type Handler struct {
	routes []*route
}

// New returns a new empty handler.
func New() *Handler {
	return &Handler{}
}

// Route adds a handler for requests that start with _prefix_.
func (h *Handler) Route(prefix string, handler http.Handler) *Handler {
	h.routes = append(h.routes, &route{
		Prefix:  prefix,
		Handler: handler,
	})
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasSuffix(r.URL.Path, "/") {
		r.URL.Path += "/"
	}
	for _, route := range h.routes {
		if !strings.HasPrefix(r.URL.Path, route.Prefix) {
			continue
		}
		r.URL.Path = strings.TrimPrefix(r.URL.Path, route.Prefix)
		route.Handler.ServeHTTP(w, r)
		return
	}
	http.NotFound(w, r)
}
