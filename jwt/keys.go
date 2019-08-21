package jwt // import "miniboard.app/jwt"

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"miniboard.app/application/storage"
)

// key is an encryption key.
type key struct {
	ID      uuid.UUID
	Private crypto.PrivateKey
	Public  crypto.PublicKey
}

type keyStorage struct {
	storage storage.Storage
}

func newKeyStorage(db storage.DB) *keyStorage {
	return &keyStorage{
		storage: db.Namespace("jwt_keys"),
	}
}

// Get returns a key by id.
func (s *keyStorage) Get(id uuid.UUID) (*key, error) {
	idBytes, err := id.MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal key id")
	}

	data, err := s.storage.Load(idBytes)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load key")
	}

	privateKey, err := x509.ParseECPrivateKey(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse key")
	}

	return &key{
		ID:      id,
		Private: privateKey,
		Public:  privateKey.Public(),
	}, nil
}

// Create returns a new key.
func (s *keyStorage) Create() (*key, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate encryption key")
	}

	data, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal encryption key")
	}

	id := uuid.New()

	idBytes, err := id.MarshalBinary()
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal key id")
	}

	if err := s.storage.Store(idBytes, data); err != nil {
		return nil, errors.Wrap(err, "failed to store encryption key")
	}

	return &key{
		ID:      id,
		Private: privateKey,
		Public:  privateKey.Public(),
	}, nil
}
