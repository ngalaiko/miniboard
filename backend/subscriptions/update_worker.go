package subscriptions

import (
	"context"
	"fmt"
)

type worker struct {
	subscriptionsIDsToUpdate <-chan string
	db                       *database

	shutdown chan struct{}
	stopped  chan struct{}
}

func newWorker(subscriptionsIDsToUpdate <-chan string, db *database) *worker {
	return &worker{
		subscriptionsIDsToUpdate: subscriptionsIDsToUpdate,
		db:                       db,
		shutdown:                 make(chan struct{}),
		stopped:                  make(chan struct{}),
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
	fmt.Printf("test: updating %s\n", subscriptionID)
	return fmt.Errorf("not implemented")
}
