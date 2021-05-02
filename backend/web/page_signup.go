package web

import (
	"errors"
	"net/http"

	"github.com/ngalaiko/miniboard/backend/httpx"
	"github.com/ngalaiko/miniboard/backend/users"
	"github.com/ngalaiko/miniboard/backend/web/render"
)

func signupHandler(log logger, usersService usersService, jwtService jwtService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Error("failed to parse form: %s", err)
			httpx.InternalError(w, log)
			return
		}

		user, err := usersService.Create(r.Context(), r.Form.Get("username"), []byte(r.Form.Get("password")))
		switch {
		case err == nil:
			token, err := jwtService.NewToken(r.Context(), user.ID)
			if err != nil {
				log.Error("failed to create a new token: %s", err)
				httpx.InternalError(w, log)
				return
			}
			setCookie(w, r.TLS != nil, token)

			http.Redirect(w, r, "/users/", http.StatusSeeOther)
		case errors.Is(err, users.ErrAlreadyExists):
			if err := render.SignupPage(w, err); err != nil {
				log.Error("failed to render signup page: %s", err)
				httpx.InternalError(w, log)
			}
		case errors.Is(err, users.ErrUsernameEmpty),
			errors.Is(err, users.ErrPasswordEmpty):
			if err := render.SignupPage(w, errors.Unwrap(err)); err != nil {
				log.Error("failed to render signup page: %s", err)
				httpx.InternalError(w, log)
			}
		default:
			log.Error("failed to create user: %s", err)
			httpx.InternalError(w, log)
		}
	}
}
