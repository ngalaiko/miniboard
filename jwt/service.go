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

const defaultIssuer = "miniboard.app"

// Service issues and validates jwt tokens.
type Service struct {
	keyStorage *keyStorage
	signer     jose.Signer
}

// New creates new jwt service instance.
func New(_ context.Context, db storage.DB) (*Service, error) {
	keyStorage := newKeyStorage(db)

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
		Issuer:   defaultIssuer,
		Subject:  subject,
		IssuedAt: jwt.NewNumericDate(now),
		Expiry:   jwt.NewNumericDate(now.Add(duration)),
	}

	return jwt.Signed(s.signer).Claims(claims).CompactSerialize()
}

// Validate returns token subject if a token is valid.
func (s *Service) Validate(raw string) (string, error) {
	token, err := jwt.ParseSigned(raw)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse token")
	}

	if len(token.Headers) == 0 {
		return "", errors.Wrap(err, "headers missing from the token")
	}

	id, err := uuid.Parse(token.Headers[0].KeyID)
	if err != nil {
		return "", errors.Wrap(err, "invalid id")
	}

	key, err := s.keyStorage.Get(id)
	if err != nil {
		return "", errors.Wrapf(err, "failed to find key '%s'", id)
	}

	claims := &jwt.Claims{}
	if err := token.Claims(key.Public, claims); err != nil {
		return "", errors.Wrapf(err, "failed to parse claims")
	}

	if err := claims.Validate(jwt.Expected{
		Issuer: defaultIssuer,
	}); err != nil {
		return "", errors.Wrap(err, "token is invalid")
	}

	return claims.Subject, nil
}
