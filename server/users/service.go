package users

import (
	"context"
	"database/sql"
	"fmt"
)

// Known errors.
var (
	ErrNotFound = fmt.Errorf("user does not exist")
)

// Service allows to manage users resource.
type Service struct {
	db *database
}

// NewService returns new users service.
func NewService(db *sql.DB) *Service {
	return &Service{
		db: newDB(db),
	}
}

// Create creates a new user with the given password.
func (s *Service) Create(ctx context.Context, password []byte) (*User, error) {
	user, err := newUser(password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err := s.db.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create used in the database: %w", err)
	}

	return user, nil
}

// Get returns a user by id.
func (s *Service) Get(ctx context.Context, id string) (*User, error) {
	user, err := s.db.Get(ctx, id)
	switch err {
	case nil:
		return user, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("failed to get user from db: %w", err)
	}
}
