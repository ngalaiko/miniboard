package static

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"path"
)

//nolint: gochecknoglobals
//go:embed files/*
var staticFS embed.FS

// NewHandler return web handler.
func NewHandler(useFS bool) http.HandlerFunc {
	var static fs.FS
	if useFS {
		static = os.DirFS("web/static/files")
	} else {
		static = &stripPrefix{"files", staticFS}
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
