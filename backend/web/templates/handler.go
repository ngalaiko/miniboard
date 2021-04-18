package templates

import (
	"io/fs"
	"net/http"
	"text/template"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/httpx"
)

type logger interface {
	Error(string, ...interface{})
}

func NewHandler(log logger, itemsService itemsService, tagsService tagsService, subscriptionsService subscriptionsService) http.HandlerFunc {
	templates := template.Must(parseTemplates(files))
	return func(w http.ResponseWriter, r *http.Request) {
		token, auth := authorizations.FromContext(r.Context())
		if !auth {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		switch r.URL.Path {
		case "/users", "/users/", "/users/index.html":
			data, err := loadUsersData(r, token.UserID, itemsService, tagsService, subscriptionsService)
			if err != nil {
				log.Error("failed to load users data: %s", err)
				httpx.InternalError(w, log)
				return
			}
			if err := templates.ExecuteTemplate(w, "files/users/index.html", data); err != nil {
				log.Error("failed to render users page: %s", err)
				httpx.InternalError(w, log)
				return
			}
		}
	}
}

func parseTemplates(files fs.FS) (*template.Template, error) {
	root := template.New("")
	return root, fs.WalkDir(files, "files", func(path string, d fs.DirEntry, err error) error {
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
