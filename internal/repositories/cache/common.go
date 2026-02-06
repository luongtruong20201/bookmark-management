// Package cache provides an abstraction layer for cache operations.
// It defines interfaces and implementations for storing, retrieving, and deleting
// cached data using a group-key pattern for efficient cache management.
package cache

import (
	"context"
	"time"
)

// DB defines the interface for cache database operations.
// It provides methods to set, get, and delete cached data using a two-level
// key structure: cacheGroupKey (for grouping related cache entries) and cacheKey
// (for individual cache entries within a group).
//
//go:generate mockery --name DB --filename db.go
type DB interface {
	// SetCacheData stores data in the cache with the specified group key and cache key.
	// The data will expire after the specified duration.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout
	//   - cacheGroupKey: The group key that identifies a collection of related cache entries
	//   - cacheKey: The specific key within the group for this cache entry
	//   - value: The data to be cached as a byte slice
	//   - exp: The expiration duration for the cache entry
	//
	// Returns:
	//   - error: An error if the cache operation fails
	SetCacheData(ctx context.Context, cacheGroupKey, cacheKey string, value []byte, exp time.Duration) error

	// GetCacheData retrieves cached data using the specified group key and cache key.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout
	//   - cacheGroupKey: The group key that identifies a collection of related cache entries
	//   - cacheKey: The specific key within the group for this cache entry
	//
	// Returns:
	//   - []byte: The cached data, or nil if not found
	//   - error: An error if the cache operation fails or the key doesn't exist
	GetCacheData(ctx context.Context, cacheGroupKey, cacheKey string) ([]byte, error)

	// DeleteCacheData removes all cache entries associated with the specified group key.
	// This is useful for invalidating all cached data for a particular entity or collection.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout
	//   - cacheGroupKey: The group key whose associated cache entries should be deleted
	//
	// Returns:
	//   - error: An error if the cache operation fails
	DeleteCacheData(ctx context.Context, cacheGroupKey string) error
}
