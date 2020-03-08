package queue

import "context"

// Queue is an interface for queue.
type Queue interface {
	// Send sends _data_ the _channel_.
	Send(ctx context.Context, channel string, data []byte) error
	// Listen executes _onMessage_ for messages from the _channel_.
	Listen(ctx context.Context, channel string, onMessage func([]byte) (bool, error)) error
}
