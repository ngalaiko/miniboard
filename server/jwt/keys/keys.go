package keys

import (
	"fmt"

	"github.com/google/uuid"
)

// Known errors
var (
	ErrKeyIsEmpty = fmt.Errorf("key is empty")
)

// Key is a public key.
type Key struct {
	ID        string
	PublicDER []byte
}

func newKey(publicDER []byte) (*Key, error) {
	if len(publicDER) == 0 {
		return nil, ErrKeyIsEmpty
	}

	return &Key{
		ID:        uuid.New().String(),
		PublicDER: publicDER,
	}, nil
}
