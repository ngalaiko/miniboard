package jwt

import (
	"sync"

	"github.com/google/uuid"
)

// in-memory cache for encryption keys.
type cache struct {
	store *sync.Map
}

func newCache() *cache {
	return &cache{
		store: &sync.Map{},
	}
}

// Save saves the _key_ by _id_
func (c *cache) Save(id uuid.UUID, k *key) {
	c.store.Store(id, k)
}

func (c *cache) Get(id uuid.UUID) (*key, bool) {
	value, ok := c.store.Load(id)
	if !ok {
		return nil, false
	}

	k, ok := value.(*key)
	if !ok {
		return nil, false
	}

	return k, true
}
