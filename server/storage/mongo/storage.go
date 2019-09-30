package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"miniboard.app/storage/resource"
)

const connectTimeout = time.Second

// Storage used to store key value data.
type Storage struct {
	db *mongo.Database
}

// New returns new mongo client instance.
func New(ctx context.Context, uri string) (*Storage, error) {
	ctx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create mongo client")
	}

	if err := client.Connect(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to connect to mongo")
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, errors.Wrap(err, "failed to ping mongo")
	}

	log("mongo").Infof("connected to %s", uri)
	return &Storage{
		db: client.Database("miniboard"),
	}, nil
}

// Store implements storage.Storage.
func (s *Storage) Store(*resource.Name, []byte) error {
	return nil
}

// Update implements storage.Storage.
func (s *Storage) Update(*resource.Name, []byte) error {
	return nil
}

// Load implements storage.Storage.
func (s *Storage) Load(*resource.Name) ([]byte, error) {
	return nil, nil
}

// Delete implements storage.Storage.
func (s *Storage) Delete(*resource.Name) error {
	return nil
}

// LoadChildren implements storage.Storage.
func (s *Storage) LoadChildren(name *resource.Name, from *resource.Name, limit int) ([]*resource.Resource, error) {
	return nil, nil
}

// ForEach implements storage.Storage.
func (s *Storage) ForEach(name *resource.Name, from *resource.Name, okFunc func(*resource.Resource) (bool, error)) error {
	return nil
}

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
