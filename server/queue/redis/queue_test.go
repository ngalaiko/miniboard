package redis

import (
	"context"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func Test_Queue(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s, err := miniredis.Run()
	assert.NoError(t, err)
	defer s.Close()

	t.Run("With redis", func(t *testing.T) {
		queue, err := New(ctx, s.Addr())
		if err != nil {
			t.Fatalf("failed to create database: %s", err)
		}

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
