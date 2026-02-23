package worker

import "context"

// Queue is the source of messages for the engine. Implementations typically
// read from a Redis list or other message broker.
//
//go:generate mockery --name Queue --filename queue.go --output ./mocks
type Queue interface {
	// PopMessage returns the next message, or an error when empty or on failure.
	PopMessage(ctx context.Context) ([]byte, error)
}
