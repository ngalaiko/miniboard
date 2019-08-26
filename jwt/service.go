package jwt // import "miniboard.app/jwt"

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

const (
	defaultIssuer = "miniboard.app"
)

// Service issues and validates jwt tokens.
type Service struct {
	keyStorage *keyStorage

	rotationInterval time.Duration
	expiryInterval   time.Duration

	signer      jose.Signer
	signerGuard *sync.RWMutex
}

// NewService creates new jwt service instance.
func NewService(ctx context.Context, db storage.Storage) *Service {
	keyStorage := newKeyStorage(db)

	s := &Service{
		keyStorage:       keyStorage,
		signerGuard:      &sync.RWMutex{},
		rotationInterval: time.Hour,
		expiryInterval:   time.Hour,
	}

	if err := s.newSigner(); err != nil {
		log("jwt").Panicf("failed to generate encryption key: %s", err)
	}

	go s.rotateKeys(ctx)

	return s
}

func (s *Service) newSigner() error {
	key, err := s.keyStorage.Create()
	if err != nil {
		return err
	}
	log("jwt").Infof("new encryption key: %s", key.ID)

	options := (&jose.SignerOptions{}).
		WithHeader("kid", key.ID).
		WithType("JWT")

	sng, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.ES256,
		Key:       key.Private,
	}, options)
	if err != nil {
		return err
	}

	s.signerGuard.Lock()
	s.signer = sng
	s.signerGuard.Unlock()

	return nil
}

// NewToken returns new authorization.
func (s *Service) NewToken(subject *resource.Name) (string, error) {
	now := time.Now()
	claims := &jwt.Claims{
		ID:       uuid.New().String(),
		Issuer:   defaultIssuer,
		Subject:  subject.String(),
		IssuedAt: jwt.NewNumericDate(now),
		Expiry:   jwt.NewNumericDate(now.Add(s.expiryInterval)),
	}

	s.signerGuard.RLock()
	defer s.signerGuard.RUnlock()
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

func (s *Service) rotateKeys(ctx context.Context) {
	timer := time.NewTicker(s.rotationInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			if err := s.newSigner(); err != nil {
				log("jwt").Errorf("failed to rotate keys: %s", err)
			}
		}
	}
}

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
