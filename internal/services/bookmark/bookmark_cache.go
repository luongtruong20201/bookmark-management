package bookmark

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/internal/repositories/cache"
	"github.com/rs/zerolog/log"
)

// bookmarkCache is a caching decorator around the bookmark Service.
//
// It caches the result of `GetBookmarks` per-user and per pagination tuple
// (offset, limit) to reduce database load. Write operations (`Create`, `Update`,
// `Delete`) invalidate the per-user cache group to keep reads consistent.
//
// Cache layout:
// - group key: `get_bookmarks_<userID>`
// - item key:  `<offset>_<limit>`
//
// Notes:
// - Cache failures are non-fatal: on cache miss/unmarshal error it falls back to
//   the underlying service; on cache set/delete errors it logs and continues.
type bookmarkCache struct {
	Service
	cache cache.DB
}

const (
	// getBookmarksCacheGroupFormat is the cache "group" key used to partition cached
	// bookmark list responses by user.
	getBookmarksCacheGroupFormat = "get_bookmarks_%s"
	// getBookmarksCacheKeyFormat is the cache "item" key used within a user group,
	// parameterized by offset and limit.
	getBookmarksCacheKeyFormat   = "%d_%d"
	// getBookmarksCacheDuration is the TTL for cached bookmark list responses.
	getBookmarksCacheDuration    = time.Hour
)

// NewBookmarkCache creates a new bookmark cache service that wraps the provided
// bookmark service with caching functionality. It uses the cache to store and
// retrieve bookmark data, reducing database queries.
func NewBookmarkCache(s Service, cache cache.DB) *bookmarkCache {
	return &bookmarkCache{
		Service: s,
		cache:   cache,
	}
}

// getCacheGroupKey generates a cache group key for a user's bookmarks
func (c *bookmarkCache) getCacheGroupKey(userID string) string {
	return fmt.Sprintf(getBookmarksCacheGroupFormat, userID)
}

// invalidateUserCache deletes the cache for a user's bookmarks.
// Errors are logged but do not prevent the operation from continuing.
func (c *bookmarkCache) invalidateUserCache(ctx context.Context, userID string) {
	cacheGroupKey := c.getCacheGroupKey(userID)
	if err := c.cache.DeleteCacheData(ctx, cacheGroupKey); err != nil {
		log.Warn().Err(err).Str("userID", userID).Msg("failed to invalidate user cache")
	}
}

// GetBookmarks retrieves bookmarks for a user with pagination support.
// It first attempts to retrieve data from cache using the cache group key
// and pagination-specific cache key. If cache miss occurs or unmarshal fails,
// it fetches from the underlying service and caches the result.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - userID: The unique identifier of the user whose bookmarks to retrieve
//   - offset: The number of records to skip (for pagination)
//   - limit: The maximum number of records to return
//
// Returns:
//   - *GetBookmarksResponse: A response containing bookmarks and total count
//   - error: An error if the service operation fails
func (s *bookmarkCache) GetBookmarks(ctx context.Context, userID string, offset, limit int) (*GetBookmarksResponse, error) {
	cacheGroupKey := s.getCacheGroupKey(userID)
	cacheKey := fmt.Sprintf(getBookmarksCacheKeyFormat, offset, limit)

	cacheData, err := s.cache.GetCacheData(ctx, cacheGroupKey, cacheKey)
	if err == nil && len(cacheData) > 0 {
		result := &GetBookmarksResponse{}
		if err := json.Unmarshal(cacheData, result); err == nil {
			return result, nil
		}
		log.Warn().Err(err).Msg("failed to unmarshal cached data, fetching from service")
	}

	result, err := s.Service.GetBookmarks(ctx, userID, offset, limit)
	if err != nil {
		return result, err
	}

	resultBytes, err := json.Marshal(result)
	if err == nil {
		if err := s.cache.SetCacheData(ctx, cacheGroupKey, cacheKey, resultBytes, getBookmarksCacheDuration); err != nil {
			log.Error().Err(err).Msg("failed to cache data")
		}
	} else {
		log.Error().Err(err).Msg("failed to marshal result for caching")
	}

	return result, nil
}

// Create creates a new bookmark for a user. It invalidates the user's bookmark
// cache before delegating to the underlying service to ensure cache consistency.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - description: The description of the bookmark
//   - url: The URL of the bookmark
//   - userId: The unique identifier of the user creating the bookmark
//
// Returns:
//   - *model.Bookmark: The created bookmark
//   - error: An error if the creation fails
func (c *bookmarkCache) Create(ctx context.Context, description, url, userId string) (*model.Bookmark, error) {
	c.invalidateUserCache(ctx, userId)
	return c.Service.Create(ctx, description, url, userId)
}

// Update updates an existing bookmark. It invalidates the user's bookmark cache
// before delegating to the underlying service to ensure cache consistency.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - bookmarkID: The unique identifier of the bookmark to update
//   - userID: The unique identifier of the user who owns the bookmark
//   - description: The new description of the bookmark
//   - url: The new URL of the bookmark
//
// Returns:
//   - *model.Bookmark: The updated bookmark
//   - error: An error if the update fails
func (c *bookmarkCache) Update(ctx context.Context, bookmarkID, userID, description, url string) (*model.Bookmark, error) {
	c.invalidateUserCache(ctx, userID)
	return c.Service.Update(ctx, bookmarkID, userID, description, url)
}

// Delete deletes a bookmark. It invalidates the user's bookmark cache before
// delegating to the underlying service to ensure cache consistency.
//
// Parameters:
//   - ctx: Context for request cancellation and timeout
//   - bookmarkID: The unique identifier of the bookmark to delete
//   - userID: The unique identifier of the user who owns the bookmark
//
// Returns:
//   - error: An error if the deletion fails
func (c *bookmarkCache) Delete(ctx context.Context, bookmarkID, userID string) error {
	c.invalidateUserCache(ctx, userID)
	return c.Service.Delete(ctx, bookmarkID, userID)
}
