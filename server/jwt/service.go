package jwt

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ngalaiko/miniboard/server/storage"
	"github.com/ngalaiko/miniboard/server/storage/resource"
	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	defaultIssuer = "github.com/ngalaiko/miniboard/server"
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

	if err := s.newSigner(ctx); err != nil {
		log().Errorf("failed to generate encryption key: %s", err)
	}

	go s.rotateKeys(ctx)
	go s.cleanupOldKeys(ctx)

	return s
}

func (s *Service) newSigner(ctx context.Context) error {
	key, err := s.keyStorage.Create(ctx)
	if err != nil {
		return err
	}
	log().Infof("%s: new encryption key", key.ID)

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
func (s *Service) Validate(ctx context.Context, raw string, typ string) (*resource.Name, error) {
	token, err := jwt.ParseSigned(raw)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if len(token.Headers) == 0 {
		return nil, fmt.Errorf("headers missing from the token: %w", err)
	}

	id := token.Headers[0].KeyID

	key, err := s.keyStorage.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find key '%s': %w", id, err)
	}

	claims := &jwt.Claims{}
	custom := &customClaims{}
	if err := token.Claims(key.Public, claims, custom); err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}

	if custom.Type != typ {
		return nil, fmt.Errorf("wrong token type, expected '%s': %w", typ, err)
	}

	if err := claims.Validate(jwt.Expected{
		Issuer: defaultIssuer,
		Time:   time.Now(),
	}); err != nil {
		return nil, fmt.Errorf("token is invalid: %w", err)
	}

	return resource.ParseName(claims.Subject), nil
}

func (s *Service) rotateKeys(ctx context.Context) {
	timer := time.NewTicker(s.rotationInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			if err := s.newSigner(ctx); err != nil {
				log().Errorf("failed to rotate keys: %s", err)
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
			if err := s.deleteOldKeys(ctx); err != nil {
				log().Errorf("failed to delete keys: %s", err)
			}
		}
	}
}

func (s *Service) deleteOldKeys(ctx context.Context) error {
	kk, err := s.keyStorage.List(ctx)
	if err != nil {
		return fmt.Errorf("can't list keys from the storage: %w", err)
	}

	deleteBefore := time.Now().Add(-2 * s.expiryInterval)
	for _, k := range kk {
		ksID, err := ksuid.Parse(k.ID)
		if err != nil {
			log().Errorf("%s: malformed kuid", k.ID)
			continue
		}

		if ksID.Time().After(deleteBefore) {
			continue
		}
		if err := s.keyStorage.Delete(ctx, k.ID); err != nil {
			log().Errorf("%s: can't delete key: %s", k.ID, err)
			continue
		}
		log().Infof("%s: key removed", k.ID)
	}

	return nil
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "jwt",
	})
}
