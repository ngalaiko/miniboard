package subscriptions

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
	errNotFound                     = fmt.Errorf("not found")
	errAlreadyExists                = fmt.Errorf("subscription already exists")
	errFailedToDownloadSubscription = fmt.Errorf("failed to download subscription")
	errFailedToParseSubscription    = fmt.Errorf("failed to parse subscription")
)

type crawler interface {
	Crawl(context.Context, *url.URL) ([]byte, error)
}

// Service allows to manage subscriptions resource.
type Service struct {
	db      *database
	crawler crawler
	parser  *gofeed.Parser
}

// NewService returns new subscriptions service.
func NewService(db *sql.DB, crawler crawler, logger logger) *Service {
	return &Service{
		db:      newDB(db, logger),
		crawler: crawler,
		parser:  gofeed.NewParser(),
	}
}

// Create creates a subscription from URL.
func (s *Service) Create(ctx context.Context, userID string, url *url.URL, tagIDs []string) (*Subscription, error) {
	if exists, err := s.db.GetByURL(ctx, userID, url.String()); err == nil && exists != nil {
		return nil, errAlreadyExists
	}

	data, err := s.crawler.Crawl(ctx, url)
	if err != nil {
		return nil, errFailedToDownloadSubscription
	}

	parsedSubscription, err := s.parser.Parse(bytes.NewReader(data))
	if err != nil {
		return nil, errFailedToParseSubscription
	}

	subscription := &Subscription{
		ID:      uuid.New().String(),
		UserID:  userID,
		URL:     url.String(),
		Title:   parsedSubscription.Title,
		Created: time.Now().Truncate(time.Nanosecond),
		TagIDs:  tagIDs,
	}

	if parsedSubscription.Image != nil {
		subscription.IconURL = &parsedSubscription.Image.URL
	}

	if err := s.db.Create(ctx, subscription); err != nil {
		return nil, fmt.Errorf("failed to store subscription: %w", err)
	}

	return subscription, nil
}

// Get returns a subscription by it's id.
func (s *Service) Get(ctx context.Context, id string, userID string) (*Subscription, error) {
	subscription, err := s.db.Get(ctx, id, userID)
	switch {
	case err == nil:
		return subscription, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, errNotFound
	default:
		return nil, err
	}
}

// List returns a list of user subscriptions.
func (s *Service) List(ctx context.Context, userID string, pageSize int, createdLT *time.Time) ([]*Subscription, error) {
	subscriptions, err := s.db.List(ctx, userID, pageSize, createdLT)
	switch {
	case err == nil:
		return subscriptions, nil
	default:
		return nil, err
	}
}
