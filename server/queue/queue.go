package queue

import "context"

// Queue is an interface for queue.
type Queue interface {
	// Send sends _data_ to the _channel_.
	Send(ctx context.Context, channel string, data []byte) error
	// Listen execute _onMessage_ for a message from _channel_.
	Receive(ctx context.Context, channel string, onMessage func([]byte) error) error
}
