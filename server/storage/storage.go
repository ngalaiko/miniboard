package storage

import (
	"context"

	"miniboard.app/storage/resource"
)

// Storage used to store key value data.
type Storage interface {
	// Store stores data by the resource name.
	Store(context.Context, *resource.Name, []byte) error
	// Update updates data by the resource name.
	Update(context.Context, *resource.Name, []byte) error
	// Load returns data by the resource name.
	Load(context.Context, *resource.Name) ([]byte, error)
	// Delete deletes data by the resource name.
	Delete(context.Context, *resource.Name) error
	// Iterates over the resource children. Stops, if _okFunc_ returns false or an error.
	// If from provided, starts from that element.
	// NOTE: sorted by DESC.
	// NOTE: IDs of the resource must be sortable to get the correct order.
	ForEach(ctx context.Context, name *resource.Name, from *resource.Name, okFunc func(*resource.Resource) (bool, error)) error
}
