package web

import (
	"net/http"
)

// Handler returns http handler for the UI.
func Handler() http.Handler {
	return http.FileServer(http.Dir("web/src"))
}
