package sources

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/api/actor"
	articles "miniboard.app/proto/users/articles/v1"
	feeds "miniboard.app/proto/users/feeds/v1"
	sources "miniboard.app/proto/users/sources/v1"
)

type articlesService interface {
	articles.ArticlesServiceServer
	CreateArticle(context.Context, io.Reader, *url.URL, *time.Time) (*articles.Article, error)
}

type feedsService interface {
	CreateFeed(context.Context, io.Reader, *url.URL) (*feeds.Feed, error)
}

type getClient interface {
	Get(string) (*http.Response, error)
}

// Service allows to add new article's sources.
// For example, a single article, or a feed.
type Service struct {
	feedsService    feedsService
	articlesService articlesService
	client          getClient
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

	wg := &sync.WaitGroup{}
	for _, source := range sources {
		source := source

		wg.Add(1)
		go func() {
			_, err := s.createSourceFromURL(ctx, source)
			if err != nil {
				log().Errorf("failed to create source from '%s': %s", source.Url, err)
			}

			wg.Done()
		}()
	}

	wg.Wait()

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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to read source response: %s", err)
	}

	ct := resp.Header.Get("Content-Type")

	switch {
	case strings.Contains(ct, "text/html"):
		now := time.Now()
		article, err := s.articlesService.CreateArticle(ctx, bytes.NewBuffer(body), url, &now)
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
