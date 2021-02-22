package operations

import (
	"context"
	"fmt"
)

type worker struct {
	logger       logger
	processQueue <-chan *Operation
	db           *database

	shutdown chan struct{}
	stopped  chan struct{}
}

func newWorker(processQueue <-chan *Operation, logger logger, db *database) *worker {
	return &worker{
		logger:       logger,
		processQueue: processQueue,
		db:           db,
		shutdown:     make(chan struct{}),
		stopped:      make(chan struct{}),
	}
}

// Start starts the worker.
func (w *worker) Start(ctx context.Context) error {
	for {
		select {
		case op := <-w.processQueue:
			if err := w.run(ctx, op); err != nil {
				w.logger.Error("operations worker: %s", err)
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

func (w *worker) run(ctx context.Context, operation *Operation) error {
	statusChan := make(chan *Operation, 10)
	errChan := make(chan error)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				w.logger.Error("panic: %s", r)

				operation.Error("internal error")

				if err := w.db.Update(ctx, operation); err != nil {
					w.logger.Error("failed to update operation: %s", err)
				}
			}
		}()

		errChan <- operation.task(ctx, operation.copy(), statusChan)

		close(errChan)
	}()

	for {
		select {
		case status := <-statusChan:
			if err := w.db.Update(ctx, status); err != nil {
				return fmt.Errorf("failed to update operation: %s", err)
			}
		case err := <-errChan:
			if err == nil {
				return nil
			}

			operation.Error(err.Error())

			if err := w.db.Update(ctx, operation); err != nil {
				return fmt.Errorf("failed to update operation: %s", err)
			}

			return nil
		}
	}
}
