package redis

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_Queue(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	host := os.Getenv("REDIS_HOST")
	if host == "" {
		t.Skip("REDIS_HOST is not set")
	}

	t.Run("With redis", func(t *testing.T) {
		queue, err := New(ctx, host)
		assert.NoError(t, err)

		t.Run("With successful message processor", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			channel := "test"

			messages := make(chan []byte)
			go func() {
				err := queue.Listen(ctx, channel, func(data []byte) (bool, error) {
					messages <- data
					return true, nil
				})

				assert.NoError(t, err)
				close(messages)
			}()

			t.Run("When sending message to a channel", func(t *testing.T) {
				msg := []byte("test")

				time.Sleep(100 * time.Millisecond)

				assert.NoError(t, queue.Send(ctx, channel, msg))
				t.Run("It should be delivered", func(t *testing.T) {
					assert.Equal(t, msg, <-messages)
				})
			})
		})

		t.Run("When retried", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			channel := "retry"

			messages := make(chan []byte)
			go func() {
				c := 0
				err := queue.Listen(ctx, channel, func(data []byte) (bool, error) {
					if c == 1 {
						messages <- data
						return true, nil
					}
					c++
					return false, nil
				})
				assert.NoError(t, err)
				close(messages)
			}()

			t.Run("When sending message to a channel", func(t *testing.T) {
				msg := []byte("retry test")

				time.Sleep(100 * time.Millisecond)

				assert.NoError(t, queue.Send(ctx, channel, msg))
				t.Run("It should be delivered", func(t *testing.T) {
					assert.Equal(t, msg, <-messages)
				})
			})
		})

		t.Run("When error", func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()

			channel := "error"

			messages := make(chan []byte)
			go func() {
				c := 0
				err := queue.Listen(ctx, channel, func(data []byte) (bool, error) {
					if c == 1 {
						messages <- data
						return true, nil
					}
					c++
					return true, errors.New("test")
				})
				assert.NoError(t, err)
				close(messages)
			}()

			t.Run("When sending message to a channel", func(t *testing.T) {
				msg := []byte("error test")

				time.Sleep(100 * time.Millisecond)

				assert.NoError(t, queue.Send(ctx, channel, msg))

				t.Run("It should not be delivered", func(t *testing.T) {
					assert.Nil(t, <-messages)
				})
			})
		})
	})
}
