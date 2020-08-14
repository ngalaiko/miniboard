package codes

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/ngalaiko/miniboard/server/email"
	"github.com/ngalaiko/miniboard/server/jwt"
	"github.com/ngalaiko/miniboard/server/storage/resource"
	"github.com/sirupsen/logrus"
	"github.com/spaolacci/murmur3"
	responsecodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Service implements codes service.
type Service struct {
	domain      string
	jwt         *jwt.Service
	emailClient email.Client
}

// NewService returns new serice instance.
func NewService(
	domain string,
	emailClient email.Client,
	jwt *jwt.Service,
) *Service {
	return &Service{
		domain:      domain,
		emailClient: emailClient,
		jwt:         jwt,
	}
}

// CreateCode creates new authorization code.
func (s *Service) CreateCode(ctx context.Context, request *CreateCodeRequest) (*Code, error) {
	h := murmur3.New128()
	_, _ = io.WriteString(h, request.Email)
	hashedEmail := fmt.Sprintf("%x", h.Sum(nil))

	token, err := s.jwt.NewToken(resource.NewName("users", hashedEmail), 10*time.Minute, "authorization_code")
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
			log("codes").Error(err)
		}
	}(msg)

	return &Code{}, nil
}

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
