package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"miniboard.app/jwt"
	"miniboard.app/storage/resource"
)

const (
	authCookie   = "auth"
	authDuration = 3 * time.Hour
)

func exchangeAuthCode(h http.Handler, jwtService *jwt.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationCodes := r.URL.Query()["authorization_code"]
		if authorizationCodes == nil {
			h.ServeHTTP(w, r)
			return
		}

		authorizationCode := authorizationCodes[0]

		subject, err := jwtService.Validate(r.Context(), authorizationCode, "authorization_code")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(fmt.Sprintf(`{"code":"16","error":"authorization code not valid","message":"%s"}`, err)))
			return
		}

		accessToken, err := jwtService.NewToken(subject, authDuration, "access_token")
		if err != nil {
			log("auth").Errorf("failed to issue new token: %s", err)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"code":"16","error":"internal error","message":"something went wrong"}`))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     authCookie,
			Value:    accessToken,
			Path:     "/",
			Expires:  time.Now().Add(authDuration),
			HttpOnly: true,
			Secure:   r.TLS != nil,
		})

		http.Redirect(w, r, fmt.Sprintf("%s/%s", r.URL.Host, subject), http.StatusTemporaryRedirect)
	})
}

func removeCookie() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:     authCookie,
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
			Secure:   r.TLS != nil,
		})
	})
}

func authorize(h http.Handler, jwtService *jwt.Service) http.Handler {
	whitelist := map[string][]*regexp.Regexp{
		http.MethodGet: {
			regexp.MustCompile(`^\/app.es5.min.js$`),
			regexp.MustCompile(`^\/favicon.ico$`),
			regexp.MustCompile(`^\/logout$`),
			regexp.MustCompile(`^\/$`),
		},
		http.MethodPost: {
			regexp.MustCompile(`^\/api\/v1\/codes$`),
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subject, errMessage := getSubject(r, jwtService)

		errorCode := http.StatusForbidden

		switch errMessage {
		case nil:
			if r.URL.Path == "/" {
				http.Redirect(w, r, fmt.Sprintf("%s/%s", r.URL.Host, subject), http.StatusTemporaryRedirect)
			}

			path := strings.TrimPrefix(r.URL.Path, "/api/v1")

			if strings.HasPrefix(path, fmt.Sprintf("/%s", subject)) {
				h.ServeHTTP(w, r)
				return
			}

			errorCode = http.StatusUnauthorized
			fallthrough
		default:
			list, ok := whitelist[r.Method]
			if ok {
				for _, whitelisted := range list {
					if whitelisted.MatchString(r.URL.Path) {
						h.ServeHTTP(w, r)
						return
					}
				}
			}

			w.WriteHeader(errorCode)
			w.Write([]byte(errMessage))
			return
		}
	})
}

func getSubject(r *http.Request, jwtService *jwt.Service) (*resource.Name, []byte) {
	authCookie, err := r.Cookie(authCookie)
	if err != nil {
		return nil, []byte(`{"code":"16","error":"authorization cookie missing","message":"authorization cookie missing"}`)
	}

	subject, err := jwtService.Validate(r.Context(), authCookie.Value, "access_token")
	if err != nil {
		return nil, []byte(fmt.Sprintf(`{"error":"invalid Authorization token","message":"%s"}`, err))
	}

	return subject, nil
}
