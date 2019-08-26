package jwt

import (
	"sync"
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
func (c *cache) Save(id string, k *key) {
	c.store.Store(id, k)
}

// Delete deletes the key by _id_.
func (c *cache) Delete(id string) {
	c.store.Delete(id)
}

// Get returns a _key_ by _id_.
func (c *cache) Get(id string) (*key, bool) {
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
