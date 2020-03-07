package redis

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Queue(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("With redis", func(t *testing.T) {
		host := os.Getenv("REDIS_HOST")
		if host == "" {
			t.Skip("no redis host provided")
		}

		queue, err := New(ctx, host)
		if err != nil {
			t.Fatalf("failed to create database: %s", err)
		}

		conn := queue.pool.Get()
		_, err = conn.Do("FLUSHALL")
		_ = conn.Close()
		assert.NoError(t, err)

		messages := make(chan []byte)

		testChan := "test_channel"
		go func() {
			err := queue.Receive(ctx, testChan, func(data []byte) error {
				messages <- data
				return nil
			})
			assert.NoError(t, err)
		}()

		t.Run("When sending message to a channel", func(t *testing.T) {
			msg := []byte("test")
			err := queue.Send(ctx, testChan, msg)
			assert.NoError(t, err)
			t.Run("it should not be found", func(t *testing.T) {
				assert.Equal(t, msg, <-messages)
			})
		})
	})
}
