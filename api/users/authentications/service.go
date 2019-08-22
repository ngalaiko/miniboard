package authentications // import "miniboard.app/services/authentications"

import (
	"context"

	"github.com/pkg/errors"
	"miniboard.app/jwt"
	"miniboard.app/proto/users/authentications/v1"
	"miniboard.app/storage"
)

// Service creates and validates new authorizations.
type Service struct {
	jwt *jwt.Service
}

// New creates a new service instance.
func New(
	ctx context.Context,
	db storage.DB,
) (*Service, error) {
	jwt, err := jwt.New(ctx, db)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init jwt service")
	}
	return &Service{
		jwt: jwt,
	}, nil
}

// CreateAuthentication returns new JWT authentiction.
func (s *Service) CreateAuthentication(
	ctx context.Context,
	request *authentications.CreateAuthenticationRequest,
) (*authentications.Authentication, error) {
	return nil, nil
}
