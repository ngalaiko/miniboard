package sources

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	articles "miniboard.app/proto/users/articles/v1"
	sources "miniboard.app/proto/users/sources/v1"
)

func Test_sources(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("With sources service", func(t *testing.T) {
		articles := &mockArticles{}
		service := New(articles)

		t.Run("When creating a source", func(t *testing.T) {
			source, err := service.CreateSource(ctx, &sources.CreateSourceRequest{
				Source: &sources.Source{
					Url: "http://example.com",
				},
			})
			assert.NoError(t, err)
			assert.Equal(t, "http://example.com", source.Url)

			t.Run("Should create an article", func(t *testing.T) {
				assert.Equal(t, len(articles.articles), 1)
			})
		})
	})
}

type mockArticles struct {
	articles []*articles.Article
}

func (s *mockArticles) CreateArticle(ctx context.Context, request *articles.CreateArticleRequest) (*articles.Article, error) {
	s.articles = append(s.articles, request.Article)
	return &articles.Article{}, nil
}

func (s *mockArticles) ListArticles(ctx context.Context, request *articles.ListArticlesRequest) (*articles.ListArticlesResponse, error) {
	return nil, nil
}

func (s *mockArticles) UpdateArticle(ctx context.Context, request *articles.UpdateArticleRequest) (*articles.Article, error) {
	return nil, nil
}

func (s *mockArticles) GetArticle(ctx context.Context, request *articles.GetArticleRequest) (*articles.Article, error) {
	return nil, nil
}

func (s *mockArticles) DeleteArticle(ctx context.Context, request *articles.DeleteArticleRequest) (*empty.Empty, error) {
	return nil, nil
}
