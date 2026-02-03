package url

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
//go:generate mockery --name URLStorage --filename url.go
type URLStorage interface {
	// Store(context.Context, string, string) error
	StoreIfNotExists(context.Context, string, string, int) (bool, error)
	// Get retrieves the URL associated with the given key from storage.
	// It returns the URL string if found, or redis.Nil error if the key does not exist.
	// Any other error indicates a storage operation failure.
	Get(context.Context, string) (string, error)
	// Exists(context.Context, string) (bool, error)
}

// urlStorage implements the URLStorage interface and provides Redis-based storage
// for shortened URL mappings. It uses a Redis client to store and retrieve URL data
// with configurable expiration times.
type urlStorage struct {
	client *redis.Client
}

// NewURLStorage creates a new URL storage repository instance with the provided Redis client.
func NewURLStorage(client *redis.Client) URLStorage {
	return &urlStorage{
		client: client,
	}
}
