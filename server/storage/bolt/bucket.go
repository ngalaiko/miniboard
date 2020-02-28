package bolt

import (
	"context"
	"fmt"
	"strings"

	bolt "github.com/coreos/bbolt"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

// buckets structure:
// <resource_type>/buckets/resource/path/<resource_id>: <resource_data>

// Store stores _data_ by _in_ in the bucket.
// Returns ErrAlreadyExists if the key _in_ already exists.
func (db *DB) Store(ctx context.Context, name *resource.Name, data []byte) error {
	name = resource.NewName(name.Type(), "bucket").AddChild(name)
	return db.update(name, func(bucket *bolt.Bucket) error {
		if bucket.Get([]byte(name.ID())) != nil {
			return storage.ErrAlreadyExists
		}
		return bucket.Put([]byte(name.ID()), data)
	})
}

// Update stores _data_ by _in_ in the bucket.
func (db *DB) Update(ctx context.Context, name *resource.Name, data []byte) error {
	name = resource.NewName(name.Type(), "bucket").AddChild(name)
	return db.update(name, func(bucket *bolt.Bucket) error {
		return bucket.Put([]byte(name.ID()), data)
	})
}

// Load returns returns data from bucket by _id_.
func (db *DB) Load(ctx context.Context, name *resource.Name) ([]byte, error) {
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

// Delete deletes data by _name_.
func (db *DB) Delete(ctx context.Context, name *resource.Name) error {
	name = resource.NewName(name.Type(), "bucket").AddChild(name)
	return db.update(name, func(bucket *bolt.Bucket) error {
		return bucket.Delete([]byte(name.ID()))
	})
}

// ForEach implements storage.Storage.
func (db *DB) ForEach(ctx context.Context, name *resource.Name, from *resource.Name, filter func(*resource.Resource) (bool, error)) error {
	return db.view(resource.NewName(name.Type(), "bucket").AddChild(name), func(bucket *bolt.Bucket) error {
		c := bucket.Cursor()

		var k, v []byte
		if from == nil {
			k, v = c.Last()
		} else {
			k, v = c.Seek([]byte(from.ID()))
		}

		if k == nil {
			return nil
		}

		for ; k != nil; k, v = c.Prev() {
			goon, err := filter(&resource.Resource{
				Name: name.Parent().Child(name.Type(), string(k)),
				Data: v,
			})
			if err != nil {
				return err
			}
			if !goon {
				return nil
			}
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
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	for _, bucketName := range path[1:] {
		childBucket := bucket.Bucket([]byte(bucketName))
		if childBucket == nil {
			var err error
			childBucket, err = bucket.CreateBucket([]byte(bucketName))
			if err != nil {
				return nil, fmt.Errorf("failed to create bucket: %w", err)
			}
		}
		bucket = childBucket
	}

	return bucket, nil
}
