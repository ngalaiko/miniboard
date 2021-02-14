package items

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
)

// Known errors.
var (
	errNotFound      = fmt.Errorf("not found")
	errAlreadyExists = fmt.Errorf("item already exists")
)

// Service allows to manage items resource.
type Service struct {
	db *database
}

// NewService returns new items service.
func NewService(db *sql.DB, logger logger) *Service {
	return &Service{
		db: newDB(db, logger),
	}
}

// Create creates a item from URL.
func (s *Service) Create(
	ctx context.Context,
	userID string,
	subscriptionID string,
	url *url.URL,
	title string,
) (*Item, error) {
	if exists, err := s.db.GetByURL(ctx, userID, url.String()); err == nil && exists != nil {
		return nil, errAlreadyExists
	}

	item := &Item{
		ID:             uuid.New().String(),
		UserID:         userID,
		URL:            url.String(),
		Title:          title,
		SubscriptionID: subscriptionID,
		Created:        time.Now().Truncate(time.Nanosecond),
	}

	if err := s.db.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to store item: %w", err)
	}

	return item, nil
}

// Get returns a item by it's id.
func (s *Service) Get(ctx context.Context, id string, userID string) (*Item, error) {
	item, err := s.db.Get(ctx, id, userID)
	switch {
	case err == nil:
		return item, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, errNotFound
	default:
		return nil, err
	}
}

// List returns a list of user items.
func (s *Service) List(ctx context.Context, userID string, pageSize int, createdLT *time.Time, subscriptionID *string) ([]*Item, error) {
	items, err := s.db.List(ctx, userID, pageSize, createdLT, subscriptionID)
	switch {
	case err == nil:
		return items, nil
	default:
		return nil, err
	}
}
