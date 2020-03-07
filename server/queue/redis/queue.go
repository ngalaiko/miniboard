package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"miniboard.app/queue"
)

var _ queue.Queue = &Queue{}

const healthCheckPeriod = time.Minute

// Queue is a redis backed queue.
type Queue struct {
	pool *redis.Pool
}

// New returns new queue instance.
func New(ctx context.Context, addr string) (*Queue, error) {
	pool := newPool(addr)

	conn := pool.Get()
	defer func() {
		_ = conn.Close()
	}()

	if _, err := conn.Do("PING"); err != nil {
		return nil, fmt.Errorf("failed to ping redis client: %w", err)
	}

	log("queue").Infof("connected to redis on: %s", addr)
	return &Queue{
		pool: pool,
	}, nil
}

// Send implements queue.Queue#Send.
func (q *Queue) Send(ctx context.Context, channel string, data []byte) error {
	conn := q.pool.Get()
	defer func() {
		_ = conn.Close()
	}()
	_, err := conn.Do("PUBLISH", channel, data)
	return err
}

// Receive implements queue.Queue#Receive.
func (q *Queue) Receive(ctx context.Context, channel string, onMessage func([]byte) error) error {
	conn := q.pool.Get()
	defer func() {
		_ = conn.Close()
	}()
	psc := redis.PubSubConn{Conn: conn}

	if err := psc.Subscribe(redis.Args{}.Add(channel)...); err != nil {
		return err
	}

	done := make(chan error, 1)

	go func() {
		for {
			switch n := psc.Receive().(type) {
			case error:
				done <- n
				return
			case redis.Message:
				if err := onMessage(n.Data); err != nil {
					log("queue").Errorf("error processing a message from %s: %s", channel, err)
					return
				}
			case redis.Subscription:
				switch n.Count {
				case 1:
					log("queue").Infof("subscribed to %s", channel)
				case 0:
					// Return from the goroutine when all channels are unsubscribed.
					done <- nil
					return
				}
			}
		}
	}()

	ticker := time.NewTicker(healthCheckPeriod)
	defer ticker.Stop()
	var err error
loop:
	for err == nil {
		select {
		case <-ticker.C:
			// Send ping to test health of connection and server. If
			// corresponding pong is not received, then receive on the
			// connection will timeout and the receive goroutine will exit.
			if err = psc.Ping(""); err != nil {
				break loop
			}
		case <-ctx.Done():
			break loop
		case err := <-done:
			// Return error from the receive goroutine.
			return err
		}
	}

	// Signal the receiving goroutine to exit by unsubscribing from all channels.
	if err := psc.Unsubscribe(); err != nil {
		return err
	}

	// Wait for goroutine to complete.
	return <-done
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

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
