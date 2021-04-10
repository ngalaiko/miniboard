package web

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path"
)

type logger interface {
	Debug(string, ...interface{})
}

//nolint: gochecknoglobals
//go:embed static/*
var staticFS embed.FS

// Config contains web configuration.
type Config struct {
	FS bool `yaml:"fs"`
}

// Handler serves web ui.
type Handler struct {
	static fs.FS
}

// NewHandler return web handler.
func NewHandler(cfg *Config, log logger) *Handler {
	var static fs.FS
	if cfg.FS {
		static = os.DirFS("web/static")
		log.Debug("serving from os filesystem")
	} else {
		static = &stripPrefix{"static", staticFS}
		log.Debug("serving from embedded filesystem")
	}

	return &Handler{
		static: static,
	}
}

// Static returns handler that serves static files.
func (h *Handler) Static() http.HandlerFunc {
	return http.FileServer(http.FS(h.static)).ServeHTTP
}

type stripPrefix struct {
	prefix string
	fs     fs.FS
}

func (sp *stripPrefix) Open(name string) (fs.File, error) {
	return sp.fs.Open(path.Join(sp.prefix, name))
}
