package rss

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/mmcdole/gofeed"
	"github.com/segmentio/ksuid"
	"github.com/spaolacci/murmur3"
	"golang.org/x/sync/errgroup"
	"miniboard.app/api/actor"
	"miniboard.app/articles"
	rss "miniboard.app/proto/users/rss/v1"
	"miniboard.app/storage"
)

// Service is an RSS service.
type Service struct {
	parser          *gofeed.Parser
	articlesService *articles.Service
	storage         storage.Storage
}

// New creates rss service.
func New(ctx context.Context, storage storage.Storage, articlesService *articles.Service) *Service {
	parser := gofeed.NewParser()
	parser.Client = &http.Client{}

	s := &Service{
		articlesService: articlesService,
		parser:          parser,
		storage:         storage,
	}
	go s.listenToUpdates(ctx)
	return s
}

// CreateFeed creates a new rss feed, fetches articles and schedules a next update.
func (s *Service) CreateFeed(ctx context.Context, reader io.Reader, url *url.URL) (*rss.Feed, error) {
	actor, _ := actor.FromContext(ctx)

	urlHash := murmur3.New128()
	_, _ = urlHash.Write([]byte(url.String()))

	// timestamp order == lexicographical order
	id, err := ksuid.FromParts(time.Now(), urlHash.Sum(nil))
	if err != nil {
		return nil, fmt.Errorf("failed to generate id: %w", err)
	}

	name := actor.Child("rss", fmt.Sprintf("%x", id.String()))

	if rawExisting, err := s.storage.Load(ctx, name); err == nil {
		feed := &rss.Feed{}
		if err := proto.Unmarshal(rawExisting, feed); err != nil {
			return nil, fmt.Errorf("failed to unmarshal the article: %w", err)
		}

		return feed, err
	}

	f := &rss.Feed{
		Name:        name.String(),
		LastFetched: ptypes.TimestampNow(),
		Url:         url.String(),
	}

	raw, err := proto.Marshal(f)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal feed: %w", err)
	}

	if err := s.storage.Store(ctx, name, raw); err != nil {
		return nil, fmt.Errorf("failed to save feed: %w", err)
	}

	if err := s.parse(ctx, reader); err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}

	return f, nil
}

func (s *Service) parse(ctx context.Context, reader io.Reader) error {
	feed, err := s.parser.Parse(reader)
	if err != nil {
		return fmt.Errorf("failed to parse feed: %w", err)
	}

	wg, ctx := errgroup.WithContext(ctx)
	for _, item := range feed.Items {
		item := item

		wg.Go(func() error {
			if err := s.saveItem(ctx, item); err != nil {
				return fmt.Errorf("failed to save item %s: %w", item.Link, err)
			}
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		return fmt.Errorf("failed to add feed: %w", err)
	}
	return nil
}

func (s *Service) saveItem(ctx context.Context, item *gofeed.Item) error {
	resp, err := s.parser.Client.Get(item.Link)
	if err != nil {
		return fmt.Errorf("failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	link, _ := url.Parse(item.Link)
	if _, err := s.articlesService.CreateArticle(ctx, resp.Body, link, item.PublishedParsed); err != nil {
		return fmt.Errorf("failed to create article: %w", err)
	}
	return nil
}
