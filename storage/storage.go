package storage // import "miniboard.app/storage"

// DB contains namespaces.
type DB interface {
	Namespace(string) Storage
}

// Storage used to store key value data.
type Storage interface {
	// Store stores data by the id.
	Store(id []byte, data []byte) error
	// Load returns data by the id.
	Load(id []byte) ([]byte, error)
	// LoadPrefix returns list of data by the prefix.
	LoadPrefix(prefix []byte) ([][]byte, error)
}
