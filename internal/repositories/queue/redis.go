package queue

import (
	"context"

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
