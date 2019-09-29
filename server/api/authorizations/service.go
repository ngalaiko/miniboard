package authorizations

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/jwt"
	"miniboard.app/proto/authorizations/v1"
	"miniboard.app/storage/resource"
)

const (
	accessToken       = "access"
	refreshToken      = "refresh"
	authorizationCode = "authorization_code"
)

// Service creates and validates new authorizations.
type Service struct {
	jwt *jwt.Service
}

// New creates a new service instance.
func New(jwtService *jwt.Service) *Service {
	return &Service{
		jwt: jwtService,
	}
}

// CreateAuthorization returns new JWT authentiction.
func (s *Service) CreateAuthorization(
	ctx context.Context,
	request *authorizations.CreateAuthorizationRequest,
) (*authorizations.Authorization, error) {
	switch request.GrantType {
	case "refresh_token":
		return s.refreshTokenAuthorization(request.RefreshToken)
	case "authorization_code":
		return s.authentictionCode(request.AuthorizationCode)
	default:
		return nil, status.New(codes.InvalidArgument, "unknown grant type").Err()
	}
}

func (s *Service) refreshTokenAuthorization(token string) (*authorizations.Authorization, error) {
	subject, err := s.jwt.Validate(token, refreshToken)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, "refresh token is not valid").Err()
	}

	return s.authorizationFor(resource.ParseName(subject))
}

func (s *Service) authentictionCode(token string) (*authorizations.Authorization, error) {
	subject, err := s.jwt.Validate(token, authorizationCode)
	if err != nil {
		return nil, status.New(codes.InvalidArgument, "refresh token is not valid").Err()
	}

	return s.authorizationFor(resource.ParseName(subject))
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
