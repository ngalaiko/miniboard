package operations

import (
	"context"
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

func (w *worker) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case op := <-w.processQueue:
				w.run(ctx, op)
			case <-w.shutdown:
				close(w.stopped)
				return
			}
		}
	}()
}

func (w *worker) Shutdown(ctx context.Context) {
	close(w.shutdown)
	<-w.stopped
}

func (w *worker) run(ctx context.Context, operation *Operation) {
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

		if err := operation.task(ctx, operation.copy(), statusChan); err != nil {
			errChan <- err
		}
		close(errChan)
	}()

	for {
		select {
		case status := <-statusChan:
			if err := w.db.Update(ctx, status); err != nil {
				w.logger.Error("failed to update operation: %s", err)
			}
		case err := <-errChan:
			switch {
			case err != nil:
				operation.Error(err.Error())
			default:
				return
			}

			if err := w.db.Update(ctx, operation); err != nil {
				w.logger.Error("failed to update operation: %s", err)
			}
		}
	}
}
