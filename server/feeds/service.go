package feeds

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
)

// Known errors.
var (
	errNotFound             = fmt.Errorf("not found")
	errFailedToDownloadFeed = fmt.Errorf("failed to download feed")
	errFailedToParseFeed    = fmt.Errorf("failed to parse feed")
)

type crawler interface {
	Crawl(context.Context, *url.URL) ([]byte, error)
}

// Service allows to manage feeds resource.
type Service struct {
	db      *database
	crawler crawler
	parser  *gofeed.Parser
}

// NewService returns new feeds service.
func NewService(db *sql.DB, crawler crawler) *Service {
	return &Service{
		db:      newDB(db),
		crawler: crawler,
		parser:  gofeed.NewParser(),
	}
}

// Create creates a feed from URL.
func (s *Service) Create(ctx context.Context, userID string, url *url.URL) (*Feed, error) {
	data, err := s.crawler.Crawl(ctx, url)
	if err != nil {
		return nil, errFailedToDownloadFeed
	}

	parsedFeed, err := s.parser.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, errFailedToParseFeed
	}

	feed := &Feed{
		ID:      uuid.New().String(),
		UserID:  userID,
		URL:     url.String(),
		Title:   parsedFeed.Title,
		Created: time.Now().Truncate(time.Nanosecond),
	}

	if err := s.db.Create(ctx, feed); err != nil {
		return nil, fmt.Errorf("failed to store feed: %w", err)
	}

	return feed, nil
}

// Get returns a feed by it's id.
func (s *Service) Get(ctx context.Context, id string, userID string) (*Feed, error) {
	feed, err := s.db.Get(ctx, id, userID)
	switch {
	case err == nil:
		return feed, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, errNotFound
	default:
		return nil, err
	}
}
