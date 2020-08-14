package redis

import (
	"context"
	"fmt"
	"time"

	redis "github.com/go-redis/redis/v7"
	"github.com/ngalaiko/miniboard/server/storage"
	"github.com/ngalaiko/miniboard/server/storage/resource"
	"github.com/sirupsen/logrus"
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
	pipe := s.db.TxPipeline()

	pipe.Set(name.String(), data, 0)
	collection, _ := name.Split()

	pipe.ZAdd(collection, &redis.Z{
		Member: name.String(),
	})

	if _, err := pipe.ExecContext(ctx); err != nil {
		return fmt.Errorf("failed to store %s: %w", name, err)
	}

	return nil
}

// Load implements storage.Storage.
func (s *Storage) Load(ctx context.Context, name *resource.Name) ([]byte, error) {
	return s.loadOne(ctx, name)
}

func (s *Storage) loadOne(ctx context.Context, name *resource.Name) ([]byte, error) {
	rr, err := s.loadMany(ctx, name.String())
	switch {
	case err != nil:
		return nil, fmt.Errorf("failed to MGET %s: %w", name, err)
	case rr[0].Data == nil:
		return nil, storage.ErrNotFound
	default:
		return rr[0].Data, nil
	}
}

func (s *Storage) loadMany(_ context.Context, names ...string) ([]*resource.Resource, error) {
	if len(names) == 0 {
		return nil, nil
	}

	dd, err := s.db.MGet(names...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to MGET %s: %w", names, err)
	}

	rr := make([]*resource.Resource, 0, len(names))
	for i, d := range dd {
		var data []byte
		switch b := d.(type) {
		case nil:
			data = nil
		case []byte:
			data = b
		case string:
			data = []byte(b)
		default:
			panic(fmt.Sprint("unexpected type"))
		}

		rr = append(rr, &resource.Resource{
			Name: resource.ParseName(names[i]),
			Data: data,
		})
	}

	return rr, nil
}

// Delete implements storage.Storage.
func (s *Storage) Delete(ctx context.Context, name *resource.Name) error {
	pipe := s.db.TxPipeline()

	pipe.Del(name.String())
	collection, _ := name.Split()

	pipe.ZRem(collection, name.String())

	if _, err := pipe.Exec(); err != nil {
		return fmt.Errorf("failed to delete %s: %w", name, err)
	}

	return nil
}

// LoadAll implements storage.Storage.
func (s *Storage) LoadAll(ctx context.Context, name *resource.Name) ([]*resource.Resource, error) {
	kk := []string{}
	iter := s.db.Scan(0, name.String(), 100).Iterator()
	for iter.Next() {
		kk = append(kk, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate: %w", err)
	}

	dd, err := s.loadMany(ctx, kk...)
	return dd, err
}

// ForEach implements storage.Storage.
func (s *Storage) ForEach(ctx context.Context, name *resource.Name, from *resource.Name, okFunc func(*resource.Resource) (bool, error)) error {
	collection, _ := name.Split()

	fromName := "+"
	if from != nil {
		fromName = "[" + from.String()
	}
	keys, err := s.db.ZRevRangeByLex(collection, &redis.ZRangeBy{
		Min: "-",
		Max: fromName,
	}).Result()
	if err != nil {
		return fmt.Errorf("failed: ZREVRANGEBYLEX %s %s : %w", collection, fromName, err)
	}

	return s.loadBatch(ctx, keys, 50, okFunc)
}

func (s *Storage) loadBatch(ctx context.Context, names []string, size int64, okFunc func(*resource.Resource) (bool, error)) error {
	if len(names) == 0 {
		return nil
	}

	l := int64(len(names))
	if l < size {
		size = l
	}

	rr, err := s.loadMany(ctx, names[:size]...)
	if err != nil {
		return fmt.Errorf("failed to load: %w", err)
	}

	for _, r := range rr {
		goon, err := okFunc(r)
		if err != nil {
			return err
		}
		if !goon {
			return nil
		}
	}

	return s.loadBatch(ctx, names[size:], size, okFunc)
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "redis",
	})
}
