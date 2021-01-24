package authorizations

import (
	"net/http"
	"time"
)

const cookieName = "authorization"

var cookieLifetime = 30 * 25 * time.Hour

func setCookie(w http.ResponseWriter, config *Config, token *Token) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    token.Token,
		Expires:  time.Now().Add(cookieLifetime),
		HttpOnly: true,
		Path:     "/",
		Secure:   config.Secure,
		SameSite: http.SameSiteNoneMode,
	}
	if config.Domain != nil {
		cookie.Domain = *config.Domain
	}
	http.SetCookie(w, cookie)
}
