package codes

import (
	"context"
	"fmt"
	"io"
	"time"

	codes "github.com/ngalaiko/miniboard/server/genproto/codes/v1"
	"github.com/ngalaiko/miniboard/server/jwt"
	"github.com/spaolacci/murmur3"
	responsecodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type emailSender interface {
	Send(to string, subject string, payload string) error
}

type logger interface {
	Error(string, ...interface{})
}

// Service implements codes service.
type Service struct {
	logger      logger
	domain      string
	jwt         *jwt.Service
	emailClient emailSender
}

// NewService returns new serice instance.
func NewService(
	logger logger,
	domain string,
	emailClient emailSender,
	jwt *jwt.Service,
) *Service {
	return &Service{
		logger:      logger,
		domain:      domain,
		emailClient: emailClient,
		jwt:         jwt,
	}
}

// CreateCode creates new authorization code.
func (s *Service) CreateCode(ctx context.Context, request *codes.CreateCodeRequest) (*codes.Code, error) {
	h := murmur3.New128()
	_, _ = io.WriteString(h, request.Email)
	hashedEmail := fmt.Sprintf("%x", h.Sum(nil))

	token, err := s.jwt.NewToken(hashedEmail, 10*time.Minute, "authorization_code")
	if err != nil {
		return nil, status.New(responsecodes.Internal, "failed to generate token").Err()
	}

	link := fmt.Sprintf("%s/codes?code=%s", s.domain, token)

	msg := fmt.Sprintf(`
Follow the link or copy code to the login form.

Code: %s

Link: %s
`, token, link)

	go func(msg string) {
		if err := s.emailClient.Send(request.Email, "Authentication link", msg); err != nil {
			s.logger.Error("%s", err)
		}
	}(msg)

	return &codes.Code{}, nil
}