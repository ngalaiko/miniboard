package templates

import (
	"embed"
	"io"
	"text/template"

	"github.com/ngalaiko/miniboard/backend/items"
)

//nolint: gochecknoglobals
var (
	//go:embed files/*
	fs   embed.FS
	tmpl *template.Template
)

//nolint: gochecknoinits
func init() {
	tmpl = template.Must(template.ParseFS(fs, "**/*"))
}

func Reader(w io.Writer, item *items.Item) error {
	return tmpl.ExecuteTemplate(w, "reader.tmpl", item)
}
