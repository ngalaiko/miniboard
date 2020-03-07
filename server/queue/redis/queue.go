package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/adjust/rmq"
	"github.com/sirupsen/logrus"
	"miniboard.app/queue"
)

var _ queue.Queue = &Queue{}

// Queue is a redis backed queue.
type Queue struct {
	name  string
	queue rmq.Queue

	PrefetchLimit int
	PollDuration  time.Duration
}

// New returns new queue instance.
func New(ctx context.Context, addr string, name string) *Queue {
	log("queue").Infof("opened queue '%s' connection to '%s'", name, addr)
	queue := rmq.OpenConnection("producer", "tcp", addr, 1).OpenQueue(name)
	queue.SetPushQueue(queue)
	go func() {
		<-ctx.Done()
		queue.Close()
	}()
	return &Queue{
		name:          name,
		queue:         queue,
		PrefetchLimit: 50,
		PollDuration:  time.Second,
	}
}

// Send implements queue.Queue#Send.
func (q *Queue) Send(ctx context.Context, data []byte) error {
	if q.queue.PublishBytes(data) {
		return nil
	}
	return fmt.Errorf("failed to send msg to '%s'", q.name)
}

// Receive implements queue.Queue#Receive.
func (q *Queue) Receive(ctx context.Context, onMessage func([]byte) (bool, error)) error {
	if !q.queue.StartConsuming(q.PrefetchLimit, q.PollDuration) {
		return fmt.Errorf("failed to start consuming '%s'", q.name)
	}

	q.queue.AddConsumerFunc("consumer", func(delivery rmq.Delivery) {
		acked, err := onMessage([]byte(delivery.Payload()))
		switch {
		case err != nil:
			log("queue").Infof("failed to process message from '%s': %s", q.name, err)
			delivery.Reject()
		case !acked:
			delivery.Push()
		default:
			delivery.Ack()
		}
	})

	<-ctx.Done()
	<-q.queue.StopConsuming()

	return nil
}

func log(src string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"source": src,
	})
}
