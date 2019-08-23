package storage // import "miniboard.app/storage"

import "miniboard.app/storage/resource"

// Storage used to store key value data.
type Storage interface {
	// Store stores data by the resource name.
	Store(*resource.Name, []byte) error
	// Load returns data by the resource name.
	Load(*resource.Name) ([]byte, error)
}
