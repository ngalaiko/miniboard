package bolt

import (
	"context"
	"os"

	bolt "github.com/coreos/bbolt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"miniboard.app/storage"
)

var _ storage.Storage = &DB{}

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

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
