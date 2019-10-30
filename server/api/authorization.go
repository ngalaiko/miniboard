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
	authDuration = 28 * 24 * time.Hour
)

func homepageRedirect(next http.Handler, jwtService *jwt.Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subject, _ := getSubjectFromCookie(r, jwtService)
		if subject == nil || r.URL.Path != "/" {
			next.ServeHTTP(w, r)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("%s%s%s",
			r.URL.Scheme, r.URL.Host, subject,
		), http.StatusTemporaryRedirect)
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
		http.MethodPost: {
			regexp.MustCompile(`^\/api\/v1\/codes$`),
		},
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		subject, errMessage := getSubject(w, r, jwtService)

		errorCode := http.StatusUnauthorized

		switch errMessage {
		case nil:
			if strings.HasPrefix(r.URL.Path, fmt.Sprintf("/api/v1/%s", subject)) {
				h.ServeHTTP(w, r)
				return
			}

			errorCode = http.StatusForbidden
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

func getSubject(w http.ResponseWriter, r *http.Request, jwtService *jwt.Service) (*resource.Name, []byte) {
	subject, errMessage := getSubjectFromCookie(r, jwtService)
	if errMessage == nil {
		return subject, nil
	}

	authorizationCodes := r.URL.Query()["authorization_code"]
	if authorizationCodes == nil {
		return nil, errMessage
	}

	authorizationCode := authorizationCodes[0]

	subject, err := jwtService.Validate(r.Context(), authorizationCode, "authorization_code")
	if err != nil {
		return nil, []byte(fmt.Sprintf(`{"code":"16","error":"authorization code not valid","message":"%s"}`, err))
	}

	accessToken, err := jwtService.NewToken(subject, authDuration, "access_token")
	if err != nil {
		log("auth").Errorf("failed to issue new token: %s", err)
		return nil, []byte(`{"code":"16","error":"internal error","message":"something went wrong"}`)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     authCookie,
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(authDuration),
		HttpOnly: true,
		Secure:   r.TLS != nil,
	})

	return subject, nil
}

func getSubjectFromCookie(r *http.Request, jwtService *jwt.Service) (*resource.Name, []byte) {
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
