package users

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 14

// Known errors.
var (
	ErrInvalidPassword = fmt.Errorf("invalid password")
)

// User is the user model.
type User struct {
	ID   string `json:"id"`
	Hash []byte `json:"hash"`
}

func newUser(password []byte) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword(password, bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password hash: %w", err)
	}

	return &User{
		ID:   uuid.New().String(),
		Hash: hash,
	}, nil
}

// ValidatePassword returns nil if a given password is valid.
func (u *User) ValidatePassword(password []byte) error {
	err := bcrypt.CompareHashAndPassword(u.Hash, password)
	switch err {
	case nil:
		return nil
	case bcrypt.ErrMismatchedHashAndPassword:
		return ErrInvalidPassword
	default:
		return fmt.Errorf("failed to validate password: %w", err)
	}
}
