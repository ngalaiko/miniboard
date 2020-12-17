package jwt

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"database/sql"
	"fmt"

	jose "gopkg.in/square/go-jose.v2"

	"github.com/ngalaiko/miniboard/server/jwt/keys"
)

type logger interface {
	Info(string, ...interface{})
}

// Service allows to issue and verify jwt tokens.
type Service struct {
	logger       logger
	keysDatabase *keys.Database

	signer jose.Signer
}

// NewService creates a new jwt service.
func NewService(db *sql.DB, logger logger) *Service {
	return &Service{
		logger:       logger,
		keysDatabase: keys.NewDatabase(db),
	}
}

// Init prepares jwt service for work by generating a jwt signer.
func (s *Service) Init(ctx context.Context) error {
	if s.signer != nil {
		return fmt.Errorf("already initialized")
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate encryption key: %w", err)
	}

	publicDER, err := x509.MarshalPKIXPublicKey(privateKey.Public())
	if err != nil {
		return fmt.Errorf("failed to marshal encryption key: %w", err)
	}

	key, err := keys.New(publicDER)
	if err != nil {
		return fmt.Errorf("failed to create a key: %w", err)
	}

	if err := s.keysDatabase.Create(ctx, key); err != nil {
		return fmt.Errorf("failed to store key in the database: %w", err)
	}

	options := (&jose.SignerOptions{}).
		WithHeader("kid", key.ID).
		WithType("JWT")

	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.ES256,
		Key:       privateKey,
	}, options)
	if err != nil {
		return err
	}

	s.logger.Info("new signer with id '%s' created", key.ID)

	s.signer = signer

	return nil
}
