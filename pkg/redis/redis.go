package redis

import "github.com/redis/go-redis/v9"

// NewClient creates a new Redis client instance using configuration from environment variables.
// The prefix parameter is currently unused but reserved for future use.
// It returns a Redis client configured with address, password, and database settings.
func NewClient(prefix string) (*redis.Client, error) {
	cfg, err := newConfig("")
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	return client, nil
}
