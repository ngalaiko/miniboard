package subscriptions

import (
	"bytes"
	"context"
	"fmt"
	"net/url"

	"github.com/mmcdole/gofeed"
	"github.com/ngalaiko/miniboard/backend/items"
)

type worker struct {
	subscriptionsIDsToUpdate <-chan string
	db                       *database
	parser                   *gofeed.Parser
	crawler                  crawler
	itemsService             itemsService
	logger                   logger

	shutdown chan struct{}
	stopped  chan struct{}
}

func newWorker(subscriptionsIDsToUpdate <-chan string, logger logger, db *database, parser *gofeed.Parser, crawler crawler, itemsService itemsService) *worker {
	return &worker{
		subscriptionsIDsToUpdate: subscriptionsIDsToUpdate,
		db:                       db,
		parser:                   parser,
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
				return err
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
	subscription, err := w.db.GetByID(ctx, subscriptionID)
	if err != nil {
		return err
	}

	sURL, err := url.Parse(subscription.URL)
	if err != nil {
		return err
	}

	data, err := w.crawler.Crawl(ctx, sURL)
	if err != nil {
		return errFailedToDownloadSubscription
	}

	parsedSubscription, err := w.parser.Parse(bytes.NewReader(data))
	if err != nil {
		return errFailedToParseSubscription
	}

	for _, item := range parsedSubscription.Items {
		if _, err := w.itemsService.Create(ctx, subscription.ID, item.Link, item.Title); err != nil && err != items.ErrAlreadyExists {
			return err
		}
	}

	return fmt.Errorf("not implemented")
}
