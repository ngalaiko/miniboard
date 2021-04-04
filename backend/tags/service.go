package tags

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
	ErrNotFound      = fmt.Errorf("not found")
	ErrAlreadyExists = fmt.Errorf("tag already exists")
)

// Service allows to manage tags resource.
type Service struct {
	db *database
}

// NewService returns new tags service.
func NewService(db *sql.DB) *Service {
	return &Service{
		db: newDB(db),
	}
}

// Create creates a tag.
func (s *Service) Create(ctx context.Context, userID string, title string) (*Tag, error) {
	if exists, err := s.db.GetByTitle(ctx, userID, title); err == nil && exists != nil {
		return nil, ErrAlreadyExists
	}

	tag := &Tag{
		ID:      uuid.New().String(),
		UserID:  userID,
		Title:   title,
		Created: time.Now().Truncate(time.Nanosecond),
	}

	if err := s.db.Create(ctx, tag); err != nil {
		return nil, fmt.Errorf("failed to store tag: %w", err)
	}

	return tag, nil
}

// Get returns a tag by it's id.
func (s *Service) Get(ctx context.Context, id string, userID string) (*Tag, error) {
	tag, err := s.db.GetByID(ctx, id, userID)
	switch {
	case err == nil:
		return tag, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// GetByTitle returns a tag by it's id.
func (s *Service) GetByTitle(ctx context.Context, userID string, title string) (*Tag, error) {
	tag, err := s.db.GetByTitle(ctx, userID, title)
	switch {
	case err == nil:
		return tag, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// List returns a list of user tags.
func (s *Service) List(ctx context.Context, userID string, pageSize int, createdLT *time.Time) ([]*Tag, error) {
	tags, err := s.db.List(ctx, userID, pageSize, createdLT)
	switch {
	case err == nil:
		return tags, nil
	default:
		return nil, err
	}
}
