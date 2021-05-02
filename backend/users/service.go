package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Known errors.
var (
	ErrNotFound      = fmt.Errorf("user does not exist")
	ErrAlreadyExists = fmt.Errorf("user already exists")
)

// Config contains users configuration.
type Config struct {
	BCryptCost int `yaml:"bcrypt_cost"`
}

// Service allows to manage users resource.
type Service struct {
	db     *database
	config *Config
}

// NewService returns new users service.
func NewService(db *sql.DB, cfg *Config) *Service {
	if cfg == nil {
		cfg = &Config{}
	}
	return &Service{
		db:     newDB(db),
		config: cfg,
	}
}

// Create creates a new user with the given password.
func (s *Service) Create(ctx context.Context, username string, password []byte) (*User, error) {
	_, err := s.db.GetByUsername(ctx, username)
	switch {
	case err == nil:
		return nil, ErrAlreadyExists
	case errors.Is(err, sql.ErrNoRows):
	default:
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	bgCryptCost := 14
	if s.config.BCryptCost > 0 {
		bgCryptCost = s.config.BCryptCost
	}
	user, err := newUser(username, password, bgCryptCost)
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
	switch {
	case err == nil:
		return user, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("failed to get user from db: %w", err)
	}
}

// GetByUsername returns a user by id.
func (s *Service) GetByUsername(ctx context.Context, username string) (*User, error) {
	user, err := s.db.GetByUsername(ctx, username)
	switch {
	case err == nil:
		return user, nil
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrNotFound
	default:
		return nil, fmt.Errorf("failed to get user from db: %w", err)
	}
}
