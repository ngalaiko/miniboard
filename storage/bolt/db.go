package bolt // import "miniboard.app/storage/bolt"

import (
	"context"

	bolt "github.com/coreos/bbolt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"miniboard.app/application/storage"
)

var _ storage.DB = &DB{}

// DB is a boltdb powered storage implementation.
type DB struct {
	db *bolt.DB
}

// New creates new storage instance. Database is storad in the _path_.
func New(ctx context.Context, path string) (*DB, error) {
	logrus.Infof("creating bolt storage in %s", path)

	db, err := bolt.Open(path, 0600, &bolt.Options{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create bolt database in '%s'", path)
	}

	go func() {
		<-ctx.Done()

		logrus.Infof("closing bolt storage %s", path)
		if err := db.Close(); err != nil {
			logrus.Errorf("closing bolt storage error: %s", err)
		}
	}()

	return &DB{
		db: db,
	}, nil
}

// Namespace creates new bucket.
func (db *DB) Namespace(name string) storage.Storage {
	logrus.Infof("creating bolt bucket '%s'", name)
	byteName := []byte(name)

	if err := db.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(byteName)
		return err
	}); err != nil {
		logrus.Panicf("failed to create bucket: %s", name)
	}

	return &Bucket{
		db:   db.db,
		name: byteName,
	}
}
