package queue

import "context"

// Queue is an interface for queue.
type Queue interface {
	// Send sends _data_ the queue.
	Send(ctx context.Context, data []byte) error
	// Listen execute _onMessage_ for messages from the queue.
	Receive(ctx context.Context, onMessage func([]byte) (bool, error)) error
}
