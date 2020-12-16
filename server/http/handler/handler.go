package handler

import (
	"net/http"
	"path"
	"strings"
)

// Handler routes and handles http requests.
type Handler struct {
	routes map[string]http.Handler
}

// New returns a new empty handler.
func New() *Handler {
	return &Handler{}
}

// Route adds a handler for requests that start with _prefix_.
func (h *Handler) Route(prefix string, handler http.Handler) *Handler {
	h.routes[prefix] = handler
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var head string
	head, r.URL.Path = shiftPath(r.URL.Path)

	if route, exists := h.routes[head]; exists {
		route.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

// shiftPath splits off the first component of p, which will be cleaned of
// relative components before processing. head will never contain a slash and
// tail will always be a rooted path without trailing slash.
func shiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}
