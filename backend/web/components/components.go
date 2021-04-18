package components

import (
	"embed"
	"io"
	"text/template"
	"time"

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
	tmpl = template.Must(template.New("").Funcs(map[string]interface{}{
		"timeformat": func(t *time.Time) string {
			return t.Format(time.RFC3339)
		},
	}).ParseFS(fs, "**/*"))
}

func Reader(w io.Writer, item *items.Item) error {
	return tmpl.ExecuteTemplate(w, "reader.tmpl", item)
}

func Item(w io.Writer, item *items.UserItem) error {
	return tmpl.ExecuteTemplate(w, "item.tmpl", item)
}
