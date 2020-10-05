package feeds

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/ngalaiko/miniboard/server/actor"
	"github.com/ngalaiko/miniboard/server/feeds/db"
	"github.com/ngalaiko/miniboard/server/feeds/parsers"
	"github.com/ngalaiko/miniboard/server/fetch"
	articles "github.com/ngalaiko/miniboard/server/genproto/articles/v1"
	feeds "github.com/ngalaiko/miniboard/server/genproto/feeds/v1"
	"github.com/spaolacci/murmur3"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type articlesService interface {
	CreateArticle(context.Context, io.Reader, *url.URL, *time.Time, *string) (*articles.Article, error)
}

type parser interface {
	Parse(feed io.Reader) (*parsers.Feed, error)
}

// Service is a Feeds service.
type Service struct {
	storage *db.DB

	parser          parser
	articlesService articlesService
	fetcher         fetch.Fetcher
}

// NewService creates feeds service.
func NewService(ctx context.Context, sqldb *sql.DB, fetcher fetch.Fetcher, articlesService articlesService, parser parser) *Service {
	s := &Service{
		articlesService: articlesService,
		parser:          parser,
		storage:         db.New(sqldb),
		fetcher:         fetcher,
	}
	go s.listenToUpdates(ctx)
	return s
}

// GetFeed returns a feed.
func (s *Service) GetFeed(ctx context.Context, request *feeds.GetFeedRequest) (*feeds.Feed, error) {
	a, _ := actor.FromContext(ctx)
	feed, err := s.storage.Get(ctx, request.Id, a.ID)
	switch {
	case err == nil:
	case errors.Is(err, sql.ErrNoRows):
		return nil, status.Errorf(codes.NotFound, "not found")
	default:
		return nil, status.Errorf(codes.Internal, "failed to load the feed")
	}

	return feed, nil
}

// ListFeeds returns a list of feeds.
func (s *Service) ListFeeds(ctx context.Context, request *feeds.ListFeedsRequest) (*feeds.ListFeedsResponse, error) {
	a, _ := actor.FromContext(ctx)

	request.PageSize++
	ff, err := s.storage.List(ctx, a.ID, request)

	var nextPageToken string
	if len(ff) == int(request.PageSize) {
		nextPageToken = base64.StdEncoding.EncodeToString([]byte(ff[len(ff)-1].Id))
		ff = ff[:request.PageSize-1]
	}

	switch err {
	case nil, sql.ErrNoRows:
		return &feeds.ListFeedsResponse{
			Feeds:         ff,
			NextPageToken: nextPageToken,
		}, nil
	default:
		log().Error(err)
		return nil, status.Errorf(codes.Internal, "failed to list feeds")
	}
}

// CreateFeed creates a new feeds feed, fetches articles and schedules a next update.
func (s *Service) CreateFeed(ctx context.Context, reader io.Reader, url *url.URL) (*feeds.Feed, error) {
	a, _ := actor.FromContext(ctx)

	urlHash := murmur3.New128()
	if _, err := urlHash.Write([]byte(url.String())); err != nil {
		return nil, fmt.Errorf("failed to hash url: %w", err)
	}

	id := fmt.Sprintf("%x", urlHash.Sum(nil))

	if feed, err := s.storage.Get(ctx, id, a.ID); err == nil {
		return feed, nil
	}

	feed := &feeds.Feed{
		Id:     id,
		UserId: a.ID,
		Url:    url.String(),
	}

	items, err := s.parse(reader, feed)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %s", err)
	}

	if err := s.storage.Create(ctx, feed); err != nil {
		return nil, fmt.Errorf("failed to save feed: %w", err)
	}

	for _, item := range items {
		if err := s.saveItem(ctx, item, feed); err != nil {
			log().Errorf("failed to save item %s: %s", item.Link, err)
			continue
		}
	}

	return feed, nil
}

func (s *Service) parse(reader io.Reader, feed *feeds.Feed) ([]*parsers.Item, error) {
	var updateLeeway = 24 * time.Hour

	parsedFeed, err := s.parser.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}

	feed.Title = parsedFeed.Title
	if parsedFeed.Image != nil {
		feed.IconUrl = &wrappers.StringValue{
			Value: parsedFeed.Image.URL,
		}
	}

	lastFetched := time.Time{}
	if feed.LastFetched != nil {
		lastFetched, err = ptypes.Timestamp(feed.LastFetched)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", err)
		}
	}

	items := make([]*parsers.Item, 0, len(parsedFeed.Items))
	for _, item := range parsedFeed.Items {
		updatedTime := latestTimestamp(item.UpdatedParsed, item.PublishedParsed)
		if updatedTime.Before(lastFetched.Add(-1 * updateLeeway)) {
			log().Infof("skipping item %s from %s: updated at %s", item.Link, feed.Id, updatedTime)
			continue
		}

		items = append(items, item)
	}

	feed.LastFetched = ptypes.TimestampNow()

	return items, nil
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

func (s *Service) saveItem(ctx context.Context, item *parsers.Item, feed *feeds.Feed) error {
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
	if _, err := s.articlesService.CreateArticle(ctx, bytes.NewReader(body), link, published, &feed.Id); err != nil {
		return fmt.Errorf("failed to create article: %w", err)
	}
	return nil
}
