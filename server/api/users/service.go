package users

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/proto/users/v1"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
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
	rawUser, err := s.usersStorage.Load(resource.ParseName(request.Name))
	switch errors.Cause(err) {
	case nil:
	case storage.ErrNotFound:
		return nil, status.New(codes.NotFound, "user not found").Err()
	default:
		return nil, status.New(codes.Internal, "failed to find user").Err()
	}

	user := &users.User{}
	if err := proto.Unmarshal(rawUser, user); err != nil {
		return nil, status.New(codes.Internal, "failed to unmarshal user").Err()
	}

	return user, nil
}
