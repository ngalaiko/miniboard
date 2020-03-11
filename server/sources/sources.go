package sources

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/api/actor"
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
	switch {
	case request.Source.Url != "":
		return s.createSourceFromURL(ctx, request.Source)
	case request.Source.Raw != nil:
		return s.createSourceFromRaw(ctx, request.Source)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "you need to provide either url or raw data")
	}
}

func (s *Service) createSourceFromRaw(ctx context.Context, source *sources.Source) (*sources.Source, error) {
	if !strings.HasPrefix(http.DetectContentType(source.Raw), "text/xml") {
		return nil, status.Errorf(codes.InvalidArgument, "only raw xml files supported")
	}

	sources, err := parseOPML(source.Raw)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "only raw opml files supported")
	}

	for _, source := range sources {
		_, err := s.createSourceFromURL(ctx, source)
		if err != nil {
			return nil, err
		}
	}

	actor, _ := actor.FromContext(ctx)

	source.Name = actor.Child("opml", ksuid.New().String()).String()
	return source, nil
}

func (s *Service) createSourceFromURL(ctx context.Context, source *sources.Source) (*sources.Source, error) {
	url, err := url.ParseRequestURI(source.Url)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "url is invalid")
	}

	resp, err := s.client.Get(source.Url)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to fetch url: %s", err)
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(ct, "text/html"):
		now := time.Now()
		article, err := s.articlesService.CreateArticle(ctx, resp.Body, url, &now)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create article from source: %s", err)
		}
		source.Name = article.Name
		return source, nil
	case strings.HasPrefix(ct, "application/rss+xml"),
		strings.HasPrefix(ct, "application/atom+xml"),
		strings.HasPrefix(ct, "application/xml"),
		strings.HasPrefix(ct, "text/xml"):
		feed, err := s.feedsService.CreateFeed(ctx, resp.Body, url)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create feed from source: %s", err)
		}
		source.Name = feed.Name
		return source, nil
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported source content type '%s'", ct)
	}
}
