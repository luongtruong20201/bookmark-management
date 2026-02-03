package url

import "context"

// Get retrieves the URL associated with the given key from Redis storage.
// It returns the URL string if the key exists, or redis.Nil error if the key is not found.
// Any other error indicates a Redis operation failure.
func (s *urlStorage) Get(ctx context.Context, key string) (string, error) {
	return s.client.Get(ctx, key).Result()
}
