package feeds

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"runtime/debug"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/mmcdole/gofeed"
	"github.com/ngalaiko/miniboard/server/api/actor"
	"github.com/ngalaiko/miniboard/server/api/articles"
	"github.com/ngalaiko/miniboard/server/fetch"
	"github.com/ngalaiko/miniboard/server/storage"
	"github.com/ngalaiko/miniboard/server/storage/resource"
	"github.com/spaolacci/murmur3"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// Known errors.
var (
	ErrAlreadyExists = errors.New("article already exists")
)

type articlesService interface {
	CreateArticle(context.Context, io.Reader, *url.URL, *time.Time) (*articles.Article, error)
}

// Service is a Feeds service.
type Service struct {
	parser          *gofeed.Parser
	articlesService articlesService
	storage         storage.Storage
	fetcher         fetch.Fetcher
}

// NewService creates feeds service.
func NewService(ctx context.Context, storage storage.Storage, fetcher fetch.Fetcher, articlesService articlesService) *Service {
	parser := gofeed.NewParser()

	s := &Service{
		articlesService: articlesService,
		parser:          parser,
		storage:         storage,
		fetcher:         fetcher,
	}
	go s.listenToUpdates(ctx)
	return s
}

// GetFeed returns a feed.
func (s *Service) GetFeed(ctx context.Context, request *GetFeedRequest) (*Feed, error) {
	name := resource.ParseName(request.Name)

	if !actor.Owns(ctx, name) {
		return nil, status.Errorf(codes.PermissionDenied, "forbidden")
	}

	raw, err := s.storage.Load(ctx, name)
	switch {
	case err == nil:
	case errors.Is(err, storage.ErrNotFound):
		return nil, status.Errorf(codes.NotFound, "not found")
	default:
		return nil, status.Errorf(codes.Internal, "failed to load the feed")
	}

	feed := &Feed{}
	if err := proto.Unmarshal(raw, feed); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unmarshal the feed")
	}

	return feed, nil
}

// ListFeeds returns a list of feeds.
func (s *Service) ListFeeds(ctx context.Context, request *ListFeedsRequest) (*ListFeedsResponse, error) {
	actor, _ := actor.FromContext(ctx)
	lookFor := actor.Child("feeds", "*")

	var from *resource.Name
	if request.PageToken != "" {
		decoded, err := base64.StdEncoding.DecodeString(request.PageToken)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid page token")
		}
		from = resource.ParseName(string(decoded))
	}

	ff := []*Feed{}
	err := s.storage.ForEach(ctx, lookFor, from, func(r *resource.Resource) (bool, error) {
		f := &Feed{}
		if err := proto.Unmarshal(r.Data, f); err != nil {
			return false, status.Errorf(codes.Internal, "failed to unmarshal feed")
		}

		ff = append(ff, f)

		if len(ff) == int(request.PageSize+1) {
			return false, nil
		}

		return true, nil
	})

	var nextPageToken string
	if len(ff) == int(request.PageSize+1) {
		nextPageToken = base64.StdEncoding.EncodeToString([]byte(ff[len(ff)-1].Name))
		ff = ff[:request.PageSize]
	}

	switch err {
	case nil, storage.ErrNotFound:
		return &ListFeedsResponse{
			Feeds:         ff,
			NextPageToken: nextPageToken,
		}, nil
	default:
		log().Error(err)
		return nil, status.Errorf(codes.Internal, "failed to list feeds")
	}
}

// CreateFeed creates a new feeds feed, fetches articles and schedules a next update.
func (s *Service) CreateFeed(ctx context.Context, reader io.Reader, url *url.URL) (*Feed, error) {
	a, _ := actor.FromContext(ctx)

	urlHash := murmur3.New128()
	_, _ = urlHash.Write([]byte(url.String()))

	name := a.Child("feeds", fmt.Sprintf("%x", urlHash.Sum(nil)))
	if _, err := s.storage.Load(ctx, name); err == nil {
		return nil, ErrAlreadyExists
	}

	f := &Feed{
		Name: name.String(),
		Url:  url.String(),
	}

	raw, err := proto.Marshal(f)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal feed: %w", err)
	}

	if err := s.storage.Store(ctx, name, raw); err != nil {
		return nil, fmt.Errorf("failed to save feed: %w", err)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log().Errorf("%s: %s", r, debug.Stack())
			}
		}()
		if err := s.parse(actor.NewContext(context.Background(), a), reader, f); err != nil {
			log().Errorf("failed to parse feed: %s", err)
		}
	}()

	return f, nil
}

func (s *Service) parse(ctx context.Context, reader io.Reader, f *Feed) error {
	var updateLeeway = 24 * time.Hour

	feed, err := s.parser.Parse(reader)
	if err != nil {
		return fmt.Errorf("failed to parse feed: %w", err)
	}

	lastFetched := time.Time{}
	if f.LastFetched != nil {
		lastFetched, err = ptypes.Timestamp(f.LastFetched)
		if err != nil {
			return fmt.Errorf("failed to parse timestamp: %w", err)
		}
	}

	updated := false
	for _, item := range feed.Items {
		item := item

		updatedTime := latestTimestamp(item.UpdatedParsed, item.PublishedParsed)
		if updatedTime.Before(lastFetched.Add(-1 * updateLeeway)) {
			log().Infof("skipping item %s from %s: updated at %s", item.Link, f.Name, updatedTime)
			continue
		}

		updated = true

		if err := s.saveItem(ctx, item); err != nil && !errors.Is(err, articles.ErrAlreadyExists) {
			log().Errorf("failed to save item %s: %s", item.Link, err)
			continue
		}
	}

	if !updated {
		return nil
	}

	f.LastFetched = ptypes.TimestampNow()
	f.Title = feed.Title

	raw, err := proto.Marshal(f)
	if err != nil {
		return fmt.Errorf("failed to marshal feed: %w", err)
	}

	if err := s.storage.Store(ctx, resource.ParseName(f.Name), raw); err != nil {
		return fmt.Errorf("failed to save feed: %w", err)
	}

	log().Infof("feed %s updated", f.Name)

	return nil
}

func latestTimestamp(ts ...*time.Time) time.Time {
	latest := time.Time{}
	for _, t := range ts {
		if t == nil {
			continue
		}

		if latest.Before(*t) {
			latest = *t
		}
	}
	return latest
}

func (s *Service) saveItem(ctx context.Context, item *gofeed.Item) error {
	resp, err := s.fetcher.Fetch(ctx, item.Link)
	if err != nil {
		return fmt.Errorf("failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	published := item.UpdatedParsed
	if item.PublishedParsed != nil {
		published = item.PublishedParsed
	}

	link, _ := url.Parse(item.Link)
	if _, err := s.articlesService.CreateArticle(ctx, bytes.NewReader(body), link, published); err != nil {
		return fmt.Errorf("failed to create article: %w", err)
	}
	return nil
}
