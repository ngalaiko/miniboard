package bolt // import "miniboard.app/storage/bolt"

import (
	"bytes"

	bolt "github.com/coreos/bbolt"
	"miniboard.app/application/storage"
)

var _ storage.Storage = &Bucket{}

// Bucket is using a boltdb bucket to store key value pairs.
type Bucket struct {
	db   *bolt.DB
	name []byte
}

// Store stores _data_ by _in_ in the bucket.
func (b *Bucket) Store(id []byte, data []byte) error {
	return b.update(func(bucket *bolt.Bucket) error {
		return bucket.Put(id, data)
	})
}

// Load returns returns data from bucket by _id_.
func (b *Bucket) Load(id []byte) ([]byte, error) {
	var data []byte
	return data, b.view(func(bucket *bolt.Bucket) error {
		data = bucket.Get([]byte(id))
		if len(data) == 0 {
			return storage.ErrNotFound
		}
		return nil
	})
}

// LoadPrefix returns returns data from bucket by _prefix_.
func (b *Bucket) LoadPrefix(prefix []byte) ([][]byte, error) {
	var data [][]byte
	return data, b.view(func(bucket *bolt.Bucket) error {
		c := bucket.Cursor()
		for k, v := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, v = c.Next() {
			data = append(data, v)
		}
		return nil
	})
}

func (b *Bucket) update(f func(*bolt.Bucket) error) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.name)
		return f(bucket)
	})
}

func (b *Bucket) view(f func(*bolt.Bucket) error) error {
	return b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(b.name)
		return f(bucket)
	})
}
