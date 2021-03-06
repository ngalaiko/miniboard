package subscriptions

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/ngalaiko/miniboard/backend/items"
	"github.com/ngalaiko/miniboard/backend/subscriptions/parser"
)

type worker struct {
	subscriptionsIDsToUpdate <-chan string
	db                       *database
	crawler                  crawler
	itemsService             itemsService
	logger                   logger

	shutdown chan struct{}
	stopped  chan struct{}
}

func newWorker(subscriptionsIDsToUpdate <-chan string, logger logger, db *database, crawler crawler, itemsService itemsService) *worker {
	return &worker{
		subscriptionsIDsToUpdate: subscriptionsIDsToUpdate,
		db:                       db,
		crawler:                  crawler,
		itemsService:             itemsService,
		shutdown:                 make(chan struct{}),
		stopped:                  make(chan struct{}),
		logger:                   logger,
	}
}

// Start starts the worker.
func (w *worker) Start(ctx context.Context) error {
	for {
		select {
		case subscriptionID := <-w.subscriptionsIDsToUpdate:
			if err := w.update(ctx, subscriptionID); err != nil {
				w.logger.Error("subscriptions worker: %s", err)
			}
		case <-w.shutdown:
			close(w.stopped)
			return nil
		}
	}
}

// Shutdown stops the worker.
func (w *worker) Shutdown(ctx context.Context) error {
	close(w.shutdown)
	<-w.stopped
	return nil
}

func (w *worker) update(ctx context.Context, subscriptionID string) error {
	w.logger.Debug("updating: subscription %s", subscriptionID)

	subscription, err := w.db.GetByID(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("get subscription %s: %w", subscriptionID, err)
	}

	sURL, err := url.Parse(subscription.URL)
	if err != nil {
		return fmt.Errorf("%s: %w", sURL, err)
	}

	data, err := w.crawler.Crawl(ctx, sURL)
	if err != nil {
		return fmt.Errorf("%s: %w", sURL, err)
	}

	parsedSubscription, err := parser.Parse(data, w.logger)
	if err != nil {
		return fmt.Errorf("%s: %w", sURL, err)
	}

	now := time.Now()
	for _, item := range parsedSubscription.Items {
		if item.Date == nil {
			item.Date = &now
		}
		if _, err := w.itemsService.Create(ctx, subscription.ID, item.Link, item.Title, item.Date, item.Content); err != nil && !errors.Is(err, items.ErrAlreadyExists) {
			return fmt.Errorf("subscription %s, item %s: %w", subscription.ID, item.Link, err)
		}
	}

	if err := w.db.Update(ctx, subscription); err != nil {
		return fmt.Errorf("failed to update subscription %s: %w", subscription.ID, err)
	}

	w.logger.Debug("updated: subscription %s", subscriptionID)

	return nil
}
