package tokens

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	tokens "github.com/ngalaiko/miniboard/server/genproto/tokens/v1"
	"github.com/ngalaiko/miniboard/server/jwt"
)

// todo: make it shorter
const authDuration = 28 * 24 * time.Hour

// Service implements tokens service.
type Service struct {
	jwtService *jwt.Service
}

// NewService returns new serice instance.
func NewService(jwt *jwt.Service) *Service {
	return &Service{
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
		log("tokens").Errorf("failed to generate token :%s", err)
		return nil, status.New(codes.InvalidArgument, "failed to generate token").Err()
	}

	return &tokens.Token{
		Token: token,
	}, nil
}
func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
