package users

import (
	"context"

	"miniboard.app/proto/users/v1"
	"miniboard.app/storage"
)

// Service controls users resource.
type Service struct {
	usersStorage storage.Storage
}

// New returns new users storage instance.
func New(db storage.Storage) *Service {
	return &Service{
		usersStorage: db,
	}
}

// GetUser returns a user if it exists.
func (s *Service) GetUser(
	ctx context.Context,
	request *users.GetUserRequest,
) (*users.User, error) {
	return &users.User{
		Name: request.Name,
	}, nil
}