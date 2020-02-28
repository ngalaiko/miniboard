package rss

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/mmcdole/gofeed"
	"github.com/segmentio/ksuid"
	"golang.org/x/sync/errgroup"
	"miniboard.app/api/actor"
	"miniboard.app/articles"
	rss "miniboard.app/proto/users/rss/v1"
)

// Service is an RSS service.
type Service struct {
	parser          *gofeed.Parser
	articlesService *articles.Service
}

// New creates rss service.
func New(articlesService *articles.Service) *Service {
	parser := gofeed.NewParser()
	parser.Client = &http.Client{}
	return &Service{
		articlesService: articlesService,
		parser:          parser,
	}
}

// CreateFeed creates a new rss feed, fetches articles and schedules a next update.
func (s *Service) CreateFeed(ctx context.Context, reader io.Reader) (*rss.Feed, error) {
	actor, _ := actor.FromContext(ctx)

	feed, err := s.parser.Parse(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
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
		return nil, fmt.Errorf("failed to add feed: %w", err)
	}

	name := actor.Child("feeds", ksuid.New().String())
	return &rss.Feed{
		Name: name.String(),
	}, nil
}

func (s *Service) saveItem(ctx context.Context, item *gofeed.Item) error {
	resp, err := s.parser.Client.Get(item.Link)
	if err != nil {
		return fmt.Errorf("failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	link, _ := url.Parse(item.Link)
	if _, err := s.articlesService.CreateArticle(ctx, resp.Body, link); err != nil {
		return fmt.Errorf("failed to create article: %w", err)
	}
	return nil
}
