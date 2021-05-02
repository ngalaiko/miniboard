package web

import (
	"context"
	"errors"
	"net/http"

	"github.com/ngalaiko/miniboard/backend/authorizations"
	"github.com/ngalaiko/miniboard/backend/httpx"
	"github.com/ngalaiko/miniboard/backend/users"
	"github.com/ngalaiko/miniboard/backend/web/templates"
)

type usersService interface {
	Create(context.Context, string, []byte) (*users.User, error)
	GetByUsername(context.Context, string) (*users.User, error)
}

type jwtService interface {
	NewToken(context.Context, string) (*authorizations.Token, error)
	Verify(context.Context, string) (*authorizations.Token, error)
}

func loginHandler(log logger, usersService usersService, jwtService jwtService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Error("failed to parse form: %s", err)
			httpx.InternalError(w, log)
			return
		}

		user, err := usersService.GetByUsername(r.Context(), r.Form.Get("username"))
		switch {
		case err == nil:
		case errors.Is(users.ErrNotFound, err):
			if err := templates.LoginPage(w, users.ErrNotFound); err != nil {
				log.Error("failed to render login page: %s", err)
				httpx.InternalError(w, log)
			}
			return
		default:
			log.Error("failed to get user: %s", err)
			httpx.InternalError(w, log)
			return
		}

		validatePasswordErr := user.ValidatePassword([]byte(r.Form.Get("password")))
		switch {
		case validatePasswordErr == nil:
		case errors.Is(validatePasswordErr, users.ErrInvalidPassword):
			if err := templates.LoginPage(w, users.ErrInvalidPassword); err != nil {
				log.Error("failed to render login page: %s", err)
				httpx.InternalError(w, log)
			}
			return
		default:
			log.Error("failed to validate password: %s", err)
			httpx.InternalError(w, log)
			return
		}

		token, err := jwtService.NewToken(r.Context(), user.ID)
		if err != nil {
			log.Error("failed to create a new token: %s", err)
			httpx.InternalError(w, log)
			return
		}

		setCookie(w, r.TLS != nil, token)

		http.Redirect(w, r, "/users/", http.StatusSeeOther)
	}
}
