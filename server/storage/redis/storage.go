package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

const poolSize = 10

// Storage used to store key value data.
type Storage struct {
	db *redis.Pool
}

// New returns new mongo client instance.
func New(ctx context.Context, addr string) (*Storage, error) {
	pool := newPool(addr)

	conn := pool.Get()
	defer func() {
		_ = conn.Close()
	}()

	if _, err := conn.Do("PING"); err != nil {
		return nil, errors.Wrap(err, "failed to ping redis client")
	}

	log("redis").Infof("connected to %s", addr)
	return &Storage{
		db: pool,
	}, nil
}

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		// Dial or DialContext must be set. When both are set, DialContext takes precedence over Dial.
		Dial: func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

// Store implements storage.Storage.
func (s *Storage) Store(ctx context.Context, name *resource.Name, data []byte) error {
	conn := s.db.Get()
	defer func() {
		_ = conn.Close()
	}()

	if _, err := conn.Do("SET", name.String(), data); err != nil {
		return errors.Wrapf(err, "failed to SET %s", name)
	}
	first, last := name.Split()
	if _, err := conn.Do("LPUSH", first, last); err != nil {
		return errors.Wrapf(err, "failed to LPUSH %s %s", first, last)
	}
	return nil
}

// Update implements storage.Storage.
func (s *Storage) Update(ctx context.Context, name *resource.Name, data []byte) error {
	conn := s.db.Get()
	defer func() {
		_ = conn.Close()
	}()

	if _, err := conn.Do("SET", name.String(), data); err != nil {
		return errors.Wrapf(err, "failed to SET %s", name)
	}
	return nil
}

// Load implements storage.Storage.
func (s *Storage) Load(ctx context.Context, name *resource.Name) ([]byte, error) {
	conn := s.db.Get()
	defer func() {
		_ = conn.Close()
	}()
	return s.load(ctx, conn, name)
}

func (s *Storage) load(ctx context.Context, conn redis.Conn, name *resource.Name) ([]byte, error) {
	data, err := redis.Bytes(conn.Do("GET", name.String()))
	switch err {
	case nil:
		return data, nil
	case redis.ErrNil:
		return nil, storage.ErrNotFound
	default:
		return nil, errors.Wrapf(err, "failed to GET %s", name)
	}
}

// Delete implements storage.Storage.
func (s *Storage) Delete(ctx context.Context, name *resource.Name) error {
	conn := s.db.Get()
	defer func() {
		_ = conn.Close()
	}()

	if _, err := conn.Do("DEL", name.String()); err != nil {
		return errors.Wrapf(err, "failed to DEL %s", name)
	}
	first, last := name.Split()
	if _, err := conn.Do("LREM", first, 0, last); err != nil {
		return errors.Wrapf(err, "failed to LREM %s 0 %s", first, last)
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
	conn := s.db.Get()
	defer func() {
		_ = conn.Close()
	}()

	first, _ := name.Split()

	var fromIndex int64
	if from != nil {
		_, id := from.Split()

		var err error
		fromIndex, err = redis.Int64(conn.Do("LINDEX", first, id))
		if err != nil {
			return errors.Wrapf(err, "invalid from value %s", from.String())
		}
	}

	keys, err := redis.Strings(conn.Do("LRANGE", first, fromIndex, -1))
	if err != nil {
		return errors.Wrapf(err, "failed: LRANGE %s %d -1", first, fromIndex)
	}

	for _, key := range keys {
		res := &resource.Resource{
			Name: resource.ParseName(fmt.Sprintf("%s/%s", first, key)),
		}

		res.Data, err = s.load(ctx, conn, res.Name)
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