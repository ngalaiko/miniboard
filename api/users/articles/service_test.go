package articles // "miniboard.app/api/articles"

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/proto/users/articles/v1"
	"miniboard.app/storage"
	"miniboard.app/storage/bolt"
	"miniboard.app/storage/resource"
)

func Test_articles(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("With articles service", func(t *testing.T) {
		service := New(testDB(ctx, t))

		t.Run("When creating an article with invalid url", func(t *testing.T) {
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
			})
			t.Run("When getting the article", func(t *testing.T) {
				resp, err := service.GetArticle(ctx, &articles.GetArticleRequest{
					Name: resp.Name,
				})
				t.Run("It should be returned", func(t *testing.T) {
					assert.NoError(t, err)
					assert.NotEmpty(t, resp.Name)
					assert.Equal(t, resp.Url, "http://localhost")
				})
			})
		})

		t.Run("When adding a few articles", func(t *testing.T) {
			parent := resource.NewName("users", "test")
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
