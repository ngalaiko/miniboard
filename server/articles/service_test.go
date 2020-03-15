package articles

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/api/actor"
	"miniboard.app/storage/redis"
	"miniboard.app/storage/resource"
)

func testArticle(replacement string) io.ReadCloser {
	file, err := os.Open("./testdata/test.html")
	if err != nil {
		return nil
	}

	dd, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}

	dd = bytes.Replace(dd, []byte("__RANDOM__"), []byte(replacement), 1)

	return ioutil.NopCloser(bytes.NewBuffer(dd))
}

type testClient struct{}

func (tc *testClient) Fetch(ctx context.Context, url string) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       testArticle(url),
	}, nil
}

func Test_articles(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()

	db, err := redis.New(ctx, s.Addr())
	if err != nil {
		t.Fatalf("failed to create database: %s", err)
	}

	testClient := &testClient{}

	t.Run("With articles service", func(t *testing.T) {
		service := NewService(db, testClient)

		t.Run("When creating an article", func(t *testing.T) {
			ctx = actor.NewContext(ctx, resource.NewName("users", "test"))

			url, err := url.Parse("http://localhost")
			assert.NoError(t, err)

			body := testArticle(url.String())

			resp, err := service.CreateArticle(ctx, body, url, nil)
			t.Run("It should be created", func(t *testing.T) {
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.Name)
				assert.Equal(t, resp.Url, "http://localhost")
				assert.NotEmpty(t, resp.Content)
			})
			t.Run("When adding the same article twice", func(t *testing.T) {
				body = testArticle(url.String())
				ts := time.Now().Add(time.Second)
				resp2, err := service.CreateArticle(ctx, body, url, &ts)
				t.Run("It should return the same article with the same content", func(t *testing.T) {
					if assert.NoError(t, err) {
						assert.Equal(t, resp2.Name, resp.Name)
						assert.Equal(t, resp2.ContentSha256Sum, resp.ContentSha256Sum)
					}
				})
			})
			t.Run("When adding the same article with different content", func(t *testing.T) {
				body = testArticle("new content here")
				ts := time.Now().Add(time.Hour)
				resp2, err := service.CreateArticle(ctx, body, url, &ts)

				t.Run("It should return the same article with different content", func(t *testing.T) {
					if assert.NoError(t, err) {
						assert.Equal(t, resp2.Name, resp.Name)
						assert.Equal(t, resp2.CreateTime, resp.CreateTime)
						assert.NotEqual(t, resp2.ContentSha256Sum, resp.ContentSha256Sum)
					}
				})
			})
			t.Run("When getting the article with full view", func(t *testing.T) {
				resp, err := service.GetArticle(ctx, &GetArticleRequest{
					Name: resp.Name,
					View: ArticleView_ARTICLE_VIEW_FULL,
				})
				t.Run("It should be returned", func(t *testing.T) {
					assert.NoError(t, err)
					assert.NotEmpty(t, resp.Name)
					assert.Equal(t, resp.Url, "http://localhost")
					assert.NotEmpty(t, resp.Title)
					assert.NotEmpty(t, resp.IconUrl)
					assert.NotEmpty(t, resp.CreateTime)
					assert.NotEmpty(t, resp.Content)
				})
			})
			t.Run("When getting the article with basic view", func(t *testing.T) {
				resp, err := service.GetArticle(ctx, &GetArticleRequest{
					Name: resp.Name,
				})
				t.Run("It should be returned", func(t *testing.T) {
					assert.NoError(t, err)
					assert.NotEmpty(t, resp.Name)
					assert.Equal(t, resp.Url, "http://localhost")
					assert.NotEmpty(t, resp.Title)
					assert.NotEmpty(t, resp.IconUrl)
					assert.NotEmpty(t, resp.CreateTime)
					assert.Empty(t, resp.Content)
				})
			})
			t.Run("When deleting the article", func(t *testing.T) {
				_, err = service.DeleteArticle(ctx, &DeleteArticleRequest{
					Name: resp.Name,
				})
				assert.NoError(t, err)
				t.Run("It should be deleted", func(t *testing.T) {
					_, err = service.GetArticle(ctx, &GetArticleRequest{
						Name: resp.Name,
					})
					assert.Error(t, err)

					status, ok := status.FromError(err)
					assert.True(t, ok)

					assert.Equal(t, codes.NotFound, status.Code())
				})
			})
		})

		t.Run("When adding a few articles", func(t *testing.T) {
			parent := resource.NewName("users", "test2")
			ctx = actor.NewContext(ctx, parent)
			for i := 0; i < 50; i++ {
				url, err := url.Parse(fmt.Sprintf("http://localhost.com/%d", i))
				assert.NoError(t, err)

				body := testArticle(url.String())

				resp, err := service.CreateArticle(ctx, body, url, nil)
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.Name)
				assert.Equal(t, resp.Url, fmt.Sprintf("http://localhost.com/%d", i))
			}

			t.Run("It should be possible to get then page by page", func(t *testing.T) {
				pageToken := ""
				for i := 0; i < 10; i++ {
					resp, err := service.ListArticles(ctx, &ListArticlesRequest{
						PageSize:  5,
						PageToken: pageToken,
					})
					assert.NoError(t, err)
					assert.Len(t, resp.Articles, 5)
					if i != 9 {
						assert.NotEmpty(t, resp.NextPageToken)
					} else {
						assert.Empty(t, resp.NextPageToken)
					}

					pageToken = resp.NextPageToken
				}
			})

			t.Run("It should be possible filter for article by title", func(t *testing.T) {
				resp, err := service.ListArticles(ctx, &ListArticlesRequest{
					Title:    &wrappers.StringValue{Value: "Building"},
					PageSize: 10,
				})
				assert.NoError(t, err)
				assert.Len(t, resp.Articles, 10)
			})
		})
	})
}
