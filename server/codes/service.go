package codes

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spaolacci/murmur3"
	responsecodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"miniboard.app/email"
	"miniboard.app/jwt"
	"miniboard.app/proto/codes/v1"
	"miniboard.app/storage/resource"
)

// Service implements codes service.
type Service struct {
	domain      string
	jwt         *jwt.Service
	emailClient email.Client
}

// New returns new serice instance.
func New(
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
func (s *Service) CreateCode(ctx context.Context, request *codes.CreateCodeRequest) (*codes.Code, error) {
	h := murmur3.New128()
	_, _ = io.WriteString(h, request.Email)
	hashedEmail := fmt.Sprintf("%x", h.Sum(nil))

	token, err := s.jwt.NewToken(resource.NewName("users", hashedEmail), 10*time.Minute, "authorization_code")
	if err != nil {
		return nil, status.New(responsecodes.Internal, "failed to generate token").Err()
	}

	link := fmt.Sprintf("%s/codes/%s", s.domain, token)

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

	return &codes.Code{}, nil
}

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
