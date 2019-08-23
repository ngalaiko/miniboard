package articles // "miniboard.app/api/articles"

import (
	"context"
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
