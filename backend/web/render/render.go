package render

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"sort"
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
	funcMap = map[string]interface{}{
		"timeformat": func(t *time.Time) string {
			return t.Format(time.RFC3339)
		},
	}
)

type Templates struct {
	root func() *template.Template
}

func Load(fs bool) *Templates {
	if fs {
		return loadEveryTime()
	}
	return loadOnce()
}

func loadEveryTime() *Templates {
	files := os.DirFS("web/render")
	return &Templates{
		root: func() *template.Template {
			return template.Must(readFiles(files))
		},
	}
}

func loadOnce() *Templates {
	root := template.Must(readFiles(files))
	return &Templates{
		root: func() *template.Template {
			return root
		},
	}
}

func readFiles(files fs.FS) (*template.Template, error) {
	root := template.New("")
	return root, fs.WalkDir(files, "files", func(path string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

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
	})
}

func (t *Templates) UsersPage(w io.Writer, i *items.UserItem, ii []*items.UserItem, tt []*tags.Tag, ss []*subscriptions.UserSubscription) error {
	type tagSubscriptions struct {
		Tag           *tags.Tag
		Subscriptions []*subscriptions.UserSubscription
	}

	subscriptionsByTagID := map[string][]*subscriptions.UserSubscription{}
	noTagSubscriptions := []*subscriptions.UserSubscription{}
	for _, s := range ss {
		if len(s.TagIDs) == 0 {
			noTagSubscriptions = append(noTagSubscriptions, s)
		}
		for _, tagID := range s.TagIDs {
			subscriptionsByTagID[tagID] = append(subscriptionsByTagID[tagID], s)
		}
	}
	tagsByTagID := map[string]*tags.Tag{}
	for _, tag := range tt {
		tagsByTagID[tag.ID] = tag
	}
	tagsSubscriptions := []*tagSubscriptions{}
	for tagID, ss := range subscriptionsByTagID {
		tagsSubscriptions = append(tagsSubscriptions, &tagSubscriptions{
			Tag:           tagsByTagID[tagID],
			Subscriptions: ss,
		})
	}
	sort.Slice(tagsSubscriptions, func(i, j int) bool {
		return tagsSubscriptions[i].Tag.Created.Before(tagsSubscriptions[j].Tag.Created)
	})
	return t.root().ExecuteTemplate(w, "files/users/index.html", map[string]interface{}{
		"Item":          i,
		"Items":         ii,
		"Tags":          tagsSubscriptions,
		"Subscriptions": noTagSubscriptions,
	})
}

func (t *Templates) SignupPage(w io.Writer, err error) error {
	return t.root().ExecuteTemplate(w, "files/signup/index.html", map[string]interface{}{
		"Error": err,
	})
}

func (t *Templates) LoginPage(w io.Writer, err error) error {
	return t.root().ExecuteTemplate(w, "files/login/index.html", map[string]interface{}{
		"Error": err,
	})
}

func (t *Templates) Reader(w io.Writer, item *items.Item) error {
	return t.root().ExecuteTemplate(w, "files/components/reader.html", item)
}

func (t *Templates) Item(w io.Writer, item *items.UserItem) error {
	return t.root().ExecuteTemplate(w, "files/components/item.html", item)
}

func (t *Templates) Subscription(w io.Writer, subscription *subscriptions.UserSubscription) error {
	return t.root().ExecuteTemplate(w, "files/components/subscription.html", subscription)
}

func (t *Templates) Tag(w io.Writer, tag *tags.Tag, ss []*subscriptions.UserSubscription) error {
	return t.root().ExecuteTemplate(w, "files/components/tag.html", map[string]interface{}{
		"Tag":           tag,
		"Subscriptions": ss,
	})
}
