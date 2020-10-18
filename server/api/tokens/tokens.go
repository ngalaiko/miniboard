package tokens

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tokens "github.com/ngalaiko/miniboard/server/genproto/tokens/v1"
	"github.com/ngalaiko/miniboard/server/jwt"
)

// todo: make it shorter
const authDuration = 28 * 24 * time.Hour

type logger interface {
	Error(string, ...interface{})
}

// Service implements tokens service.
type Service struct {
	logger     logger
	jwtService *jwt.Service
}

// NewService returns new serice instance.
func NewService(logger logger, jwt *jwt.Service) *Service {
	return &Service{
		logger:     logger,
		jwtService: jwt,
	}
}

// CreateToken creates new authorization token from a code.
func (s *Service) CreateToken(ctx context.Context, request *tokens.CreateTokenRequest) (*tokens.Token, error) {
	subject, err := s.jwtService.Validate(ctx, request.Code, "authorization_code")
	if err != nil {
		return nil, status.New(codes.InvalidArgument, "authorization code not valid").Err()
	}

	token, err := s.jwtService.NewToken(subject, authDuration, "access_token")
	if err != nil {
		s.logger.Error("failed to generate token :%s", err)
		return nil, status.New(codes.InvalidArgument, "failed to generate token").Err()
	}

	return &tokens.Token{
		Token: token,
	}, nil
}
