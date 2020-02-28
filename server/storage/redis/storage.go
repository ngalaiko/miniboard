package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
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
		return nil, fmt.Errorf("failed to ping redis client: %w", err)
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
		return fmt.Errorf("failed to SET %s: %w", name, err)
	}

	first, last := name.Split()

	// add element to to the list
	index, err := redis.Int(conn.Do("LPUSH", first, last))
	if err != nil {
		return fmt.Errorf("failed to LPUSH %s %s: %w", first, last, err)
	}

	// add save element index
	if _, err := conn.Do("HSET", first+"/hash", last, index); err != nil {
		return fmt.Errorf("failed to HSET %s %s %d: %w", first, last, index, err)
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
		return fmt.Errorf("failed to SET %s: %w", name, err)
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
		return nil, fmt.Errorf("failed to MGET %s: %w", name, err)
	case dd[0] == nil:
		return nil, storage.ErrNotFound
	default:
		return dd[0], nil
	}
}

func (s *Storage) loadMany(conn redis.Conn, names ...*resource.Name) ([][]byte, error) {
	if len(names) == 0 {
		return nil, nil
	}

	ns := make([]string, 0, len(names))
	for _, name := range names {
		ns = append(ns, name.String())
	}

	data, err := redis.ByteSlices(conn.Do("MGET", redis.Args{}.AddFlat(ns)...))
	if err != nil {
		return nil, fmt.Errorf("failed to MGET %s: %w", names, err)
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
		return fmt.Errorf("failed to DEL %s: %w", name, err)
	}
	return nil
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
			return fmt.Errorf("failed to LLEN %s: %w", first, err)
		}

		index, err := redis.Int(conn.Do("HGET", first+"/hash", last))
		switch err {
		case nil:
			start = len - index
		case redis.ErrNil:
		default:
			return fmt.Errorf("failed to HGET %s %s: %w", first, last, err)
		}
	}

	first, _ := name.Split()
	keys, err := redis.Strings(conn.Do("LRANGE", first, start, -1))
	if err != nil {
		return fmt.Errorf("failed: LRANGE %s %d -1: %w", first, start, err)
	}

	nn := make([]*resource.Name, 0, len(keys))
	for _, key := range keys {
		nn = append(nn, resource.ParseName(fmt.Sprintf("%s/%s", first, key)))
	}

	// TODO: load in batches

	dd, err := s.loadMany(conn, nn...)
	if err != nil {
		return fmt.Errorf("failed to load: %w", err)
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
