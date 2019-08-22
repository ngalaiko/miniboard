package authentications // import "miniboard.app/services/authentications"

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/jwt"
	"miniboard.app/passwords"
	"miniboard.app/proto/users/authentications/v1"
	"miniboard.app/storage"
)

const tokenDuration = time.Hour

// Service creates and validates new authorizations.
type Service struct {
	jwt       *jwt.Service
	passwords *passwords.Service
}

// New creates a new service instance.
func New(jwtService *jwt.Service, passwordsService *passwords.Service) *Service {
	return &Service{
		jwt:       jwtService,
		passwords: passwordsService,
	}
}

// CreateAuthentication returns new JWT authentiction.
func (s *Service) CreateAuthentication(
	ctx context.Context,
	request *authentications.CreateAuthenticationRequest,
) (*authentications.Authentication, error) {
	valid, err := s.passwords.Validate(request.Parent, request.Password)
	switch errors.Cause(err) {
	case nil:
	case storage.ErrNotFound:
		return nil, status.New(codes.NotFound, "user not found").Err()
	default:
		return nil, status.New(codes.Internal, "failed to validate password").Err()
	}

	if !valid {
		return nil, status.New(codes.InvalidArgument, "password is not valid").Err()
	}

	token, err := s.jwt.NewToken(request.Parent, tokenDuration)
	if err != nil {
		return nil, status.New(codes.Internal, "failed to generage token").Err()
	}

	return &authentications.Authentication{
		Type:  authentications.Authentication_TYPE_BEARER,
		Token: token,
	}, nil
}
