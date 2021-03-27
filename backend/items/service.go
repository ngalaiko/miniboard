package items

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Known errors.
var (
	errNotFound      = fmt.Errorf("not found")
	errURLIsEmpty    = fmt.Errorf("item url must not be empty")
	errTitleIsEmpty  = fmt.Errorf("item title must not be empty")
	ErrAlreadyExists = fmt.Errorf("item already exists")
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
	subscriptionID string,
	url string,
	title string,
	date time.Time,
	summary string,
) (*Item, error) {
	if exists, err := s.db.GetByURL(ctx, url); err == nil && exists != nil {
		return nil, ErrAlreadyExists
	}

	if url == "" {
		return nil, errURLIsEmpty
	}

	if title == "" {
		return nil, errTitleIsEmpty
	}

	item := &Item{
		ID:             uuid.New().String(),
		URL:            url,
		Title:          title,
		SubscriptionID: subscriptionID,
		Created:        date,
		Summary:        summary,
	}

	if err := s.db.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to store item %+v: %w", item, err)
	}

	return item, nil
}

// Get returns a item by it's id.
func (s *Service) Get(ctx context.Context, id string, userID string) (*UserItem, error) {
	item, err := s.db.Get(ctx, userID, id)
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
func (s *Service) List(ctx context.Context, userID string, pageSize int, createdLT *time.Time, subscriptionID *string, tagID *string) ([]*UserItem, error) {
	items, err := s.db.List(ctx, userID, pageSize, createdLT, subscriptionID, tagID)
	switch {
	case err == nil:
		return items, nil
	default:
		return nil, err
	}
}
