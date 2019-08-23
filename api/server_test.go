package api // import "miniboard.app/api"

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"miniboard.app/proto/users/articles/v1"
	"miniboard.app/proto/users/authorizations/v1"
	"miniboard.app/proto/users/v1"
	"miniboard.app/storage"
	"miniboard.app/storage/bolt"
)

func Test_server(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := testDB(ctx, t)

	t.Run("Given a server", func(t *testing.T) {
		server := httptest.NewServer(NewServer(ctx, db).httpServer.Handler)
		defer server.Close()

		t.Run("When getting a non existing user without a header", func(t *testing.T) {
			resp, err := http.Get(server.URL + "/api/v1/users/random-id")
			t.Run("It should return 401", func(t *testing.T) {
				assert.NoError(t, err)
				assert.Equal(t, resp.StatusCode, http.StatusUnauthorized)
			})
		})

		t.Run("When creating a user", func(t *testing.T) {
			username := "user"
			password := "password"

			resp, err := http.DefaultClient.Do(postJSON(
				t,
				server.URL+"/api/v1/users",
				map[string]interface{}{
					"username": username,
					"password": password,
				},
				nil,
			))

			t.Run("It should return new user", func(t *testing.T) {
				assert.NoError(t, err)
				assert.Equal(t, resp.StatusCode, http.StatusOK)

				user := &users.User{}
				parseResponse(t, resp, user)

				assert.Equal(t, user.Name, "users/"+username)

				t.Run("When creating an authorization", func(t *testing.T) {
					resp, err = http.DefaultClient.Do(postJSON(t,
						fmt.Sprintf("%s/api/v1/%s/authorizations", server.URL, user.Name),
						map[string]interface{}{
							"password": "password",
						},
						nil,
					))
					assert.NoError(t, err)
					assert.Equal(t, resp.StatusCode, http.StatusOK)

					authorization := &authorizations.Authorization{}
					parseResponse(t, resp, authorization)

					assert.NotEmpty(t, authorization.Token)
					assert.Equal(t, authorization.Type, "Bearer")

					t.Run("When getting the user with the token", func(t *testing.T) {
						resp, err := http.DefaultClient.Do(getAuth(t,
							fmt.Sprintf("%s/api/v1/%s", server.URL, user.Name),
							authorization,
						))
						t.Run("It should return the user", func(t *testing.T) {
							assert.NoError(t, err)
							assert.Equal(t, resp.StatusCode, http.StatusOK)

							got := &users.User{}
							parseResponse(t, resp, got)

							assert.Equal(t, got.Name, "users/"+username)
						})
					})

					t.Run("When crating an article with the token", func(t *testing.T) {
						resp, err := http.DefaultClient.Do(postJSON(t,
							fmt.Sprintf("%s/api/v1/%s/articles", server.URL, user.Name),
							map[string]interface{}{
								"url": "http://localhost",
							},
							authorization,
						))
						t.Run("It should create the article", func(t *testing.T) {
							assert.NoError(t, err)

							article := &articles.Article{}

							parseResponse(t, resp, article)

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
									parseResponse(t, resp, &a)

									assert.NotEmpty(t, a.Name)
									assert.Equal(t, a.Url, "http://localhost")
								})
							})
						})
						t.Run("When listing articles", func(t *testing.T) {
							resp, err = http.DefaultClient.Do(getAuth(t,
								fmt.Sprintf("%s/api/v1/%s/articles?page_size=1", server.URL, user.Name),
								authorization,
							))
							t.Run("It should be in the list", func(t *testing.T) {
								assert.NoError(t, err)

								aa := struct {
									Articles      []*articles.Article `json:"articles"`
									NextPageToken string              `json:"next_page_token"`
								}{}
								parseResponse(t, resp, &aa)

								assert.Len(t, aa.Articles, 1)
								assert.Empty(t, aa.NextPageToken)
								assert.NotEmpty(t, aa.Articles[0].Name)
								assert.Equal(t, aa.Articles[0].Url, "http://localhost")
							})
						})
					})
				})
			})
		})
	})
}

func parseResponse(t *testing.T, resp *http.Response, dst interface{}) {
	raw, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	t.Log(string(raw))
	assert.NoError(t, json.Unmarshal(raw, dst))
}

func getAuth(t *testing.T, url string, auth *authorizations.Authorization) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)

	req.Header.Add("Authorization", fmt.Sprintf("%s %s", auth.Type, auth.Token))

	return req
}

func postJSON(t *testing.T, url string, body interface{}, auth *authorizations.Authorization) *http.Request {
	data, err := json.Marshal(body)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	assert.NoError(t, err)

	if auth != nil {
		req.Header.Add("Authorization", fmt.Sprintf("%s %s", auth.Type, auth.Token))
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
