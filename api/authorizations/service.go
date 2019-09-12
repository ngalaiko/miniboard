package authorizations

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/jwt"
	"miniboard.app/passwords"
	"miniboard.app/proto/authorizations/v1"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

const (
	accessToken  = "access"
	refreshToken = "refresh"
)

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

// CreateAuthorization returns new JWT authentiction.
func (s *Service) CreateAuthorization(
	ctx context.Context,
	request *authorizations.CreateAuthorizationRequest,
) (*authorizations.Authorization, error) {
	switch request.GrantType {
	case "password":
		return s.passwordAuthorization(resource.NewName("users", request.Username), request.Password)
	case "refresh_token":
		return s.refreshTokenAuthorization(resource.NewName("users", request.Username), request.RefreshToken)
	default:
		return nil, status.New(codes.InvalidArgument, "unknown grant type").Err()
	}
}

func (s *Service) passwordAuthorization(user *resource.Name, password string) (*authorizations.Authorization, error) {
	valid, err := s.passwords.Validate(user, password)
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

	return s.authorizationFor(user)
}

func (s *Service) refreshTokenAuthorization(user *resource.Name, token string) (*authorizations.Authorization, error) {
	subject, err := s.jwt.Validate(token, refreshToken)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, "refresh token is not valid").Err()
	}

	if subject != user.String() {
		return nil, status.New(codes.InvalidArgument, "refresh token: wrong subject").Err()
	}

	return s.authorizationFor(user)
}

func (s *Service) authorizationFor(user *resource.Name) (*authorizations.Authorization, error) {
	accessToken, err := s.jwt.NewToken(user, 3*time.Hour, accessToken)
	if err != nil {
		return nil, status.New(codes.Internal, "failed to generage token").Err()
	}

	refreshToken, err := s.jwt.NewToken(user, 72*time.Hour, refreshToken)
	if err != nil {
		return nil, status.New(codes.Internal, "failed to generage token").Err()
	}

	return &authorizations.Authorization{
		TokenType:    "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
