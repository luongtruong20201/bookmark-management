package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// expireTime is the default expiration time for stored URLs.
	expireTime = 24 * time.Hour
)

// URLStorage defines the interface for URL storage repositories.
// It provides methods to store and retrieve shortened URLs.
//
//go:generate mockery --name URLStorage --filename url_storage.go
type URLStorage interface {
	// Store(context.Context, string, string) error
	StoreIfNotExists(context.Context, string, string, int) (bool, error)
	// Get(context.Context, string) (string, error)
	// Exists(context.Context, string) (bool, error)
}

type urlStorage struct {
	client *redis.Client
}

// NewURLStorage creates a new URL storage repository instance with the provided Redis client.
func NewURLStorage(client *redis.Client) URLStorage {
	return &urlStorage{
		client: client,
	}
}

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
