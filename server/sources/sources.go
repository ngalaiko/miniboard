package sources

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/ngalaiko/miniboard/server/actor"
	"github.com/ngalaiko/miniboard/server/articles"
	"github.com/ngalaiko/miniboard/server/feeds"
	"github.com/ngalaiko/miniboard/server/fetch"
	"github.com/ngalaiko/miniboard/server/operations"
	"github.com/sirupsen/logrus"
	longrunning "google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type articlesService interface {
	CreateArticle(context.Context, io.Reader, *url.URL, *time.Time, *string) (*articles.Article, error)
}

type feedsService interface {
	CreateFeed(context.Context, io.Reader, *url.URL) (*feeds.Feed, error)
}

type operationsService interface {
	CreateOperation(context.Context, *any.Any, operations.Operation) (*longrunning.Operation, error)
}

// Service allows to add new article's sources.
// For example, a single article, or a feed.
type Service struct {
	feedsService      feedsService
	articlesService   articlesService
	operationsService operationsService
	client            fetch.Fetcher
}

// NewService returns new sources instance.
func NewService(
	articlesService articlesService,
	feedsService feedsService,
	operationsService operationsService,
	fetch fetch.Fetcher) *Service {
	return &Service{
		articlesService:   articlesService,
		feedsService:      feedsService,
		operationsService: operationsService,
		client:            fetch,
	}
}

// CreateSource creates a new source.
func (s *Service) CreateSource(ctx context.Context, request *CreateSourceRequest) (*longrunning.Operation, error) {
	a, err := ptypes.MarshalAny(request)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse request")
	}
	return s.operationsService.CreateOperation(ctx, a, s.processSource)
}

func (s *Service) processSource(ctx context.Context, operation *longrunning.Operation) (*any.Any, error) {
	request := &CreateSourceRequest{}
	if err := ptypes.UnmarshalAny(operation.Metadata, request); err != nil {
		return nil, fmt.Errorf("failed to unmarshal operation metadata: %w", err)
	}
	return s.createSource(ctx, request)
}

func (s *Service) createSource(ctx context.Context, request *CreateSourceRequest) (*any.Any, error) {
	switch {
	case request.Source.Url != "":
		return s.createSourceFromURL(ctx, request.Source)
	case request.Source.Raw != nil:
		return s.createSourceFromRaw(ctx, request.Source)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "you need to provide either url or raw data")
	}
}

func (s *Service) createSourceFromRaw(ctx context.Context, source *Source) (*any.Any, error) {
	a, _ := actor.FromContext(ctx)

	if !strings.HasPrefix(http.DetectContentType(source.Raw), "text/xml") {
		return nil, status.Errorf(codes.InvalidArgument, "only raw xml files supported")
	}

	sources, err := parseOPML(source.Raw)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "only raw opml files supported")
	}

	start := time.Now()
	log().Infof("adding %d sources from opml", len(sources))

	ctx = actor.NewContext(context.Background(), a.ID)

	list := &SourceList{}
	for _, source := range sources {
		any, err := s.createSourceFromURL(ctx, source)
		if err != nil {
			return nil, fmt.Errorf("failed to create source from '%s': %s", source.Url, err)
		}

		s := &Source{}
		if err := ptypes.UnmarshalAny(any, s); err != nil {
			return nil, fmt.Errorf("failed to unmarshal source to list: %w", err)
		}
		list.Sources = append(list.Sources, s)
	}
	log().Infof("added %d sources from opml in %s", len(sources), time.Since(start))

	return ptypes.MarshalAny(list)
}

func (s *Service) createSourceFromURL(ctx context.Context, source *Source) (*any.Any, error) {
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
		article, err := s.articlesService.CreateArticle(ctx, bytes.NewBuffer(body), url, &now, nil)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to create article from source: %s", err)
		}
		return ptypes.MarshalAny(article)
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
		return ptypes.MarshalAny(feed)
	default:
		return nil, status.Errorf(codes.InvalidArgument, "unsupported source content type '%s'", ct)
	}
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "sources",
	})
}
