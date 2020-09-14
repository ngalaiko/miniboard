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
	"runtime/debug"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/mmcdole/gofeed"
	"github.com/ngalaiko/miniboard/server/actor"
	"github.com/ngalaiko/miniboard/server/articles"
	"github.com/ngalaiko/miniboard/server/fetch"
	"github.com/spaolacci/murmur3"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// FromID returns page token converted to id.
func (request *ListFeedsRequest) FromID() (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(request.PageToken)
	if err != nil {
		return "", status.Errorf(codes.InvalidArgument, "invalid page token")
	}
	return string(decoded), nil
}

// Known errors.
var (
	ErrAlreadyExists = errors.New("article already exists")
)

type articlesService interface {
	CreateArticle(context.Context, io.Reader, *url.URL, *time.Time, string) (*articles.Article, error)
}

// Service is a Feeds service.
type Service struct {
	storage *feedsDB

	parser          *gofeed.Parser
	articlesService articlesService
	fetcher         fetch.Fetcher
}

// NewService creates feeds service.
func NewService(ctx context.Context, sqldb *sql.DB, fetcher fetch.Fetcher, articlesService articlesService) *Service {
	parser := gofeed.NewParser()

	s := &Service{
		articlesService: articlesService,
		parser:          parser,
		storage:         newDB(sqldb),
		fetcher:         fetcher,
	}
	go s.listenToUpdates(ctx)
	return s
}

// GetFeed returns a feed.
func (s *Service) GetFeed(ctx context.Context, request *GetFeedRequest) (*Feed, error) {
	feed, err := s.storage.Get(ctx, request.Id)
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
func (s *Service) ListFeeds(ctx context.Context, request *ListFeedsRequest) (*ListFeedsResponse, error) {
	request.PageSize++
	ff, err := s.storage.List(ctx, request)

	var nextPageToken string
	if len(ff) == int(request.PageSize+1) {
		nextPageToken = base64.StdEncoding.EncodeToString([]byte(ff[len(ff)-1].Id))
		ff = ff[:request.PageSize-1]
	}

	switch err {
	case nil, sql.ErrNoRows:
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
	if _, err := urlHash.Write([]byte(url.String())); err != nil {
		return nil, fmt.Errorf("failed to hash url: %w", err)
	}

	id := fmt.Sprintf("%x", urlHash.Sum(nil))

	if _, err := s.storage.Get(ctx, id); err == nil {
		return nil, ErrAlreadyExists
	}

	feed := &Feed{
		Id:     id,
		UserId: a.ID(),
		Url:    url.String(),
	}

	if err := s.storage.Create(ctx, feed); err != nil {
		return nil, fmt.Errorf("failed to save feed: %w", err)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log().Errorf("%s: %s", r, debug.Stack())
			}
		}()
		if err := s.parse(actor.NewContext(context.Background(), a), reader, feed); err != nil {
			log().Errorf("failed to parse feed: %s", err)
		}
	}()

	return feed, nil
}

func (s *Service) parse(ctx context.Context, reader io.Reader, feed *Feed) error {
	var updateLeeway = 24 * time.Hour

	parsedFeed, err := s.parser.Parse(reader)
	if err != nil {
		return fmt.Errorf("failed to parse feed: %w", err)
	}

	lastFetched := time.Time{}
	if feed.LastFetched != nil {
		lastFetched, err = ptypes.Timestamp(feed.LastFetched)
		if err != nil {
			return fmt.Errorf("failed to parse timestamp: %w", err)
		}
	}

	updated := false
	for _, item := range parsedFeed.Items {
		item := item

		updatedTime := latestTimestamp(item.UpdatedParsed, item.PublishedParsed)
		if updatedTime.Before(lastFetched.Add(-1 * updateLeeway)) {
			log().Infof("skipping item %s from %s: updated at %s", item.Link, feed.Id, updatedTime)
			continue
		}

		updated = true

		if err := s.saveItem(ctx, item, feed); err != nil {
			log().Errorf("failed to save item %s: %s", item.Link, err)
			continue
		}
	}

	if !updated {
		return nil
	}

	feed.LastFetched = ptypes.TimestampNow()

	if err := s.storage.Update(ctx, feed); err != nil {
		return fmt.Errorf("failed to save feed: %w", err)
	}

	log().Infof("feed %s updated", feed.Id)

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

func (s *Service) saveItem(ctx context.Context, item *gofeed.Item, feed *Feed) error {
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
	if _, err := s.articlesService.CreateArticle(ctx, bytes.NewReader(body), link, published, feed.Id); err != nil {
		return fmt.Errorf("failed to create article: %w", err)
	}
	return nil
}
