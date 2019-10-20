package api

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"miniboard.app/jwt"
)

const (
	authCookie   = "auth"
	authDuration = 3 * time.Hour
)

func withAuthorization(h http.Handler, jwtService *jwt.Service) http.Handler {
	whitelist := map[string][]*regexp.Regexp{
		http.MethodPost: {
			regexp.MustCompile(`^\/api\/v1\/codes$`),
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		list, ok := whitelist[r.Method]
		if ok {
			for _, whitelisted := range list {
				if whitelisted.MatchString(r.URL.Path) {
					h.ServeHTTP(w, r)
					return
				}
			}
		}

		authCookie, err := r.Cookie(authCookie)
		switch err {
		case nil:
		case http.ErrNoCookie:
			if err := setCookie(w, r, jwtService); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(`{"code":"16","error":"authorization cookie missing","message":"authorization cookie missing"}`))
				return
			}
			return
		default:
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"code":"16","error":"authorization cookie missing","message":"authorization cookie missing"}`))
			return
		}

		subject, err := jwtService.Validate(r.Context(), authCookie.Value, "access_token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf(`{"error":"invalid Authorization token","message":"%s"}`, err)))
			return
		}

		if !strings.HasPrefix(r.URL.Path, fmt.Sprintf("/api/v1/%s", subject)) {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(`{"code":"7","error":"you don't have access to the resource","message":"you don't have access to the resource"}`))
			return
		}

		h.ServeHTTP(w, r)
	})
}

func setCookie(w http.ResponseWriter, r *http.Request, jwtService *jwt.Service) error {
	authorizationCodes := r.URL.Query()["authorization_code"]
	if authorizationCodes == nil {
		return errors.New("authorization code is empty")
	}

	authorizationCode := authorizationCodes[0]

	subject, err := jwtService.Validate(r.Context(), authorizationCode, "authorization_code")
	if err != nil {
		return errors.New("authorization code is not valid")
	}

	accessToken, err := jwtService.NewToken(subject, authDuration, "access_token")
	if err != nil {
		return errors.New("faield to generate access code")
	}

	http.SetCookie(w, &http.Cookie{
		Name:     authCookie,
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(authDuration),
		HttpOnly: true,
	})

	http.Redirect(w, r, fmt.Sprintf("%s/%s", r.URL.Host, subject), http.StatusTemporaryRedirect)

	return nil
}
