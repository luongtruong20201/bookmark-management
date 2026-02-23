package queue

import (
	"context"
	"errors"
)

// Queue name constants for Redis list keys.
const (
	QueueNameBookmarkImport = "bookmark_import"
)

// Repository defines the storage interface for queue messages.
// Implementations are responsible for pushing serialized messages to the
// underlying queueing system.
//
//go:generate mockery --name Repository --filename queue.go
type Repository interface {
	// PushMessage enqueues the given message bytes into the queue.
	PushMessage(ctx context.Context, msg []byte) error
	PopMessage(ctx context.Context) ([]byte, error)
}

var (
	ErrNoMessage = errors.New("no message")
)
