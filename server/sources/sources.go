package sources

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	articles "miniboard.app/proto/users/articles/v1"
	feeds "miniboard.app/proto/users/feeds/v1"
	sources "miniboard.app/proto/users/sources/v1"
	"miniboard.app/reader"
)

type articlesService interface {
	articles.ArticlesServiceServer
	CreateArticle(context.Context, io.Reader, *url.URL, *time.Time) (*articles.Article, error)
}

type feedsService interface {
	CreateFeed(context.Context, io.Reader, *url.URL) (*feeds.Feed, error)
}

// Service allows to add new article's sources.
// For example, a single article, or a feed.
type Service struct {
	feedsService    feedsService
	articlesService articlesService
	client          reader.GetClient
}

// New returns new sources instance.
func New(articlesService articlesService, feedsService feedsService) *Service {
	return &Service{
		articlesService: articlesService,
		feedsService:    feedsService,
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
		return nil, fmt.Errorf("failed to fetch url: %w", err)
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(ct, "text/html"):
		now := time.Now()
		article, err := s.articlesService.CreateArticle(ctx, resp.Body, url, &now)
		if err != nil {
			return nil, fmt.Errorf("failed to create article from source: %w", err)
		}
		request.Source.Name = article.Name
		return request.Source, nil
	case strings.HasPrefix(ct, "application/rss+xml"),
		strings.HasPrefix(ct, "application/atom+xml"),
		strings.HasPrefix(ct, "text/xm"):
		feed, err := s.feedsService.CreateFeed(ctx, resp.Body, url)
		if err != nil {
			return nil, fmt.Errorf("failed to create feed from source: %w", err)
		}
		request.Source.Name = feed.Name
		return request.Source, nil
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported format %s", ct)
	}
}
