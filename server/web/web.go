package web

import (
	"net/http"
)

// Handler returns http handler for the UI.
func Handler() http.Handler {
	return http.FileServer(&fs{
		rootFS: http.Dir("./web"),
	})
}

type fs struct {
	rootFS http.FileSystem
}

// Open opens a file if it exists, or 'index.html' otherweise.
func (fs *fs) Open(name string) (http.File, error) {
	file, err := fs.rootFS.Open(name)
	if err != nil {
		return fs.rootFS.Open("index.html")
	}
	return file, err
}
