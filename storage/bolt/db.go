package bolt // import "miniboard.app/storage/bolt"

import (
	"context"
	"os"

	bolt "github.com/coreos/bbolt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"miniboard.app/storage"
)

var _ storage.DB = &DB{}

// DB is a boltdb powered storage implementation.
type DB struct {
	db *bolt.DB
}

// New creates new storage instance. Database is storad in the _path_.
func New(ctx context.Context, path string) (*DB, error) {
	if _, err := os.Open(path); err != nil {
		log("bolt").Infof("creating storage in %s", path)
	} else {
		log("bolt").Infof("found storage in %s", path)
	}

	db, err := bolt.Open(path, 0600, &bolt.Options{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create bolt database in '%s'", path)
	}

	go func() {
		<-ctx.Done()

		log("bolt").Infof("[bolt] closing storage %s", path)
		if err := db.Close(); err != nil {
			log("bolt").Errorf("[bolt] closing storage error: %s", err)
		}
	}()

	return &DB{
		db: db,
	}, nil
}

// Namespace creates new bucket.
func (db *DB) Namespace(name string) storage.Storage {
	byteName := []byte(name)

	if err := db.db.Update(func(tx *bolt.Tx) error {
		if tx.Bucket(byteName) != nil {
			log("bolt").Infof("found bucket '%s'", name)
			return nil
		}
		_, err := tx.CreateBucket(byteName)
		log("bolt").Infof("created bucket '%s'", name)
		return err
	}); err != nil {
		log("bolt").Panicf("failed to create bucket: %s", name)
	}

	return &Bucket{
		db:   db.db,
		name: byteName,
	}
}
func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
