package web

import (
	"net/http"

	packr "github.com/gobuffalo/packr/v2"
)

// Handler returns http handler for the UI.
func Handler() http.Handler {
	box := packr.New("src", "../../web/src")
	return http.FileServer(box)
}
