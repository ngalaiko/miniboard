package web

import (
	"net/http"
	"time"

	"github.com/ngalaiko/miniboard/backend/authorizations"
)

const cookieName = "authorization"

func removeCookie(w http.ResponseWriter, secure bool) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, cookie)
}

func setCookie(w http.ResponseWriter, secure bool, token *authorizations.Token) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    token.Token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, cookie)
}
