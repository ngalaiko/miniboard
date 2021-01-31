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
	errAlreadyExists        = fmt.Errorf("feed already exists")
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
func NewService(db *sql.DB, crawler crawler, logger logger) *Service {
	return &Service{
		db:      newDB(db, logger),
		crawler: crawler,
		parser:  gofeed.NewParser(),
	}
}

// Create creates a feed from URL.
func (s *Service) Create(ctx context.Context, userID string, url *url.URL, tagIDs []string) (*Feed, error) {
	if exists, err := s.db.GetByURL(ctx, userID, url.String()); err == nil && exists != nil {
		return nil, errAlreadyExists
	}

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
		TagIDs:  tagIDs,
	}

	if parsedFeed.Image != nil {
		feed.IconURL = &parsedFeed.Image.URL
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

// List returns a list of user feeds.
func (s *Service) List(ctx context.Context, userID string, pageSize int, createdLT *time.Time) ([]*Feed, error) {
	feeds, err := s.db.List(ctx, userID, pageSize, createdLT)
	switch {
	case err == nil:
		return feeds, nil
	default:
		return nil, err
	}
}
