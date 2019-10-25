package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/golang/protobuf/jsonpb"
	"github.com/stretchr/testify/assert"
	"miniboard.app/email/mock"
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

	jar, err := cookiejar.New(nil)
	assert.NoError(t, err)

	httpClient := &http.Client{
		Jar: jar,
	}

	t.Run("Given a server", func(t *testing.T) {
		server := httptest.NewServer(NewServer(ctx, db, emailClient, "http://localhost").httpServer.Handler)
		defer server.Close()

		t.Run("When getting a non existing user without a header", func(t *testing.T) {
			resp, err := http.Get(server.URL + "/api/v1/users/random-id")
			t.Run("It should return 401", func(t *testing.T) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusForbidden, resp.StatusCode)
			})
		})

		t.Run("When asking for auth code", func(t *testing.T) {
			email := "user@example.com"

			resp, err := httpClient.Do(postJSON(
				t,
				server.URL+"/api/v1/codes",
				map[string]interface{}{
					"email": email,
				},
			))

			t.Run("It should send auth code", func(t *testing.T) {
				assert.NoError(t, err)
				assert.Equal(t, resp.StatusCode, http.StatusOK)

				msg := emailClient.LastMessage()
				assert.NotEmpty(t, msg)

				loginURL, err := url.Parse(msg)
				assert.NoError(t, err)

				t.Run("When clicking on link from email", func(t *testing.T) {
					t.Run("It should set cookie and redirect", func(t *testing.T) {
						_, err := httpClient.Do(get(t,
							server.URL+loginURL.RequestURI(),
						))
						assert.NoError(t, err)
					})
				})

				parts := strings.Split(loginURL.Path, "/")
				username := fmt.Sprintf("%s/%s", parts[3], parts[4])

				t.Run("When getting the user with the token", func(t *testing.T) {
					resp, err := httpClient.Do(get(t,
						fmt.Sprintf("%s/api/v1/%s", server.URL, username),
					))
					t.Run("It should return the user", func(t *testing.T) {
						assert.NoError(t, err)
						assert.Equal(t, resp.StatusCode, http.StatusOK)

						got := &users.User{}
						assert.NoError(t, jsonpb.Unmarshal(resp.Body, got))

						assert.Equal(t, got.Name, username)
					})
				})

				t.Run("When crating an article with the token", func(t *testing.T) {
					resp, err := httpClient.Do(postJSON(t,
						fmt.Sprintf("%s/api/v1/%s/articles", server.URL, username),
						map[string]interface{}{
							"url": "http://localhost",
						},
					))
					t.Run("It should create the article", func(t *testing.T) {
						assert.NoError(t, err)

						article := &articles.Article{}
						assert.NoError(t, jsonpb.Unmarshal(resp.Body, article))

						assert.Equal(t, article.Url, "http://localhost")
						assert.NotEmpty(t, article.Name)

						t.Run("When getting the article with the token", func(t *testing.T) {
							resp, err = httpClient.Do(get(t,
								fmt.Sprintf("%s/api/v1/%s", server.URL, article.Name),
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
						resp, err = httpClient.Do(get(t,
							fmt.Sprintf("%s/api/v1/%s/articles?page_size=1", server.URL, username),
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
}

func get(t *testing.T, url string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)
	return req
}

func postJSON(t *testing.T, url string, body interface{}) *http.Request {
	data, err := json.Marshal(body)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	assert.NoError(t, err)

	return req
}

func patchJSON(t *testing.T, url string, body interface{}) *http.Request {
	data, err := json.Marshal(body)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(data))
	assert.NoError(t, err)

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
