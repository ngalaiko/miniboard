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
	ErrUsernameEmpty   = fmt.Errorf("username is empty")
	ErrPasswordEmpty   = fmt.Errorf("password is empty")
)

// User is the user model.
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Hash     []byte `json:"-"`
}

func newUser(username string, password []byte) (*User, error) {
	if username == "" {
		return nil, ErrUsernameEmpty
	}

	if len(password) == 0 {
		return nil, ErrPasswordEmpty
	}

	hash, err := bcrypt.GenerateFromPassword(password, bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password hash: %w", err)
	}

	return &User{
		ID:       uuid.New().String(),
		Username: username,
		Hash:     hash,
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
