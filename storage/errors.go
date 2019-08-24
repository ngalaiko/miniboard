package storage // import "miniboard.app/storage"

import "errors"

// Common errors.
var (
	ErrAlreadyExists = errors.New("already exists")
	ErrNotFound      = errors.New("not found")
)
