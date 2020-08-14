package users

import (
	"context"

	"github.com/ngalaiko/miniboard/server/api/actor"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service controls users resource.
type Service struct {
}

// NewService returns new users storage instance.
func NewService() *Service {
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
