package jwt

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"miniboard.app/application/storage"
)

// Service issues and validates jwt tokens.
type Service struct {
	keyStorage *keyStorage
	signer     jose.Signer
}

// New creates new jwt service instance.
func New(_ context.Context, publicKeyStorage storage.Storage) (*Service, error) {
	keyStorage := newKeyStorage(publicKeyStorage)

	key, err := keyStorage.Create()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate key")
	}
	logrus.Infof("generated key: %s", key.ID)

	options := (&jose.SignerOptions{}).
		WithHeader("kid", key.ID).
		WithType("JWT")

	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.ES256,
		Key:       key.Private,
	}, options)

	return &Service{
		keyStorage: keyStorage,
		signer:     signer,
	}, nil
}

// NewToken returns new authorization.
func (s *Service) NewToken(subject string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := &jwt.Claims{
		ID:       uuid.New().String(),
		Issuer:   "miniboard.app",
		Subject:  subject,
		IssuedAt: jwt.NewNumericDate(now),
		Expiry:   jwt.NewNumericDate(now.Add(duration)),
	}

	return jwt.Signed(s.signer).Claims(claims).CompactSerialize()
}
