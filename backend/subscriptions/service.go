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
	"golang.org/x/sync/errgroup"
)

// Known errors.
var (
	errNotFound                     = fmt.Errorf("not found")
	errAlreadyExists                = fmt.Errorf("subscription already exists")
	errFailedToDownloadSubscription = fmt.Errorf("failed to download subscription")
	errFailedToParseSubscription    = fmt.Errorf("failed to parse subscription")
)

// Config is a subscriptions service configuration.
type Config struct {
	Update *struct {
		Workers int `yaml:"workers"`
	} `yaml:"update"`
}

type crawler interface {
	Crawl(context.Context, *url.URL) ([]byte, error)
}

// Service allows to manage subscriptions resource.
type Service struct {
	db      *database
	crawler crawler
	parser  *gofeed.Parser
	cfg     *Config

	logger logger

	workers []*worker
}

// NewService returns new subscriptions service.
func NewService(db *sql.DB, crawler crawler, logger logger, cfg *Config) *Service {
	if cfg == nil {
		cfg = &Config{}
	}
	return &Service{
		db:      newDB(db, logger),
		crawler: crawler,
		parser:  gofeed.NewParser(),
		cfg:     cfg,
		logger:  logger,
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

// Start starts background subscriptions updates.
func (s *Service) Start(ctx context.Context) error {
	nWorkers := 10
	if s.cfg.Update != nil && s.cfg.Update.Workers > 0 {
		nWorkers = s.cfg.Update.Workers
	}

	subscriptionIDsToUpdate := make(chan string)

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < nWorkers; i++ {
		worker := newWorker(subscriptionIDsToUpdate, s.db)
		s.workers = append(s.workers, worker)

		g.Go(restartingWorker(ctx, worker, s.logger))
	}

	return g.Wait()
}

func restartingWorker(ctx context.Context, worker *worker, logger logger) func() error {
	return func() error {
		if err := worker.Start(ctx); err != nil {
			logger.Error("subscriptions worker: %w", err)
			return restartingWorker(ctx, worker, logger)()
		}
		return nil
	}
}

// Shutdown stops background subscriptions updates.
func (s *Service) Shutdown(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)
	for _, worker := range s.workers {
		worker := worker
		g.Go(func() error {
			if err := worker.Shutdown(ctx); err != nil {
				return fmt.Errorf("failed to shutdown a worker")
			}
			return nil
		})
	}
	return g.Wait()
}
