package passwords

import (
	"github.com/pkg/errors"
	"github.com/raja/argon2pw"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

// Service controlls user's passwords.
type Service struct {
	storage storage.Storage
}

// NewService returns new service instance
func NewService(db storage.Storage) *Service {
	return &Service{
		storage: db,
	}
}

// Set sets _user_ password to _password_.
func (s *Service) Set(userName string, password string) error {
	hash, err := argon2pw.GenerateSaltedHash(password)
	if err != nil {
		return errors.Wrap(err, "failed to calculate password hash")
	}

	if err := s.storage.Store(resource.NewName("passwords", userName), []byte(hash)); err != nil {
		return errors.Wrap(err, "failed to store password hash")
	}

	return nil
}

// Validate validates user's password.
func (s *Service) Validate(userName string, password string) (bool, error) {
	hash, err := s.storage.Load(resource.NewName("passwords", userName))
	switch err {
	case nil:
		valid, _ := argon2pw.CompareHashWithPassword(string(hash), password)
		return valid, nil
	default:
		return false, errors.Wrap(err, "failed to calculate password hash")
	}
}
