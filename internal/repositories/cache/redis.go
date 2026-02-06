// Package cache provides Redis-based implementation of the cache interface.
package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisCache is a Redis-based implementation of the DB interface.
// It uses Redis hash operations (HSET/HGET) to store cache entries within
// cache groups, allowing efficient bulk invalidation by deleting the entire group.
type redisCache struct {
	c *redis.Client
}

// NewRedisCache creates a new Redis cache implementation with the provided Redis client.
//
// Parameters:
//   - c: A configured Redis client instance
//
// Returns:
//   - DB: An implementation of the cache DB interface backed by Redis
func NewRedisCache(c *redis.Client) DB {
	return &redisCache{
		c: c,
	}
}

// SetCacheData stores data in Redis using a hash structure.
// The cacheGroupKey is used as the hash key, and cacheKey is the field within that hash.
// The entire hash (cache group) is set to expire after the specified duration.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - cacheGroupKey: The Redis hash key that groups related cache entries
//   - cacheKey: The field name within the hash for this specific cache entry
//   - value: The data to be cached as a byte slice
//   - exp: The expiration duration for the entire cache group
//
// Returns:
//   - error: An error if the Redis operation fails
func (db *redisCache) SetCacheData(ctx context.Context, cacheGroupKey, cacheKey string, value []byte, exp time.Duration) error {
	if err := db.c.HSet(ctx, cacheGroupKey, cacheKey, value).Err(); err != nil {
		return err
	}

	return db.c.Expire(ctx, cacheGroupKey, exp).Err()
}

// GetCacheData retrieves cached data from Redis using hash get operation.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - cacheGroupKey: The Redis hash key that groups related cache entries
//   - cacheKey: The field name within the hash for this specific cache entry
//
// Returns:
//   - []byte: The cached data as a byte slice, or nil if not found
//   - error: An error if the Redis operation fails or the key doesn't exist
func (db *redisCache) GetCacheData(ctx context.Context, cacheGroupKey, cacheKey string) ([]byte, error) {
	val, err := db.c.HGet(ctx, cacheGroupKey, cacheKey).Bytes()
	if err != nil {
		return nil, err
	}

	return val, nil
}

// DeleteCacheData removes the entire cache group (Redis hash) from Redis.
// This effectively invalidates all cache entries within that group.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - cacheGroupKey: The Redis hash key whose associated cache entries should be deleted
//
// Returns:
//   - error: An error if the Redis operation fails
func (db *redisCache) DeleteCacheData(ctx context.Context, cacheGroupKey string) error {
	return db.c.Del(ctx, cacheGroupKey).Err()
}
