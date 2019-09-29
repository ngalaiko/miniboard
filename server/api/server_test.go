package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/stretchr/testify/assert"
	"gopkg.in/square/go-jose.v2/jwt"
	"miniboard.app/email/mock"
	"miniboard.app/proto/authorizations/v1"
	"miniboard.app/proto/users/articles/v1"
	"miniboard.app/proto/users/v1"
	"miniboard.app/storage"
	"miniboard.app/storage/bolt"
)

func Test_server(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := testDB(ctx, t)
	emailClient := mock.New()

	t.Run("Given a server", func(t *testing.T) {
		server := httptest.NewServer(NewServer(ctx, db, emailClient, "localhost").httpServer.Handler)
		defer server.Close()

		t.Run("When getting a non existing user without a header", func(t *testing.T) {
			resp, err := http.Get(server.URL + "/api/v1/users/random-id")
			t.Run("It should return 401", func(t *testing.T) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
			})
		})

		t.Run("When asking for auth code", func(t *testing.T) {
			email := "user@example.com"

			resp, err := http.DefaultClient.Do(postJSON(
				t,
				server.URL+"/api/v1/authorizations/codes",
				map[string]interface{}{
					"email": email,
				},
				nil,
			))

			t.Run("It should send auth code", func(t *testing.T) {
				assert.NoError(t, err)
				assert.Equal(t, resp.StatusCode, http.StatusOK)

				msg := emailClient.LastMessage()
				assert.NotEmpty(t, msg)

				url, err := url.Parse(msg)
				assert.NoError(t, err)

				code := url.Query().Get("authorization_code")
				assert.NotEmpty(t, code)

				t.Run("When creating an authorization", func(t *testing.T) {
					resp, err = http.DefaultClient.Do(postJSON(t,
						fmt.Sprintf("%s/api/v1/authorizations", server.URL),
						map[string]interface{}{
							"authorization_code": code,
							"grant_type":         "authorization_code",
						},
						nil,
					))

					assert.NoError(t, err)
					assert.Equal(t, resp.StatusCode, http.StatusOK)

					authorization := &authorizations.Authorization{}
					assert.NoError(t, jsonpb.Unmarshal(resp.Body, authorization))

					assert.NotEmpty(t, authorization.AccessToken)
					assert.Equal(t, authorization.TokenType, "Bearer")

					token, err := jwt.ParseSigned(authorization.AccessToken)
					assert.NoError(t, err)
					claims := &jwt.Claims{}
					assert.NoError(t, token.UnsafeClaimsWithoutVerification(claims))

					username := claims.Subject

					t.Run("When getting the user with the token", func(t *testing.T) {
						resp, err := http.DefaultClient.Do(getAuth(t,
							fmt.Sprintf("%s/api/v1/%s", server.URL, username),
							authorization,
						))
						t.Run("It should return the user", func(t *testing.T) {
							assert.NoError(t, err)
							assert.Equal(t, resp.StatusCode, http.StatusOK)

							got := &users.User{}
							assert.NoError(t, jsonpb.Unmarshal(resp.Body, got))

							assert.Equal(t, got.Name, "users/"+username)
						})
					})

					t.Run("When crating an article with the token", func(t *testing.T) {
						resp, err := http.DefaultClient.Do(postJSON(t,
							fmt.Sprintf("%s/api/v1/%s/articles", server.URL, username),
							map[string]interface{}{
								"url": "http://localhost",
							},
							authorization,
						))
						t.Run("It should create the article", func(t *testing.T) {
							assert.NoError(t, err)

							article := &articles.Article{}
							assert.NoError(t, jsonpb.Unmarshal(resp.Body, article))

							assert.Equal(t, article.Url, "http://localhost")
							assert.NotEmpty(t, article.Name)

							t.Run("When getting the article with the token", func(t *testing.T) {
								resp, err = http.DefaultClient.Do(getAuth(t,
									fmt.Sprintf("%s/api/v1/%s", server.URL, article.Name),
									authorization,
								))
								t.Run("It should be returned", func(t *testing.T) {
									assert.NoError(t, err)

									a := &articles.Article{}
									assert.NoError(t, jsonpb.Unmarshal(resp.Body, a))

									assert.NotEmpty(t, a.Name)
									assert.Equal(t, a.Url, "http://localhost")
								})
							})
						})
						t.Run("When listing articles", func(t *testing.T) {
							resp, err = http.DefaultClient.Do(getAuth(t,
								fmt.Sprintf("%s/api/v1/%s/articles?page_size=1", server.URL, username),
								authorization,
							))
							t.Run("It should be in the list", func(t *testing.T) {
								assert.NoError(t, err)

								aa := &articles.ListArticlesResponse{}
								assert.NoError(t, jsonpb.Unmarshal(resp.Body, aa))

								assert.Empty(t, aa.NextPageToken)
								if assert.Len(t, aa.Articles, 1) {
									assert.NotEmpty(t, aa.Articles[0].Name)
									assert.Equal(t, aa.Articles[0].Url, "http://localhost")
								}
							})
						})
					})
				})
			})
		})
	})
}

func getAuth(t *testing.T, url string, auth *authorizations.Authorization) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)

	req.Header.Add("Authorization", fmt.Sprintf("%s %s", auth.TokenType, auth.AccessToken))

	return req
}

func postJSON(t *testing.T, url string, body interface{}, auth *authorizations.Authorization) *http.Request {
	data, err := json.Marshal(body)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	assert.NoError(t, err)

	if auth != nil {
		req.Header.Add("Authorization", fmt.Sprintf("%s %s", auth.TokenType, auth.AccessToken))
	}
	return req
}

func patchJSON(t *testing.T, url string, body interface{}, auth *authorizations.Authorization) *http.Request {
	data, err := json.Marshal(body)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(data))
	assert.NoError(t, err)

	if auth != nil {
		req.Header.Add("Authorization", fmt.Sprintf("%s %s", auth.TokenType, auth.AccessToken))
	}
	return req
}

func testDB(ctx context.Context, t *testing.T) storage.Storage {
	tmpfile, err := ioutil.TempFile("", "bolt")
	if err != nil {
		t.Fatalf("failed to create database: %s", err)
	}
	go func() {
		<-ctx.Done()
		defer os.Remove(tmpfile.Name()) // clean up
	}()

	db, err := bolt.New(ctx, tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to create database: %s", err)
	}
	return db
}
