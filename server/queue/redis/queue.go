package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"miniboard.app/queue"
)

var _ queue.Queue = &Queue{}

// Queue is a redis backed queue.
type Queue struct {
	name   string
	client *redis.Client
}

// New returns new queue instance.
func New(ctx context.Context, addr string) (*Queue, error) {
	redisdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := redisdb.Ping().Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to '%s': %w", addr, err)
	}

	go func() {
		<-ctx.Done()
		log().Infof("stopping client")
		_ = redisdb.Close()
	}()

	log().Infof("connected to %s", addr)

	return &Queue{
		client: redisdb,
	}, nil
}

// Send implements queue.Queue#Send.
func (q *Queue) Send(ctx context.Context, channel string, data []byte) error {
	client := q.client.WithContext(ctx)

	if err := client.Publish(channel, data).Err(); err != nil {
		return fmt.Errorf("failed to send msg to '%s': %w", q.name, err)
	}

	return nil
}

// Listen implements queue.Queue#Receive.
func (q *Queue) Listen(ctx context.Context, channel string, onMessage func([]byte) (bool, error)) error {
	client := q.client.WithContext(ctx)

	pubsub := client.Subscribe(channel)

	// Wait for confirmation that subscription is created before publishing anything.
	if _, err := pubsub.Receive(); err != nil {
		return fmt.Errorf("failed to subscribe to '%s': %w", channel, err)
	}

	log().Infof("listening to '%s'", channel)

	for {
		select {
		case <-ctx.Done():
			if err := pubsub.Close(); err != nil {
				return fmt.Errorf("failed to close channel '%s': %w", channel, err)
			}
			return nil
		case msg := <-pubsub.Channel():
			q.processMessage(ctx, msg, onMessage)
		}
	}
}

func (q *Queue) processMessage(ctx context.Context, msg *redis.Message, onMessage func(data []byte) (bool, error)) {
	defer func() {
		if r := recover(); r != nil {
			log().Panicf("failed to process message from '%s': %s", msg.Channel, r)
		}
	}()

	payload := []byte(msg.Payload)

	acked, err := onMessage(payload)
	if err != nil {
		log().Errorf("faield to process message from '%s': %s", msg.Channel, err)
		return
	}

	if !acked {
		if err := q.Send(ctx, msg.Channel, payload); err != nil {
			log().Errorf("faield to reschedule message to %s: %s", msg.Channel, err)
		}
		return
	}
}

func log() *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": "queue",
	})
}
