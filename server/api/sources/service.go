package sources

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime/debug"
	"strings"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/api/actor"
	"miniboard.app/api/articles"
	"miniboard.app/api/feeds"
	"miniboard.app/fetch"
)

type articlesService interface {
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
	client          fetch.Fetcher
}

// NewService returns new sources instance.
func NewService(articlesService articlesService, feedsService feedsService, fetch fetch.Fetcher) *Service {
	return &Service{
		articlesService: articlesService,
		feedsService:    feedsService,
		client:          fetch,
	}
}

// CreateSource creates a new source.
func (s *Service) CreateSource(ctx context.Context, request *CreateSourceRequest) (*Source, error) {
	switch {
	case request.Source.Url != "":
		return s.createSourceFromURL(ctx, request.Source)
	case request.Source.Raw != nil:
		return s.createSourceFromRaw(ctx, request.Source)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "you need to provide either url or raw data")
	}
}

func (s *Service) createSourceFromRaw(ctx context.Context, source *Source) (*Source, error) {
	a, _ := actor.FromContext(ctx)

	if !strings.HasPrefix(http.DetectContentType(source.Raw), "text/xml") {
		return nil, status.Errorf(codes.InvalidArgument, "only raw xml files supported")
	}

	sources, err := parseOPML(source.Raw)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "only raw opml files supported")
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log().Errorf("%s: %s", r, debug.Stack())
			}
		}()

		start := time.Now()
		log().Infof("adding %d sources from opml", len(sources))

		ctx = actor.NewContext(context.Background(), a)

		for _, source := range sources {
			source := source

			_, err := s.createSourceFromURL(ctx, source)
			if err != nil {
				log().Errorf("failed to create source from '%s': %s", source.Url, err)
			}
		}
		log().Infof("added %d sources from opml in %s", len(sources), time.Since(start))
	}()

	source.Name = a.Child("opml", ksuid.New().String()).String()
	return source, nil
}

func (s *Service) createSourceFromURL(ctx context.Context, source *Source) (*Source, error) {
	url, err := url.ParseRequestURI(source.Url)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "url is invalid")
	}

	resp, err := s.client.Fetch(ctx, source.Url)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to fetch url: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, status.Errorf(codes.InvalidArgument, "failde to fetch source: response code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read source response: %s", err)
	}

	ct := resp.Header.Get("Content-Type")

	switch {
	case strings.Contains(ct, "text/html"):
		now := time.Now()
		article, err := s.articlesService.CreateArticle(ctx, bytes.NewBuffer(body), url, &now)
		if errors.Is(err, articles.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "article already exists")
		}
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create article from source: %s", err)
		}
		source.Name = article.Name
		return source, nil
	case strings.Contains(ct, "application/rss+xml"),
		strings.Contains(ct, "application/atom+xml"),
		strings.Contains(ct, "application/xml"),
		strings.Contains(ct, "text/xml"):
		feed, err := s.feedsService.CreateFeed(ctx, bytes.NewBuffer(body), url)
		if errors.Is(err, feeds.ErrAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "article already exists")
		}
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create feed from source: %s", err)
		}
		source.Name = feed.Name
		return source, nil
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported source content type '%s'", ct)
	}
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "sources",
	})
}
