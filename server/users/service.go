package users

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/api/actor"
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
	request *GetMeRequest,
) (*User, error) {
	actor, ok := actor.FromContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "not authenticated")
	}
	return &User{
		Name: actor.String(),
	}, nil
}
