package operations

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

var (
	errNotFound = fmt.Errorf("not found")
)

type logger interface {
	Info(string, ...interface{})
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
	config *Config

	processQueue chan *Operation
	workers      []*worker
}

// NewService returns new service.
func NewService(logger logger, sqldb *sql.DB, cfg *Config) *Service {
	if cfg == nil {
		cfg = &Config{}
	}

	return &Service{
		db:           newDatabase(sqldb),
		logger:       logger,
		processQueue: make(chan *Operation),
		config:       cfg,
	}
}

// Start starts processing workers.
func (s *Service) Start(ctx context.Context) error {
	nWorkers := 10
	if s.config.Workers > 0 {
		nWorkers = s.config.Workers
	}

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < nWorkers; i++ {
		worker := newWorker(s.processQueue, s.logger, s.db)
		s.workers = append(s.workers, worker)

		g.Go(restartingWorker(ctx, worker, s.logger))
	}

	return g.Wait()
}

func restartingWorker(ctx context.Context, worker *worker, logger logger) func() error {
	return func() error {
		if err := worker.Start(ctx); err != nil {
			logger.Error("operations worker: %s", err)
			return restartingWorker(ctx, worker, logger)()
		}
		return nil
	}
}

// Shutdown stops all running workers.
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

// Create creates an operation, and runs it.
func (s *Service) Create(ctx context.Context, userID string, opFunc Task) (*Operation, error) {
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
