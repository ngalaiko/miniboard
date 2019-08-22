package users // import "miniboard.app/api/users"

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/proto/users/v1"
	"miniboard.app/storage"
)

const bcryptCost = 10

// Service controlls users resource.
type Service struct {
	usersStorage     storage.Storage
	passwordsStorage storage.Storage
}

// New returns new users storage instance.
func New(db storage.DB) *Service {
	return &Service{
		usersStorage:     db.Namespace("users"),
		passwordsStorage: db.Namespace("passwords"),
	}
}

// GetUser returns a user if it exists.
func (s *Service) GetUser(
	ctx context.Context,
	request *users.GetUserRequest,
) (*users.User, error) {
	rawUser, err := s.usersStorage.Load([]byte(request.Name))
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
	if request.Name == "" {
		return nil, status.New(codes.InvalidArgument, "name is empty").Err()
	}

	if request.Password == "" {
		return nil, status.New(codes.InvalidArgument, "password is empty").Err()
	}

	hashName := fmt.Sprintf("%s/password", request.Name)

	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcryptCost)
	if err != nil {
		return nil, status.New(codes.Internal, "failed to calculate password hash").Err()
	}

	if err := s.passwordsStorage.Store([]byte(hashName), hash); err != nil {
		return nil, status.New(codes.Internal, "failed to store password hash").Err()
	}

	user := &users.User{
		Name: request.Name,
	}

	rawUser, err := proto.Marshal(user)
	if err != nil {
		return nil, status.New(codes.Internal, "failed to marshal user").Err()
	}

	if err := s.usersStorage.Store([]byte(user.Name), rawUser); err != nil {
		return nil, status.New(codes.Internal, "failed to store user").Err()
	}

	return user, nil
}
