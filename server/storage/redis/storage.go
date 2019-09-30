package redis

import (
	"context"
	"sort"

	"github.com/mediocregopher/radix/v3"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

const poolSize = 10

// Storage used to store key value data.
type Storage struct {
	db *radix.Pool
}

// New returns new mongo client instance.
func New(ctx context.Context, uri string) (*Storage, error) {
	pool, err := radix.NewPool("tcp", uri, poolSize)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to redis")
	}

	if err := pool.Do(radix.Cmd(nil, "PING")); err != nil {
		return nil, errors.Wrap(err, "failed to ping redis client")
	}

	log("redis").Infof("connected to %s", uri)
	return &Storage{
		db: pool,
	}, nil
}

// Store implements storage.Storage.
func (s *Storage) Store(ctx context.Context, name *resource.Name, data []byte) error {
	if err := s.db.Do(radix.FlatCmd(nil, "SET", name.String(), data)); err != nil {
		return errors.Wrap(err, "failed to SET")
	}
	return nil
}

// Update implements storage.Storage.
func (s *Storage) Update(ctx context.Context, name *resource.Name, data []byte) error {
	if err := s.db.Do(radix.FlatCmd(nil, "SET", name.String(), data)); err != nil {
		return errors.Wrap(err, "failed to SET")
	}
	return nil
}

// Load implements storage.Storage.
func (s *Storage) Load(ctx context.Context, name *resource.Name) ([]byte, error) {
	var data []byte
	if err := s.db.Do(radix.FlatCmd(&data, "GET", name.String())); err != nil {
		return nil, errors.Wrap(err, "failed to GET")
	}
	if data == nil {
		return nil, storage.ErrNotFound
	}
	return data, nil
}

// Delete implements storage.Storage.
func (s *Storage) Delete(ctx context.Context, name *resource.Name) error {
	if err := s.db.Do(radix.FlatCmd(nil, "DEL", name.String())); err != nil {
		return errors.Wrap(err, "failed to DEL")
	}
	return nil
}

// LoadChildren implements storage.Storage.
func (s *Storage) LoadChildren(ctx context.Context, name *resource.Name, from *resource.Name, limit int) ([]*resource.Resource, error) {
	data := make([]*resource.Resource, 0, limit)
	return data, s.ForEach(ctx, name, from, func(r *resource.Resource) (bool, error) {
		if len(data) == limit {
			return false, nil
		}
		data = append(data, r)
		return true, nil
	})
}

// ForEach implements storage.Storage.
func (s *Storage) ForEach(ctx context.Context, name *resource.Name, from *resource.Name, okFunc func(*resource.Resource) (bool, error)) error {
	scanner := radix.NewScanner(s.db, radix.ScanOpts{
		Command: "SCAN",
		Pattern: name.String(),
	})
	defer func() {
		if err := scanner.Close(); err != nil {
			log("redis").Errorf("failed to close scanner: %s", err)
		}
	}()

	var key string
	var err error
	keys := []string{}
	for scanner.Next(&key) {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})

	start := false
	var fromString string
	if from == nil {
		start = true
	} else {
		fromString = from.String()
	}

	for _, key := range keys {
		if fromString == key {
			start = true
		}
		if !start {
			continue
		}

		res := &resource.Resource{
			Name: resource.ParseName(key),
		}

		res.Data, err = s.Load(ctx, res.Name)
		if err != nil {
			return errors.Wrap(err, "failed to GET")
		}

		goon, err := okFunc(res)
		if err != nil {
			return err
		}
		if !goon {
			return nil
		}
	}
	return nil
}

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
