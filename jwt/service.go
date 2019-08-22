package jwt // import "miniboard.app/jwt"

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"miniboard.app/storage"
)

const defaultIssuer = "miniboard.app"

// Service issues and validates jwt tokens.
type Service struct {
	keyStorage *keyStorage
	signer     jose.Signer
}

// NewService creates new jwt service instance.
func NewService(db storage.DB) *Service {
	keyStorage := newKeyStorage(db)

	key, err := keyStorage.Create()
	if err != nil {
		log("jwt").Panicf("failed to generate key: %s", err)
	}
	log("jwt").Infof("generated encryption key: %s", key.ID)

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
	}
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

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
