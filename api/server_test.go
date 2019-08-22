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
	"miniboard.app/proto/users/authentications/v1"
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
			username := "test user"
			password := "password"

			resp, err := http.DefaultClient.Do(postJSON(
				t,
				server.URL+"/api/v1/users",
				map[string]interface{}{
					"username": username,
					"password": password,
				},
			))

			t.Run("It should return new user", func(t *testing.T) {
				assert.NoError(t, err)
				assert.Equal(t, resp.StatusCode, http.StatusOK)

				user := &users.User{}
				parseResponse(t, resp, user)

				assert.Equal(t, user.Username, username)

				t.Run("When creating an authorization", func(t *testing.T) {
					resp, err = http.DefaultClient.Do(postJSON(t,
						fmt.Sprintf("%s/api/v1/%s/authentications", server.URL, user.Name),
						map[string]interface{}{
							"password": "password",
						},
					))
					assert.NoError(t, err)
					assert.Equal(t, resp.StatusCode, http.StatusOK)

					authentication := &authentications.Authentication{}
					parseResponse(t, resp, authentication)

					assert.NotEmpty(t, authentication.Token)
					assert.Equal(t, authentication.Type, "Bearer")

					t.Run("When getting the user with the token", func(t *testing.T) {
						resp, err := http.DefaultClient.Do(getAuth(t,
							fmt.Sprintf("%s/api/v1/%s", server.URL, user.Name),
							authentication,
						))
						t.Run("It should return the user", func(t *testing.T) {
							assert.NoError(t, err)
							assert.Equal(t, resp.StatusCode, http.StatusOK)

							got := &users.User{}
							parseResponse(t, resp, got)

							assert.Equal(t, got.Username, username)
						})
					})
					// TODO: get another user
				})
			})
		})
	})
}

func parseResponse(t *testing.T, resp *http.Response, dst interface{}) {
	raw, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.NoError(t, json.Unmarshal(raw, dst))
}

func getAuth(t *testing.T, url string, auth *authentications.Authentication) *http.Request {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.NoError(t, err)

	req.Header.Add("Authorization", fmt.Sprintf("%s %s", auth.Type, auth.Token))

	return req
}

func postJSON(t *testing.T, url string, body interface{}) *http.Request {
	data, err := json.Marshal(body)
	assert.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	assert.NoError(t, err)

	req.Header.Add("Content-Type", "application/json")
	return req
}

func testDB(ctx context.Context, t *testing.T) storage.DB {
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
