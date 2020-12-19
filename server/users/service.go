package users

import (
	"context"
	"database/sql"
	"fmt"
)

// Known errors.
var (
	ErrNotFound = fmt.Errorf("user does not exist")
	// TODO: implement
	ErrAlreadyExists = fmt.Errorf("user already exists")
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
func (s *Service) Create(ctx context.Context, username string, password []byte) (*User, error) {
	_, err := s.db.GetByUsername(ctx, username)
	switch err {
	case nil:
		return nil, ErrAlreadyExists
	case sql.ErrNoRows:
	default:
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user, err := newUser(username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if err := s.db.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create used in the database: %w", err)
	}

	return user, nil
}

// GetByID returns a user by id.
func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	user, err := s.db.GetByID(ctx, id)
	switch err {
	case nil:
		return user, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("failed to get user from db: %w", err)
	}
}

// GetByUsername returns a user by id.
func (s *Service) GetByUsername(ctx context.Context, username string) (*User, error) {
	user, err := s.db.GetByUsername(ctx, username)
	switch err {
	case nil:
		return user, nil
	case sql.ErrNoRows:
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("failed to get user from db: %w", err)
	}
}
