package users

import "testing"

func Test_User_newUser__invalid_username(t *testing.T) {
	_, err := newUser("", []byte("password"))
	if err != ErrUsernameEmpty {
		t.Fatalf("expected %s, got %s", ErrUsernameEmpty, err)
	}
}

func Test_User_newUser__invalid_password(t *testing.T) {
	_, err := newUser("username", []byte(""))
	if err != ErrPasswordEmpty {
		t.Fatalf("expected %s, got %s", ErrPasswordEmpty, err)
	}
}

func Test_User_ValidatePassword(t *testing.T) {
	user, err := newUser("username", []byte("password"))
	if err != nil {
		t.Fatalf("failed to create a user: %s", err)
	}

	if err := user.ValidatePassword([]byte("password")); err != nil {
		t.Fatalf("failed to validate password: %s", err)
	}
}

func Test_User_ValidatePassword__invalid_password(t *testing.T) {
	user, err := newUser("username", []byte("password"))
	if err != nil {
		t.Fatalf("failed to create a user: %s", err)
	}

	if err := user.ValidatePassword([]byte("wrong")); err != ErrInvalidPassword {
		t.Fatalf("expected %s, got %s", ErrInvalidPassword, err)
	}
}

func Test_User_ValidatePassword__invalid_hash(t *testing.T) {
	user, err := newUser("username", []byte("password"))
	if err != nil {
		t.Fatalf("failed to create a user: %s", err)
	}

	user.Hash = []byte{}
	if err := user.ValidatePassword([]byte("password")); err == nil {
		t.Fatalf("expected error")
	}
}
