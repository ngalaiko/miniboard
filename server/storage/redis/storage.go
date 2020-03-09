package redis

import (
	"context"
	"fmt"
	"strings"
	"time"

	redis "github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"miniboard.app/storage"
	"miniboard.app/storage/resource"
)

// Storage used to store key value data.
type Storage struct {
	db redis.UniversalClient
}

// New returns new mongo client instance.
func New(ctx context.Context, addr string) (*Storage, error) {
	redisdb := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:              []string{addr},
		MaxRetries:         10,
		MinRetryBackoff:    time.Millisecond,
		MaxRetryBackoff:    time.Second,
		PoolSize:           100,
		MaxConnAge:         10 * time.Minute,
		MinIdleConns:       10,
		IdleTimeout:        5 * time.Minute,
		IdleCheckFrequency: 10 * time.Second,
	})

	go func() {
		<-ctx.Done()
		log().Infof("stopping client")
		_ = redisdb.Close()
	}()

	if _, err := redisdb.Ping().Result(); err != nil {
		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	log().Infof("connected to %s", addr)

	return &Storage{
		db: redisdb,
	}, nil
}

// Store implements storage.Storage.
func (s *Storage) Store(ctx context.Context, name *resource.Name, data []byte) error {
	if err := s.db.Set(name.String(), data, 0).Err(); err != nil {
		return fmt.Errorf("failed to SET %s: %w", name, err)
	}

	first, last := name.Split()

	// add element to to the list
	index, err := s.db.LPush(first, last).Result()
	if err != nil {
		return fmt.Errorf("failed to LPUSH %s %s: %w", first, last, err)
	}

	// add save element index
	if _, err := s.db.HSet(first+"/hash", last, index).Result(); err != nil {
		return fmt.Errorf("failed to HSET %s %s %d: %w", first, last, index, err)
	}

	return nil
}

// Update implements storage.Storage.
func (s *Storage) Update(ctx context.Context, name *resource.Name, data []byte) error {
	if err := s.db.Set(name.String(), data, 0).Err(); err != nil {
		return fmt.Errorf("failed to SET %s: %w", name, err)
	}
	return nil
}

// Load implements storage.Storage.
func (s *Storage) Load(ctx context.Context, name *resource.Name) ([]byte, error) {
	return s.loadOne(ctx, name)
}

func (s *Storage) loadOne(ctx context.Context, name *resource.Name) ([]byte, error) {
	dd, err := s.loadMany(ctx, name.String())
	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to MGET %s: %w", name, err)
	case dd[0] == nil:
		return nil, storage.ErrNotFound
	default:
		return dd[0], nil
	}
}

func (s *Storage) loadMany(ctx context.Context, names ...string) ([][]byte, error) {
	if len(names) == 0 {
		return nil, nil
	}

	dd, err := s.db.MGet(names...).Result()
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
	if _, err := s.db.Del(name.String()).Result(); err != nil {
		return fmt.Errorf("failed to DEL %s: %w", name, err)
	}
	return nil
}

// LoadAll implements storage.Storage.
func (s *Storage) LoadAll(ctx context.Context, name *resource.Name) ([][]byte, error) {
	kk := []string{}
	iter := s.db.Scan(0, name.String(), 100).Iterator()
	for iter.Next() {
		if strings.HasSuffix(iter.Val(), "/hash") {
			continue
		}
		kk = append(kk, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate: %w", err)
	}

	dd, err := s.loadMany(ctx, kk...)
	return dd, err
}

// ForEach implements storage.Storage.
func (s *Storage) ForEach(ctx context.Context, name *resource.Name, from *resource.Name, limit int64, okFunc func(*resource.Resource) (bool, error)) error {
	var start int64
	// get start position from a list
	if from != nil {
		first, last := from.Split()

		len, err := s.db.LLen(first).Result()
		if err != nil {
			return fmt.Errorf("failed to LLEN %s: %w", first, err)
		}

		index, err := s.db.HGet(first+"/hash", last).Int64()
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
	ss := make([]string, 0, len(keys))
	for _, key := range keys {
		nn = append(nn, resource.ParseName(fmt.Sprintf("%s/%s", first, key)))
		ss = append(ss, fmt.Sprintf("%s/%s", first, key))
	}

	// TODO: load in batches

	dd, err := s.loadMany(ctx, ss...)
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

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "redis",
	})
}
