package worker

import "context"

// Handler processes a single message from the queue. Used by the pool
// workers; implementations may decode the payload and call services (e.g.
// bookmark import).
//
//go:generate mockery --name Handler --filename handler.go --output ./mocks
type Handler interface {
	// Handle processes the message. Returning an error is logged but does
	// not stop the worker; panics are recovered and may trigger a restart.
	Handle(ctx context.Context, message []byte) error
}
