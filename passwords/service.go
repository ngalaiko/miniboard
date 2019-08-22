package passwords // import "miniboard.app/passwords"

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"miniboard.app/storage"
)

const bcryptCost = 10

// Service controlls user's passwords.
type Service struct {
	storage storage.Storage
}

// NewService returns new service instance
func NewService(db storage.DB) *Service {
	return &Service{
		storage: db.Namespace("passwords"),
	}
}

// Set sets _user_ password to _password_.
func (s *Service) Set(userName string, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return errors.Wrap(err, "failed to calculate password hash")
	}

	if err := s.storage.Store([]byte(userName), hash); err != nil {
		return errors.Wrap(err, "failed to store password hash")
	}

	return nil
}

// Validate validates user's password.
func (s *Service) Validate(userName string, password string) (bool, error) {
	hash, err := s.storage.Load([]byte(userName))
	switch err {
	case nil:
		return bcrypt.CompareHashAndPassword(hash, []byte(password)) == nil, nil
	default:
		return false, errors.Wrap(err, "failed to calculate password hash")
	}
}
