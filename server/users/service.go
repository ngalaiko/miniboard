package users

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/api/actor"
	"miniboard.app/proto/users/v1"
)

// Service controls users resource.
type Service struct {
}

// New returns new users storage instance.
func New() *Service {
	return &Service{}
}

// GetMe returns authenticated user.
func (s *Service) GetMe(
	ctx context.Context,
	request *users.GetMeRequest,
) (*users.User, error) {
	actor, ok := actor.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "not authenticated")
	}
	return &users.User{
		Name: actor.String(),
	}, nil
}
