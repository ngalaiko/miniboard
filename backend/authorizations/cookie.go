package authorizations

import (
	"net/http"
)

const cookieName = "authorization"

func setCookie(w http.ResponseWriter, config *Config, token *Token) {
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    token.Token,
		Expires:  token.ExpiresAt,
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
