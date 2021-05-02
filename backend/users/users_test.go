package users

import (
	"errors"
	"testing"
)

func Test_User_newUser__invalid_username(t *testing.T) {
	t.Parallel()

	_, err := newUser("", []byte("password"), 10)
	if !errors.Is(err, ErrUsernameEmpty) {
		t.Fatalf("expected %s, got %s", ErrUsernameEmpty, err)
	}
}

func Test_User_newUser__invalid_password(t *testing.T) {
	t.Parallel()

	_, err := newUser("username", []byte(""), 10)
	if !errors.Is(err, ErrPasswordEmpty) {
		t.Fatalf("expected %s, got %s", ErrPasswordEmpty, err)
	}
}

func Test_User_ValidatePassword(t *testing.T) {
	t.Parallel()

	user, err := newUser("username", []byte("password"), 10)
	if err != nil {
		t.Fatalf("failed to create a user: %s", err)
	}

	if err := user.ValidatePassword([]byte("password")); err != nil {
		t.Fatalf("failed to validate password: %s", err)
	}
}

func Test_User_ValidatePassword__invalid_password(t *testing.T) {
	t.Parallel()

	user, err := newUser("username", []byte("password"), 10)
	if err != nil {
		t.Fatalf("failed to create a user: %s", err)
	}

	if err := user.ValidatePassword([]byte("wrong")); !errors.Is(err, ErrInvalidPassword) {
		t.Fatalf("expected %s, got %s", ErrInvalidPassword, err)
	}
}

func Test_User_ValidatePassword__invalid_hash(t *testing.T) {
	t.Parallel()

	user, err := newUser("username", []byte("password"), 10)
	if err != nil {
		t.Fatalf("failed to create a user: %s", err)
	}

	user.Hash = []byte{}
	if err := user.ValidatePassword([]byte("password")); err == nil {
		t.Fatalf("expected error")
	}
}
