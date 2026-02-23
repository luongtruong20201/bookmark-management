package queue

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

// redisQueue is a Redis-backed implementation of the Repository interface.
// It uses a Redis list identified by queueName to store queued messages.
type redisQueue struct {
	redis     *redis.Client
	queueName string
}

// NewRedisQueue creates a new Redis-based implementation of the queue Repository.
// Messages are stored in a Redis list under the given queueName key.
func NewRedisQueue(redis *redis.Client, queueName string) Repository {
	return &redisQueue{
		redis:     redis,
		queueName: queueName,
	}
}

// PushMessage enqueues the given serialized message into the Redis list
// associated with this queue.
func (r *redisQueue) PushMessage(ctx context.Context, msg []byte) error {
	return r.redis.LPush(ctx, r.queueName, msg).Err()
}

// PopMessage removes and returns one message from the right side of the Redis list
// (FIFO with respect to PushMessage, which uses LPush). Returns the message bytes,
// or ErrNoMessage when the list is empty. Other Redis errors are returned as-is.
func (r *redisQueue) PopMessage(ctx context.Context) ([]byte, error) {
	val, err := r.redis.RPop(ctx, r.queueName).Bytes()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return nil, ErrNoMessage
		default:
			return nil, err
		}
	}
	return val, nil
}
