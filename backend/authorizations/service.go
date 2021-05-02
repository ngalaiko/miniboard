package authorizations

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	jose "gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"

	"github.com/ngalaiko/miniboard/backend/authorizations/keys"
)

// Known errors.
var (
	ErrInvalidToken = fmt.Errorf("token is invalid")
	ErrTokenExpired = fmt.Errorf("token is expired")
)

const (
	defaultIssuer = "miniboard.app"
	validFor      = 15 * time.Minute
)

type logger interface {
	Error(string, ...interface{})
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

	s.signer = signer

	return nil
}

// NewToken creates a new signed JWT.
func (s *Service) NewToken(ctx context.Context, userID string) (*Token, error) {
	if s.signer == nil {
		if err := s.Init(ctx); err != nil {
			return nil, err
		}
	}

	now := time.Now()
	claims := &jwt.Claims{
		ID:       uuid.New().String(),
		Issuer:   defaultIssuer,
		Subject:  userID,
		IssuedAt: jwt.NewNumericDate(now),
		Expiry:   jwt.NewNumericDate(now.Add(validFor)),
	}

	token, err := jwt.Signed(s.signer).Claims(claims).CompactSerialize()
	if err != nil {
		return nil, fmt.Errorf("failed to create a signed token: %w", err)
	}

	return &Token{
		Token:     token,
		UserID:    userID,
		ExpiresAt: claims.Expiry.Time(),
	}, nil
}

// Verify checks token signature and returns it's meaningful content.
func (s *Service) Verify(ctx context.Context, token string) (*Token, error) {
	jwtoken, err := jwt.ParseSigned(token)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if len(jwtoken.Headers) == 0 {
		return nil, ErrInvalidToken
	}

	id := jwtoken.Headers[0].KeyID

	pubicKey, err := s.get(ctx, id)
	switch {
	case err == nil:
	case errors.Is(err, sql.ErrNoRows):
		return nil, ErrInvalidToken
	default:
		return nil, fmt.Errorf("failed to find key '%s': %w", id, err)
	}

	claims := &jwt.Claims{}
	if err := jwtoken.Claims(pubicKey, claims); err != nil {
		return nil, ErrInvalidToken
	}

	validateErr := claims.ValidateWithLeeway(jwt.Expected{
		Time:   time.Now(),
		Issuer: defaultIssuer,
	}, time.Second)
	switch {
	case validateErr == nil:
		return &Token{
			Token:     token,
			UserID:    claims.Subject,
			ExpiresAt: claims.Expiry.Time(),
		}, nil
	case errors.Is(validateErr, jwt.ErrExpired):
		return &Token{
			Token:     token,
			UserID:    claims.Subject,
			ExpiresAt: claims.Expiry.Time(),
		}, ErrTokenExpired
	default:
		return nil, ErrInvalidToken
	}
}

func (s *Service) get(ctx context.Context, id string) (*ecdsa.PublicKey, error) {
	key, err := s.keysDatabase.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find key '%s': %w", id, err)
	}

	untypedResult, err := x509.ParsePKIXPublicKey(key.PublicDER)
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
