package url

import (
	"context"
	"time"
)

// StoreIfNotExists stores a URL with the given code if the code does not already exist.
// It returns true if the code was successfully stored, false if it already exists, and an error if storage fails.
// The expire parameter specifies the expiration time in seconds (0 means default expiration).
func (s *urlStorage) StoreIfNotExists(ctx context.Context, code, url string, expire int) (bool, error) {
	duration := expireTime
	if expire > 0 {
		duration = time.Duration(expire) * time.Second
	}

	return s.client.SetNX(ctx, code, url, duration).Result()
}
