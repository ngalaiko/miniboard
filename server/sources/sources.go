package sources

import (
	"context"

	"github.com/pkg/errors"
	articles "miniboard.app/proto/users/articles/v1"
	sources "miniboard.app/proto/users/sources/v1"
)

// Service allows to add new article's sources.
// For example, a single article, or an RSS feed.
type Service struct {
	articlesService articles.ArticlesServiceServer
}

// New returns new sources instance.
func New(articlesService articles.ArticlesServiceServer) *Service {
	return &Service{
		articlesService: articlesService,
	}
}

// CreateSource creates a new source.
func (s *Service) CreateSource(ctx context.Context, request *sources.CreateSourceRequest) (*sources.Source, error) {
	article, err := s.articlesService.CreateArticle(ctx, &articles.CreateArticleRequest{
		Article: &articles.Article{
			Url: request.Source.Url,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create article from source")
	}
	request.Source.Name = article.Name
	return request.Source, nil
}
