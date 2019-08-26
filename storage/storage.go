package storage // import "miniboard.app/storage"

import "miniboard.app/storage/resource"

// Storage used to store key value data.
type Storage interface {
	// Store stores data by the resource name.
	Store(*resource.Name, []byte) error
	// Load returns data by the resource name.
	Load(*resource.Name) ([]byte, error)
	// Delete deletes data by the resource name.
	Delete(*resource.Name) error
	// LoadChildren returns the list of resouce's children.
	// If _from_ is not nil, returns _limit_ documents starting with _from_ ordered by descanding.
	// If _from_ is nil, returns _limit_ first documents ordered by descanding.
	// NOTE: name.ID() doesn't matter in this case and treated as a wildcard.
	// NOTE: IDs of the resource must be sortable to get the right order.
	LoadChildren(name *resource.Name, from *resource.Name, limit int) ([][]byte, error)
}
