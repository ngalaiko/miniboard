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

	// add element data
	if _, err := conn.Do("SET", name.String(), data); err != nil {
		return errors.Wrapf(err, "failed to SET %s", name)
	}

	first, last := name.Split()

	// add element to to the list
	index, err := redis.Int(conn.Do("LPUSH", first, last))
	if err != nil {
		return errors.Wrapf(err, "failed to LPUSH %s %s", first, last)
	}

	// add save element index
	if _, err := conn.Do("HSET", first+"/hash", last, index); err != nil {
		return errors.Wrapf(err, "failed to HSET %s %s %d", first, last, index)
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
	return s.loadOne(conn, name)
}

func (s *Storage) loadOne(conn redis.Conn, name *resource.Name) ([]byte, error) {
	dd, err := s.loadMany(conn, name)
	switch {
	case err != nil:
		return nil, errors.Wrapf(err, "failed to MGET %s", name)
	case dd[0] == nil:
		return nil, storage.ErrNotFound
	default:
		return dd[0], nil
	}
}

func (s *Storage) loadMany(conn redis.Conn, names ...*resource.Name) ([][]byte, error) {
	ns := make([]string, 0, len(names))
	for _, name := range names {
		ns = append(ns, name.String())
	}

	data, err := redis.ByteSlices(conn.Do("MGET", redis.Args{}.AddFlat(ns)...))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to MGET %s", names)
	}

	return data, nil
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

	start := 0
	// get start position from a list
	if from != nil {
		first, last := from.Split()

		len, err := redis.Int(conn.Do("LLEN", first))
		if err != nil {
			return errors.Wrapf(err, "failed to LLEN %s", first)
		}

		index, err := redis.Int(conn.Do("HGET", first+"/hash", last))
		switch err {
		case nil:
			start = len - index
		case redis.ErrNil:
		default:
			return errors.Wrapf(err, "failed to HGET %s %s", first, last)
		}
	}

	first, _ := name.Split()
	keys, err := redis.Strings(conn.Do("LRANGE", first, start, -1))
	if err != nil {
		return errors.Wrapf(err, "failed: LRANGE %s %d -1", first, start)
	}

	nn := make([]*resource.Name, 0, len(keys))
	for _, key := range keys {
		nn = append(nn, resource.ParseName(fmt.Sprintf("%s/%s", first, key)))
	}

	// TODO: load in batches

	dd, err := s.loadMany(conn, nn...)
	if err != nil {
		return errors.Wrapf(err, "failed to load")
	}

	for i, d := range dd {
		if d == nil {
			continue
		}

		res := &resource.Resource{
			Name: nn[i],
			Data: d,
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
