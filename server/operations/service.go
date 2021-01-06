package operations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var (
	errNotFound = fmt.Errorf("not found")
)

type logger interface {
	Error(string, ...interface{})
}

// Config contains configuration for the operations service.
type Config struct {
	Workers int `yaml:"workers"`
}

// Service is used to manage longrunning opeartions.
type Service struct {
	db     *database
	logger logger

	processQueue chan *Operation
}

// NewService returns new service.
func NewService(ctx context.Context, logger logger, sqldb *sql.DB, cfg *Config) *Service {
	if cfg == nil {
		cfg = &Config{}
	}

	s := &Service{
		db:           newDatabase(sqldb),
		logger:       logger,
		processQueue: make(chan *Operation),
	}

	nWorkers := 10
	if cfg.Workers > 0 {
		nWorkers = cfg.Workers
	}

	for i := 0; i < nWorkers; i++ {
		go s.runWorker(ctx)
	}

	return s
}

// Create creates an operation, and runs it.
func (s *Service) Create(ctx context.Context, userID string, opFunc operationFunc) (*Operation, error) {
	operation := &Operation{
		ID:     uuid.New().String(),
		UserID: userID,
		task:   opFunc,
	}

	if err := s.db.Create(ctx, operation); err != nil {
		return nil, fmt.Errorf("failed to create operation: %w", err)
	}

	s.processQueue <- operation

	return operation, nil
}

func (s *Service) runWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case op := <-s.processQueue:
			s.runOperation(ctx, op)
		}
	}
}

func (s *Service) runOperation(ctx context.Context, operation *Operation) {
	statusChan := make(chan *Operation, 10)

	errChan := make(chan error)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error("panic: %s", r)

				operation.Error("internal error")

				if err := s.db.Update(ctx, operation); err != nil {
					s.logger.Error("failed to update operation: %s", err)
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
			if err := s.db.Update(ctx, status); err != nil {
				s.logger.Error("failed to update operation: %s", err)
			}
		case err := <-errChan:
			switch {
			case err != nil:
				operation.Error(err.Error())
			default:
				return
			}

			if err := s.db.Update(ctx, operation); err != nil {
				s.logger.Error("failed to update operation: %s", err)
			}
		}
	}
}

// Get returns a single operation.
func (s *Service) Get(ctx context.Context, id string, userID string) (*Operation, error) {
	operation, err := s.db.Get(ctx, id, userID)
	switch {
	case err == nil:
		return operation, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, errNotFound
	default:
		return nil, err
	}
}
