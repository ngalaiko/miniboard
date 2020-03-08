package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

// Storage used to store key value data.
type Storage struct {
	db *redis.Client
}

// New returns new mongo client instance.
func New(ctx context.Context, addr string) (*Storage, error) {
	redisdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	go func() {
		<-ctx.Done()
		log("storage").Infof("stopping client")
		_ = redisdb.Close()
	}()

	if _, err := redisdb.Ping().Result(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	log("redis").Infof("connected to %s", addr)

	return &Storage{
		db: redisdb,
	}, nil
}

// Store implements storage.Storage.
func (s *Storage) Store(ctx context.Context, name *resource.Name, data []byte) error {
	db := s.db.WithContext(ctx)

	if err := db.Set(name.String(), data, 0).Err(); err != nil {
		return fmt.Errorf("failed to SET %s: %w", name, err)
	}

	first, last := name.Split()

	// add element to to the list
	index, err := db.LPush(first, last).Result()
	if err != nil {
		return fmt.Errorf("failed to LPUSH %s %s: %w", first, last, err)
	}

	// add save element index
	if _, err := db.HSet(first+"/hash", last, index).Result(); err != nil {
		return fmt.Errorf("failed to HSET %s %s %d: %w", first, last, index, err)
	}

	return nil
}

// Update implements storage.Storage.
func (s *Storage) Update(ctx context.Context, name *resource.Name, data []byte) error {
	db := s.db.WithContext(ctx)

	if err := db.Set(name.String(), data, 0).Err(); err != nil {
		return fmt.Errorf("failed to SET %s: %w", name, err)
	}
	return nil
}

// Load implements storage.Storage.
func (s *Storage) Load(ctx context.Context, name *resource.Name) ([]byte, error) {
	return s.loadOne(ctx, name)
}

func (s *Storage) loadOne(ctx context.Context, name *resource.Name) ([]byte, error) {
	dd, err := s.loadMany(ctx, name)
	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to MGET %s: %w", name, err)
	case dd[0] == nil:
		return nil, storage.ErrNotFound
	default:
		return dd[0], nil
	}
}

func (s *Storage) loadMany(ctx context.Context, names ...*resource.Name) ([][]byte, error) {
	if len(names) == 0 {
		return nil, nil
	}

	ns := make([]string, 0, len(names))
	for _, name := range names {
		ns = append(ns, name.String())
	}

	db := s.db.WithContext(ctx)

	dd, err := db.MGet(ns...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to MGET %s: %w", names, err)
	}

	data := make([][]byte, 0, len(dd))
	for _, d := range dd {
		switch b := d.(type) {
		case nil:
			data = append(data, nil)
		case []byte:
			data = append(data, b)
		case string:
			data = append(data, []byte(b))
		default:
			panic(fmt.Sprint("unexpected type"))
		}
	}
	return data, nil
}

// Delete implements storage.Storage.
func (s *Storage) Delete(ctx context.Context, name *resource.Name) error {
	db := s.db.WithContext(ctx)

	if _, err := db.Del("DEL", name.String()).Result(); err != nil {
		return fmt.Errorf("failed to DEL %s: %w", name, err)
	}
	return nil
}

// ForEach implements storage.Storage.
func (s *Storage) ForEach(ctx context.Context, name *resource.Name, from *resource.Name, limit int64, okFunc func(*resource.Resource) (bool, error)) error {
	db := s.db.WithContext(ctx)

	var start int64
	// get start position from a list
	if from != nil {
		first, last := from.Split()

		len, err := db.LLen(first).Result()
		if err != nil {
			return fmt.Errorf("failed to LLEN %s: %w", first, err)
		}

		index, err := db.HGet(first+"/hash", last).Int64()
		switch err {
		case nil:
			start = len - index
		case redis.Nil:
		default:
			return fmt.Errorf("failed to HGET %s %s: %w", first, last, err)
		}
	}

	first, _ := name.Split()
	keys, err := s.db.LRange(first, start, -1).Result()
	if err != nil {
		return fmt.Errorf("failed: LRANGE %s %d -1: %w", first, start, err)
	}

	nn := make([]*resource.Name, 0, len(keys))
	for _, key := range keys {
		nn = append(nn, resource.ParseName(fmt.Sprintf("%s/%s", first, key)))
	}

	// TODO: load in batches

	dd, err := s.loadMany(ctx, nn...)
	if err != nil {
		return fmt.Errorf("failed to load: %w", err)
	}

	var c int64
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

		c++
		if limit != 0 && c == limit {
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
