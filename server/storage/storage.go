package storage

import (
	"context"

	"miniboard.app/storage/resource"
)

// Storage used to store key value data.
type Storage interface {
	// Store stores data by the resource name.
	Store(context.Context, *resource.Name, []byte) error
	// Load returns data by the resource name.
	Load(context.Context, *resource.Name) ([]byte, error)
	// LoadAll returns all data by the resource name.
	// NOTE: order is not guaranteed.
	LoadAll(context.Context, *resource.Name) ([][]byte, error)
	// Delete deletes data by the resource name.
	Delete(context.Context, *resource.Name) error
	// Iterates over the resource children. Stops, if _okFunc_ returns false or an error.
	// If from provided, starts from that element.
	// NOTE: sorted lexicographically by keys DESC.
	ForEach(ctx context.Context, name *resource.Name, from *resource.Name, limit int64, okFunc func(*resource.Resource) (bool, error)) error
}
