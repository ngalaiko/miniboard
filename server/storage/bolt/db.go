package bolt

import (
	"context"
	"fmt"
	"os"

	bolt "github.com/coreos/bbolt"
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
		log().Infof("creating storage in %s", path)
	} else {
		log().Infof("found storage in %s", path)
	}

	db, err := bolt.Open(path, 0600, &bolt.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to create bolt database in '%s': %w", path, err)
	}

	go func() {
		<-ctx.Done()

		log().Infof("closing storage %s", path)
		if err := db.Close(); err != nil {
			log().Errorf("closing storage error: %s", err)
		}
	}()

	return &DB{
		db: db,
	}, nil
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "bolt",
	})
}
