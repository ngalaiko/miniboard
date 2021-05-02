package templates

import (
	"embed"
	"io"
	"io/fs"
	"text/template"
	"time"

	"github.com/ngalaiko/miniboard/backend/items"
	"github.com/ngalaiko/miniboard/backend/subscriptions"
	"github.com/ngalaiko/miniboard/backend/tags"
)

//nolint: gochecknoglobals
var (
	//go:embed files
	files   embed.FS
	root    = template.New("")
	funcMap = map[string]interface{}{
		"timeformat": func(t *time.Time) string {
			return t.Format(time.RFC3339)
		},
	}
)

//nolint: gochecknoinits
func init() {
	if err := fs.WalkDir(files, "files", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}

		content, err := fs.ReadFile(files, path)
		if err != nil {
			return err
		}
		if _, err := root.New(path).Funcs(funcMap).Parse(string(content)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}
}

func SignupPage(w io.Writer, err error) error {
	return root.ExecuteTemplate(w, "files/signup/index.html", map[string]interface{}{
		"Error": err,
	})
}

func LoginPage(w io.Writer, err error) error {
	return root.ExecuteTemplate(w, "files/login/index.html", map[string]interface{}{
		"Error": err,
	})
}

func Reader(w io.Writer, item *items.Item) error {
	return root.ExecuteTemplate(w, "files/components/reader.html", item)
}

func Item(w io.Writer, item *items.UserItem) error {
	return root.ExecuteTemplate(w, "files/components/item.html", item)
}

func Subscription(w io.Writer, subscription *subscriptions.UserSubscription) error {
	return root.ExecuteTemplate(w, "files/components/subscription.html", subscription)
}

func Tag(w io.Writer, tag *tags.Tag) error {
	return root.ExecuteTemplate(w, "files/components/tag.html", tag)
}
