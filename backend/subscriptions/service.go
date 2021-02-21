package subscriptions

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	"github.com/ngalaiko/miniboard/backend/items"
	"golang.org/x/sync/errgroup"
)

// Known errors.
var (
	errNotFound                     = fmt.Errorf("not found")
	errAlreadyExists                = fmt.Errorf("subscription already exists")
	errFailedToDownloadSubscription = fmt.Errorf("failed to download subscription")
	errFailedToParseSubscription    = fmt.Errorf("failed to parse subscription")
	errFailedToStoreItem            = fmt.Errorf("failed to store item")
)

type updateConfig struct {
	Workers  int           `yaml:"workers"`
	Interval time.Duration `yaml:"interval"`
}

// Config is a subscriptions service configuration.
type Config struct {
	Update *updateConfig `yaml:"update"`
}

type crawler interface {
	Crawl(context.Context, *url.URL) ([]byte, error)
}

type itemsService interface {
	Create(ctx context.Context, subscriptionID string, url string, title string) (*items.Item, error)
}

// Service allows to manage subscriptions resource.
type Service struct {
	db      *database
	crawler crawler
	parser  *gofeed.Parser
	cfg     *Config

	logger       logger
	itemsService itemsService

	workers         []*worker
	subscriptionIDs sync.Map
}

// NewService returns new subscriptions service.
func NewService(db *sql.DB, crawler crawler, logger logger, cfg *Config, itemsService itemsService) *Service {
	if cfg == nil {
		cfg = &Config{}
	}
	if cfg.Update == nil {
		cfg.Update = &updateConfig{}
	}
	if cfg.Update.Interval == 0 {
		cfg.Update.Interval = 5 * time.Minute
	}
	return &Service{
		db:           newDB(db, logger),
		crawler:      crawler,
		parser:       gofeed.NewParser(),
		cfg:          cfg,
		logger:       logger,
		itemsService: itemsService,
	}
}

// Create creates a subscription from URL.
func (s *Service) Create(ctx context.Context, userID string, url *url.URL, tagIDs []string) (*UserSubscription, error) {
	if exists, err := s.db.GetByURL(ctx, userID, url.String()); err == nil && exists != nil {
		return nil, errAlreadyExists
	}

	data, err := s.crawler.Crawl(ctx, url)
	if err != nil {
		s.logger.Error("failed to fetch subscription %s: %s", url, err)
		return nil, errFailedToDownloadSubscription
	}

	parsedSubscription, err := s.parser.Parse(bytes.NewReader(data))
	if err != nil {
		s.logger.Error("failed to parse subscription %s: %s", url, err)
		return nil, errFailedToParseSubscription
	}

	subscription := &UserSubscription{}
	subscription.ID = uuid.New().String()
	subscription.UserID = userID
	subscription.URL = url.String()
	subscription.Title = parsedSubscription.Title
	subscription.Created = time.Now().Truncate(time.Nanosecond)
	subscription.TagIDs = tagIDs

	if parsedSubscription.Image != nil {
		subscription.IconURL = &parsedSubscription.Image.URL
	}

	if err := s.db.Create(ctx, subscription); err != nil {
		return nil, fmt.Errorf("failed to store subscription: %w", err)
	}

	for _, item := range parsedSubscription.Items {
		if _, err := s.itemsService.Create(ctx, subscription.ID, item.Link, item.Title); err != nil && err != items.ErrAlreadyExists {
			s.logger.Error("failed to store feed %s item %s: %s", subscription.ID, item.Link, err)
			return nil, errFailedToStoreItem
		}
	}

	s.watchUpdates(&subscription.Subscription)

	return subscription, nil
}

// Get returns a subscription by it's id.
func (s *Service) Get(ctx context.Context, id string, userID string) (*UserSubscription, error) {
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
func (s *Service) List(ctx context.Context, userID string, pageSize int, createdLT *time.Time) ([]*UserSubscription, error) {
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
	subscriptionIDsToUpdate := make(chan string)

	subscriptions, err := s.db.ListAll(ctx)
	if err != nil {
		return err
	}
	for _, subscription := range subscriptions {
		s.watchUpdates(subscription)
	}

	updateCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go s.startUpdating(updateCtx, subscriptionIDsToUpdate)

	return s.startWorkers(ctx, subscriptionIDsToUpdate)
}

func (s *Service) watchUpdates(subscription *Subscription) {
	updated := subscription.Created
	if subscription.Updated != nil {
		updated = *subscription.Updated
	}
	s.subscriptionIDs.Store(subscription.ID, updated)
}

func (s *Service) startUpdating(ctx context.Context, subscriptionIDsToUpdate chan<- string) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			s.subscriptionIDs.Range(func(id interface{}, updated interface{}) bool {
				due := updated.(time.Time).Add(s.cfg.Update.Interval)
				if due.Before(time.Now()) {
					subscriptionIDsToUpdate <- id.(string)
					s.subscriptionIDs.Store(id, time.Now())
				}
				return true
			})
		}
	}
}

func (s *Service) startWorkers(ctx context.Context, subscriptionIDsToUpdate chan string) error {
	nWorkers := 10
	if s.cfg.Update != nil && s.cfg.Update.Workers > 0 {
		nWorkers = s.cfg.Update.Workers
	}

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < nWorkers; i++ {
		worker := newWorker(subscriptionIDsToUpdate, s.logger, s.db, s.parser, s.crawler, s.itemsService)
		s.workers = append(s.workers, worker)

		g.Go(restartingWorker(ctx, worker, s.logger))
	}

	return g.Wait()
}

func restartingWorker(ctx context.Context, worker *worker, logger logger) func() error {
	return func() error {
		if err := worker.Start(ctx); err != nil {
			logger.Error("subscriptions worker: %s", err)
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
