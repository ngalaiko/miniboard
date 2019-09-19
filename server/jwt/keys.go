package jwt

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"

	"github.com/pkg/errors"
	"github.com/segmentio/ksuid"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

// key is an encryption key.
type key struct {
	ID      string
	Private crypto.PrivateKey
	Public  crypto.PublicKey
}

type keyStorage struct {
	storage storage.Storage
	cache   *cache
}

func newKeyStorage(db storage.Storage) *keyStorage {
	return &keyStorage{
		storage: db,
		cache:   newCache(),
	}
}

// Get returns a key by id.
func (s *keyStorage) Get(id string) (*key, error) {
	fromCache, cached := s.cache.Get(id)
	if cached {
		return fromCache, nil
	}

	data, err := s.storage.Load(resource.NewName("jwt-key", id))
	if err != nil {
		return nil, errors.Wrap(err, "failed to load key")
	}

	privateKey, err := x509.ParseECPrivateKey(data)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse key")
	}

	k := &key{
		ID:      id,
		Private: privateKey,
		Public:  privateKey.Public(),
	}
	s.cache.Save(id, k)

	return k, nil
}

// Delete deletes a key by id.
func (s *keyStorage) Delete(id string) error {
	if err := s.storage.Delete(resource.NewName("jwt-key", id)); err != nil {
		return errors.Wrap(err, "failed to delete key")
	}
	s.cache.Delete(id)
	return nil
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

	id := ksuid.New().String()

	if err := s.storage.Store(resource.NewName("jwt-key", id), data); err != nil {
		return nil, errors.Wrap(err, "failed to store encryption key")
	}

	k := &key{
		ID:      id,
		Private: privateKey,
		Public:  privateKey.Public(),
	}
	s.cache.Save(id, k)

	return k, nil
}

// List returns all keys from the storage.
func (s *keyStorage) List() ([]*key, error) {
	dd, err := s.storage.LoadChildren(resource.NewName("jwt-key", "*"), nil, 50)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load keys")
	}

	kk := make([]*key, 0, len(dd))
	for _, d := range dd {
		privateKey, err := x509.ParseECPrivateKey(d.Data)
		if err != nil {
			log("jwt").Errorf("failed to parse key '%s': %s", d.Name, err)
			continue
		}
		kk = append(kk, &key{
			ID:      d.Name.ID(),
			Public:  privateKey.Public,
			Private: privateKey,
		})
	}

	return kk, nil
}
