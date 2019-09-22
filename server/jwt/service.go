package jwt

import (
	"context"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
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

type customClaims struct {
	Type string `json:"type"`
}

// NewService creates new jwt service instance.
func NewService(ctx context.Context, db storage.Storage) *Service {
	keyStorage := newKeyStorage(db)

	s := &Service{
		keyStorage:       keyStorage,
		signerGuard:      &sync.RWMutex{},
		rotationInterval: 2 * time.Hour,
		expiryInterval:   72 * 3 * time.Hour,
	}

	if err := s.newSigner(); err != nil {
		log("jwt").Panicf("failed to generate encryption key: %s", err)
	}

	go s.rotateKeys(ctx)
	go s.cleanupOldKeys(ctx)

	return s
}

func (s *Service) newSigner() error {
	key, err := s.keyStorage.Create()
	if err != nil {
		return err
	}
	log("jwt").Infof("%s: new encryption key", key.ID)

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
func (s *Service) NewToken(subject *resource.Name, validFor time.Duration, typ string) (string, error) {
	now := time.Now()
	claims := &jwt.Claims{
		ID:       ksuid.New().String(),
		Issuer:   defaultIssuer,
		Subject:  subject.String(),
		IssuedAt: jwt.NewNumericDate(now),
		Expiry:   jwt.NewNumericDate(now.Add(validFor)),
	}

	custom := &customClaims{
		Type: typ,
	}

	s.signerGuard.RLock()
	defer s.signerGuard.RUnlock()
	return jwt.Signed(s.signer).Claims(claims).Claims(custom).CompactSerialize()
}

// Validate returns token subject if a token is valid.
func (s *Service) Validate(raw string, typ string) (string, error) {
	token, err := jwt.ParseSigned(raw)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse token")
	}

	if len(token.Headers) == 0 {
		return "", errors.Wrap(err, "headers missing from the token")
	}

	id := token.Headers[0].KeyID

	key, err := s.keyStorage.Get(id)
	if err != nil {
		return "", errors.Wrapf(err, "failed to find key '%s'", id)
	}

	claims := &jwt.Claims{}
	custom := &customClaims{}
	if err := token.Claims(key.Public, claims, custom); err != nil {
		return "", errors.Wrapf(err, "failed to parse claims")
	}

	if custom.Type != typ {
		return "", errors.Wrapf(err, "wrong token type, expected '%s'", typ)
	}

	if err := claims.Validate(jwt.Expected{
		Issuer: defaultIssuer,
		Time:   time.Now(),
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

func (s *Service) cleanupOldKeys(ctx context.Context) {
	timer := time.NewTicker(s.expiryInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			if err := s.deleteOldKeys(); err != nil {
				log("jwt").Errorf("failed to delete keys: %s", err)
			}
		}
	}
}

func (s *Service) deleteOldKeys() error {
	kk, err := s.keyStorage.List()
	if err != nil {
		return errors.Wrap(err, "can't list keys from the storage")
	}

	deleteBefore := time.Now().Add(-2 * s.expiryInterval)
	for _, k := range kk {
		ksID, err := ksuid.Parse(k.ID)
		if err != nil {
			log("jwt").Errorf("%s: malformed kuid", k.ID)
			continue
		}

		if ksID.Time().After(deleteBefore) {
			continue
		}
		if err := s.keyStorage.Delete(k.ID); err != nil {
			log("jwt").Errorf("%s: can't delete key: %s", k.ID, err)
			continue
		}
		log("jwt").Infof("%s: key removed", k.ID)
	}

	return nil
}

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}