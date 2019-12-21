package web

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// Handler returns http handler for the UI.
func Handler() http.Handler {
	fileHandler := http.FileServer(&fs{
		rootFS: http.Dir("./web"),
	})
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if pusher, ok := w.(http.Pusher); ok {
			if err := pusher.Push("/app/app.js", nil); err != nil {
				log("web").Errorf("failed to push /app/app.js: %s", err)
			}
		}

		fileHandler.ServeHTTP(w, r)
	})
}

type fs struct {
	rootFS http.FileSystem
}

// Open opens a file if it exists, or 'index.html' otherweise.
func (fs *fs) Open(name string) (http.File, error) {
	file, err := fs.rootFS.Open(name)
	if err != nil {
		return fs.rootFS.Open("app/index.html")
	}
	return file, err
}

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
