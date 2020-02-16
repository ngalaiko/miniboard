package sources

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	articles "miniboard.app/proto/users/articles/v1"
	rss "miniboard.app/proto/users/rss/v1"
	sources "miniboard.app/proto/users/sources/v1"
	"miniboard.app/reader"
)

type articlesService interface {
	articles.ArticlesServiceServer
	CreateArticle(context.Context, io.Reader, *url.URL) (*articles.Article, error)
}

type rssService interface {
	CreateFeed(context.Context, io.Reader) (*rss.Feed, error)
}

// Service allows to add new article's sources.
// For example, a single article, or an RSS feed.
type Service struct {
	rssService      rssService
	articlesService articlesService
	client          reader.GetClient
}

// New returns new sources instance.
func New(articlesService articlesService, rssService rssService) *Service {
	return &Service{
		articlesService: articlesService,
		rssService:      rssService,
		client:          &http.Client{},
	}
}

// CreateSource creates a new source.
func (s *Service) CreateSource(ctx context.Context, request *sources.CreateSourceRequest) (*sources.Source, error) {
	url, err := url.ParseRequestURI(request.Source.Url)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "url is invalid")
	}

	resp, err := s.client.Get(request.Source.Url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch url")
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(ct, "text/html"):
		article, err := s.articlesService.CreateArticle(ctx, resp.Body, url)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create article from source")
		}
		request.Source.Name = article.Name
		return request.Source, nil
	case strings.HasPrefix(ct, "application/rss+xml"),
		strings.HasPrefix(ct, "application/atom+xml"),
		strings.HasPrefix(ct, "text/xm"):
		feed, err := s.rssService.CreateFeed(ctx, resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create feed from source")
		}
		request.Source.Name = feed.Name
		return request.Source, nil
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported format %s", ct)
	}
}
