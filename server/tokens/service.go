package tokens

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"miniboard.app/jwt"
	"miniboard.app/proto/tokens/v1"
)

// todo: make it shorter
const authDuration = 28 * 24 * time.Hour

// Service implements tokens service.
type Service struct {
	jwtService *jwt.Service
}

// New returns new serice instance.
func New(jwt *jwt.Service) *Service {
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

	//grpc.SetHeader(ctx, metadata.Pairs(
	//"Set-Cookie", fmt.Sprintf("Set-Cookie: jwt=%s; Secure; HttpOnly", token)),
	//)

	return &tokens.Token{
		Token: token,
	}, nil
}
func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
