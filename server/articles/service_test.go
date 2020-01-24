package articles

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/api/actor"
	"miniboard.app/images"
	"miniboard.app/proto/users/articles/v1"
	"miniboard.app/storage"
	"miniboard.app/storage/bolt"
	"miniboard.app/storage/resource"
)

type testClient struct{}

func (tc *testClient) Get(url string) (*http.Response, error) {
	file, err := os.Open("./testdata/test.html")
	if err != nil {
		return nil, err
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       file,
	}, nil
}

func Test_articles(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db := testDB(ctx, t)

	t.Run("With articles service", func(t *testing.T) {
		service := New(db, images.New(db))
		service.client = &testClient{}

		t.Run("When creating an article with invalid url", func(t *testing.T) {
			ctx = actor.NewContext(ctx, resource.NewName("users", "test"))
			resp, err := service.CreateArticle(ctx, &articles.CreateArticleRequest{
				Parent: resource.NewName("users", "test").String(),
				Article: &articles.Article{
					Url: "invalid :(",
				},
			})
			t.Run("Error should be InvalidArgument", func(t *testing.T) {
				assert.Nil(t, resp)
				assert.Error(t, err)
				status, ok := status.FromError(err)
				assert.True(t, ok)
				assert.Equal(t, status.Code(), codes.InvalidArgument)
			})
		})

		t.Run("When creating an article", func(t *testing.T) {
			ctx = actor.NewContext(ctx, resource.NewName("users", "test"))
			resp, err := service.CreateArticle(ctx, &articles.CreateArticleRequest{
				Parent: resource.NewName("users", "test1").String(),
				Article: &articles.Article{
					Url: "http://localhost",
				},
			})
			t.Run("It should be created", func(t *testing.T) {
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.Name)
				assert.Equal(t, resp.Url, "http://localhost")
				assert.NotEmpty(t, resp.Content)
			})
			t.Run("When getting the article with full view", func(t *testing.T) {
				resp, err := service.GetArticle(ctx, &articles.GetArticleRequest{
					Name: resp.Name,
					View: articles.ArticleView_ARTICLE_VIEW_FULL,
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
				resp, err := service.GetArticle(ctx, &articles.GetArticleRequest{
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
				_, err = service.DeleteArticle(ctx, &articles.DeleteArticleRequest{
					Name: resp.Name,
				})
				assert.NoError(t, err)
				t.Run("It should be deleted", func(t *testing.T) {
					_, err = service.GetArticle(ctx, &articles.GetArticleRequest{
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
			parent := resource.NewName("users", "test")
			ctx = actor.NewContext(ctx, parent)
			for i := 0; i < 50; i++ {
				resp, err := service.CreateArticle(ctx, &articles.CreateArticleRequest{
					Parent: parent.String(),
					Article: &articles.Article{
						Url: fmt.Sprintf("http://localhost.com/%d", i),
					},
				})
				assert.NotEmpty(t, resp.Name)
				assert.Equal(t, resp.Url, fmt.Sprintf("http://localhost.com/%d", i))
				assert.NoError(t, err)
			}

			t.Run("It should be possible to get then page by page", func(t *testing.T) {
				pageToken := ""
				for i := 0; i < 10; i++ {
					resp, err := service.ListArticles(ctx, &articles.ListArticlesRequest{
						Parent:    parent.String(),
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
				resp, err := service.ListArticles(ctx, &articles.ListArticlesRequest{
					Parent:   parent.String(),
					Title:    &wrappers.StringValue{Value: "Building"},
					PageSize: 10,
				})
				assert.NoError(t, err)
				assert.Len(t, resp.Articles, 10)
			})

			t.Run("It should be possible filter for article by url", func(t *testing.T) {
				resp, err := service.ListArticles(ctx, &articles.ListArticlesRequest{
					Parent:   parent.String(),
					Url:      &wrappers.StringValue{Value: "localhost"},
					PageSize: 5,
				})
				assert.NoError(t, err)
				assert.Len(t, resp.Articles, 5)
			})
		})
	})
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
