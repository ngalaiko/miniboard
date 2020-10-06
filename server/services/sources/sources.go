package sources

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/ngalaiko/miniboard/server/fetch"
	articles "github.com/ngalaiko/miniboard/server/genproto/articles/v1"
	feeds "github.com/ngalaiko/miniboard/server/genproto/feeds/v1"
	longrunning "github.com/ngalaiko/miniboard/server/genproto/google/longrunning"
	sources "github.com/ngalaiko/miniboard/server/genproto/sources/v1"
	"github.com/ngalaiko/miniboard/server/services/operations"
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
func (s *Service) CreateSource(ctx context.Context, request *sources.CreateSourceRequest) (*longrunning.Operation, error) {
	a, err := ptypes.MarshalAny(request)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse request")
	}
	return s.operationsService.CreateOperation(ctx, a, s.processSource)
}

func (s *Service) processSource(ctx context.Context, operation *longrunning.Operation, status chan<- *longrunning.Operation) error {
	request := &sources.CreateSourceRequest{}
	if err := ptypes.UnmarshalAny(operation.Metadata, request); err != nil {
		return fmt.Errorf("failed to unmarshal operation metadata: %w", err)
	}

	return s.createSource(ctx, request, operation, status)
}

func (s *Service) createSource(ctx context.Context, request *sources.CreateSourceRequest, operation *longrunning.Operation, status chan<- *longrunning.Operation) error {
	switch {
	case request.Source.Url != "":
		result, err := s.createSourceFromURL(ctx, request.Source)
		if err != nil {
			return err
		}

		operation.Result = &longrunning.Operation_Response{
			Response: result,
		}
		operation.Done = true

		status <- operation
		return nil
	default:
		return fmt.Errorf("you need to provide url ")
	}
}

func (s *Service) createSourceFromURL(ctx context.Context, source *sources.Source) (*any.Any, error) {
	url, err := url.ParseRequestURI(source.Url)
	if err != nil {
		return nil, fmt.Errorf("url is invalid")
	}

	resp, err := s.client.Fetch(ctx, source.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch url: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failde to fetch source: response code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read source response: %s", err)
	}

	ct := resp.Header.Get("Content-Type")

	switch {
	case strings.Contains(ct, "text/html"):
		now := time.Now()
		article, err := s.articlesService.CreateArticle(ctx, bytes.NewBuffer(body), url, &now, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create article from source: %s", err)
		}
		return ptypes.MarshalAny(article)
	case strings.Contains(ct, "application/rss+xml"),
		strings.Contains(ct, "application/atom+xml"),
		strings.Contains(ct, "application/xml"),
		strings.Contains(ct, "text/xml"):
		feed, err := s.feedsService.CreateFeed(ctx, bytes.NewBuffer(body), url)
		if err != nil {
			return nil, fmt.Errorf("failed to create feed from source: %s", err)
		}
		return ptypes.MarshalAny(feed)
	default:
		return nil, fmt.Errorf("unsupported source content type '%s'", ct)
	}
}
