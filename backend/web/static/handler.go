package static

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
//go:embed files/*
var staticFS embed.FS

// NewHandler return web handler.
func NewHandler(useFS bool, log logger) http.HandlerFunc {
	var static fs.FS
	if useFS {
		static = os.DirFS("web/static/files")
		log.Debug("serving from os filesystem")
	} else {
		static = &stripPrefix{"files", staticFS}
		log.Debug("serving from embedded filesystem")
	}
	return http.FileServer(http.FS(static)).ServeHTTP
}

type stripPrefix struct {
	prefix string
	fs     fs.FS
}

func (sp *stripPrefix) Open(name string) (fs.File, error) {
	return sp.fs.Open(path.Join(sp.prefix, name))
}
