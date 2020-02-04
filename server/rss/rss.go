package rss

import (
	"context"
	"io"
	"net/http"
	"net/url"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "failed to parse feed")
	}

	wg, ctx := errgroup.WithContext(ctx)
	for _, item := range feed.Items {
		item := item
		wg.Go(func() error {
			if err := s.saveItem(ctx, item); err != nil {
				return errors.Wrapf(err, "failed to save item %s", item.Link)
			}
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		return nil, errors.Wrap(err, "failed to add feed")
	}

	name := actor.Child("feeds", ksuid.New().String())
	return &rss.Feed{
		Name: name.String(),
	}, nil
}

func (s *Service) saveItem(ctx context.Context, item *gofeed.Item) error {
	resp, err := s.parser.Client.Get(item.Link)
	if err != nil {
		return errors.Wrap(err, "failed to fetch")
	}
	defer resp.Body.Close()

	link, _ := url.Parse(item.Link)
	if _, err := s.articlesService.CreateArticle(ctx, resp.Body, link); err != nil {
		return errors.Wrap(err, "failed to create article")
	}
	return nil
}
