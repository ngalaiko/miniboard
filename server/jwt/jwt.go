package jwt

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	"github.com/ngalaiko/miniboard/server/jwt/db"
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
	keyStorage *db.DB

	rotationInterval time.Duration
	expiryInterval   time.Duration

	signer      jose.Signer
	signerGuard *sync.RWMutex
}

type customClaims struct {
	Type string `json:"type"`
}

// NewService creates new jwt service instance.
func NewService(ctx context.Context, sqldb *sql.DB) *Service {
	s := &Service{
		keyStorage:       db.New(sqldb),
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
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate encryption key: %w", err)
	}

	der, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return fmt.Errorf("failed to marshal encryption key: %w", err)
	}

	publicKey := &db.PublicKey{
		ID:        ksuid.New().String(),
		DerBase64: base64.StdEncoding.EncodeToString(der),
	}

	if err := s.keyStorage.Create(ctx, publicKey); err != nil {
		return err
	}
	log().Infof("%s: new encryption key", publicKey.ID)

	options := (&jose.SignerOptions{}).
		WithHeader("kid", publicKey.ID).
		WithType("JWT")

	sng, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.ES256,
		Key:       privateKey,
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
func (s *Service) NewToken(subject string, validFor time.Duration, typ string) (string, error) {
	now := time.Now()
	claims := &jwt.Claims{
		ID:       ksuid.New().String(),
		Issuer:   defaultIssuer,
		Subject:  subject,
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
func (s *Service) Validate(ctx context.Context, raw string, typ string) (string, error) {
	token, err := jwt.ParseSigned(raw)
	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if len(token.Headers) == 0 {
		return "", fmt.Errorf("headers missing from the token: %w", err)
	}

	id := token.Headers[0].KeyID

	pubicKey, err := s.get(ctx, id)
	if err != nil {
		return "", fmt.Errorf("failed to find key '%s': %w", id, err)
	}

	claims := &jwt.Claims{}
	custom := &customClaims{}
	if err := token.Claims(pubicKey, claims, custom); err != nil {
		return "", fmt.Errorf("failed to parse claims: %w", err)
	}

	if custom.Type != typ {
		return "", fmt.Errorf("wrong token type, expected '%s': %w", typ, err)
	}

	if err := claims.Validate(jwt.Expected{
		Issuer: defaultIssuer,
		Time:   time.Now(),
	}); err != nil {
		return "", fmt.Errorf("token is invalid: %w", err)
	}

	return claims.Subject, nil
}

func (s *Service) get(ctx context.Context, id string) (*ecdsa.PublicKey, error) {
	key, err := s.keyStorage.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find key '%s': %w", id, err)
	}

	der, err := base64.StdEncoding.DecodeString(key.DerBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	untypedResult, err := x509.ParsePKIXPublicKey(der)
	if err != nil {
		return nil, fmt.Errorf("unable to parse PKIX public key: %w", err)
	}

	switch v := untypedResult.(type) {
	case *ecdsa.PublicKey:
		return v, nil
	default:
		return nil, fmt.Errorf("unknown public key type: %T", v)
	}
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
