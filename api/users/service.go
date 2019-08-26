package users // import "miniboard.app/api/users"

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/passwords"
	"miniboard.app/proto/users/v1"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

// Service controlls users resource.
type Service struct {
	usersStorage     storage.Storage
	passwordsService *passwords.Service
}

// New returns new users storage instance.
func New(db storage.Storage, passwordsService *passwords.Service) *Service {
	return &Service{
		usersStorage:     db,
		passwordsService: passwordsService,
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

// CreateUser creates a new user.
func (s *Service) CreateUser(
	ctx context.Context,
	request *users.CreateUserRequest,
) (*users.User, error) {
	if request.Username == "" {
		return nil, status.New(codes.InvalidArgument, "name is empty").Err()
	}

	if request.Password == "" {
		return nil, status.New(codes.InvalidArgument, "password is empty").Err()
	}

	name := resource.NewName("users", request.Username)

	user := &users.User{
		Name: name.String(),
	}

	rawUser, err := proto.Marshal(user)
	if err != nil {
		return nil, status.New(codes.Internal, "failed to marshal user").Err()
	}

	switch errors.Cause(s.usersStorage.Store(name, rawUser)) {
	case nil:
	case storage.ErrAlreadyExists:
		return nil, status.New(codes.AlreadyExists, "user already exists").Err()
	default:
		return nil, status.New(codes.Internal, "failed to store user").Err()
	}

	if err := s.passwordsService.Set(name, request.Password); err != nil {
		return nil, status.New(codes.Internal, "failed to store password hash").Err()
	}

	return user, nil
}
