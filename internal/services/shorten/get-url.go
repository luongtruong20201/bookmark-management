package shorten

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

// GetURL retrieves the original URL associated with the given short code from the repository.
// It returns the original URL if found, or ErrCodeNotFound if the code does not exist in storage.
// Any other error from the repository is returned as-is.
func (s *shortenURL) GetURL(ctx context.Context, code string) (string, error) {
	url, err := s.repository.Get(ctx, code)
	if errors.Is(err, redis.Nil) {
		return "", ErrCodeNotFound
	}

	return url, err
}
