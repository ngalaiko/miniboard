package users

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"miniboard.app/api/actor"
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

// GetMe returns authenticated user.
func (s *Service) GetMe(
	ctx context.Context,
	request *users.GetMeRequest,
) (*users.User, error) {
	actor, ok := actor.FromContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.NotFound, "not authenticated")
	}
	return &users.User{
		Name: actor.String(),
	}, nil
}
