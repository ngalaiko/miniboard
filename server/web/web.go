package web

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// Handler returns http handler for the UI.
func Handler(filePath string) http.Handler {
	log().Infof("filepath: %s", filePath)
	return http.FileServer(http.Dir(filePath))
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "web",
	})
}
