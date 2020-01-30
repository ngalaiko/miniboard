package sources

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	articles "miniboard.app/proto/users/articles/v1"
	sources "miniboard.app/proto/users/sources/v1"
	"miniboard.app/reader"
)

type articlesService interface {
	articles.ArticlesServiceServer
	CreateArticle(context.Context, io.Reader, *url.URL) (*articles.Article, error)
}

// Service allows to add new article's sources.
// For example, a single article, or an RSS feed.
type Service struct {
	articlesService articlesService
	client          reader.GetClient
}

// New returns new sources instance.
func New(articlesService articlesService) *Service {
	return &Service{
		articlesService: articlesService,
		client:          &http.Client{},
	}
}

// CreateSource creates a new source.
func (s *Service) CreateSource(ctx context.Context, request *sources.CreateSourceRequest) (*sources.Source, error) {
	url, err := url.ParseRequestURI(request.Source.Url)
	if err != nil {
		return nil, grpc.Errorf(codes.InvalidArgument, "url is invalid")
	}

	resp, err := s.client.Get(request.Source.Url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch url")
	}
	defer resp.Body.Close()

	article, err := s.articlesService.CreateArticle(ctx, resp.Body, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create article from source")
	}
	request.Source.Name = article.Name
	return request.Source, nil
}
