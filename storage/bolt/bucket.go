package bolt // import "miniboard.app/storage/bolt"

import (
	"strings"

	bolt "github.com/coreos/bbolt"
	"github.com/pkg/errors"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

// buckets structure:
// <resource_type>/buckets/resource/path/<resource_id>: <resource_data>

// Store stores _data_ by _in_ in the bucket.
func (db *DB) Store(name *resource.Name, data []byte) error {
	name = resource.NewName(name.Type(), "bucket").AddChild(name)
	return db.update(name, func(bucket *bolt.Bucket) error {
		return bucket.Put([]byte(name.ID()), data)
	})
}

// Load returns returns data from bucket by _id_.
func (db *DB) Load(name *resource.Name) ([]byte, error) {
	var data []byte
	name = resource.NewName(name.Type(), "bucket").AddChild(name)
	return data, db.view(name, func(bucket *bolt.Bucket) error {
		data = bucket.Get([]byte(name.ID()))
		if len(data) == 0 {
			return storage.ErrNotFound
		}
		return nil
	})
}

// LoadChildren implements storage.Storage.
func (db *DB) LoadChildren(name *resource.Name, from *resource.Name, limit int) ([][]byte, error) {
	var data [][]byte
	name = resource.NewName(name.Type(), "bucket").AddChild(name)
	return data, db.view(name, func(bucket *bolt.Bucket) error {
		c := bucket.Cursor()

		var k, v []byte
		if from == nil {
			k, v = c.First()
		} else {
			k, v = c.Seek([]byte(from.ID()))
		}

		if k == nil {
			return nil
		}

		data = make([][]byte, 0, limit)
		for ; k != nil && len(data) < limit; k, v = c.Next() {
			data = append(data, v)
		}
		return nil
	})
}

func (db *DB) update(name *resource.Name, f func(*bolt.Bucket) error) error {
	return db.db.Update(func(tx *bolt.Tx) error {
		b, err := bucket(tx, name)
		if err != nil {
			return err
		}
		return f(b)
	})
}

func (db *DB) view(name *resource.Name, f func(*bolt.Bucket) error) error {
	return db.db.View(func(tx *bolt.Tx) error {
		b, err := bucket(tx, name)
		if err != nil {
			return storage.ErrNotFound
		}
		return f(b)
	})
}

// bucket returns a bucket where resources with a given name are stored.
func bucket(tx *bolt.Tx, name *resource.Name) (*bolt.Bucket, error) {
	path := strings.Split(name.String(), "/")
	path = path[:len(path)-1]

	bucket := tx.Bucket([]byte(path[0]))
	if bucket == nil {
		var err error
		bucket, err = tx.CreateBucket([]byte(path[0]))
		if err != nil {
			return nil, errors.Wrap(err, "failed to create bucket")
		}
	}

	for _, bucketName := range path[1:] {
		childBucket := bucket.Bucket([]byte(bucketName))
		if childBucket == nil {
			var err error
			childBucket, err = bucket.CreateBucket([]byte(bucketName))
			if err != nil {
				return nil, errors.Wrap(err, "failed to create bucket")
			}
		}
		bucket = childBucket
	}

	return bucket, nil
}
