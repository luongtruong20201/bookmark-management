package bookmark_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	models "github.com/luongtruong20201/bookmark-management/internal/models"
	bookmark "github.com/luongtruong20201/bookmark-management/internal/services/bookmark"
	cacheMocks "github.com/luongtruong20201/bookmark-management/internal/repositories/cache/mocks"
	serviceMocks "github.com/luongtruong20201/bookmark-management/internal/services/bookmark/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewBookmarkCache(t *testing.T) {
	t.Parallel()

	service := serviceMocks.NewService(t)
	cache := cacheMocks.NewDB(t)

	cacheService := bookmark.NewBookmarkCache(service, cache)

	assert.NotNil(t, cacheService)
}

func TestBookmarkCache_GetBookmarks(t *testing.T) {
	t.Parallel()

	var (
		testErrService = errors.New("service error")
		testErrCache   = errors.New("cache error")
		mockUserID     = "550e8400-e29b-41d4-a716-446655440000"
	)

	testCases := []struct {
		name           string
		userID         string
		offset         int
		limit          int
		setupService   func(t *testing.T, ctx context.Context, userID string, offset, limit int) *serviceMocks.Service
		setupCache     func(t *testing.T, ctx context.Context, cacheGroupKey, cacheKey string, shouldHit bool, shouldUnmarshalFail bool, shouldSetFail bool, shouldMarshalFail bool) *cacheMocks.DB
		expectedError  error
		expectedCount  int
		expectedTotal  int64
		verifyResponse func(t *testing.T, resp *bookmark.GetBookmarksResponse)
	}{
		{
			name:   "success - cache hit with valid data",
			userID: mockUserID,
			offset: 0,
			limit:  10,
			setupCache: func(t *testing.T, ctx context.Context, cacheGroupKey, cacheKey string, shouldHit bool, shouldUnmarshalFail bool, shouldSetFail bool, shouldMarshalFail bool) *cacheMocks.DB {
				cache := cacheMocks.NewDB(t)
				response := &bookmark.GetBookmarksResponse{
					Data: []*models.Bookmark{
						{
							Base: models.Base{
								ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
							},
							Description: "Facebook",
							URL:         "https://www.facebook.com",
							Code:        "abc1234",
							UserID:      mockUserID,
						},
					},
					Total: 1,
				}
				data, _ := json.Marshal(response)
				cache.On("GetCacheData", ctx, cacheGroupKey, cacheKey).Return(data, nil).Once()
				return cache
			},
			setupService: func(t *testing.T, ctx context.Context, userID string, offset, limit int) *serviceMocks.Service {
				return serviceMocks.NewService(t)
			},
			expectedError: nil,
			expectedCount: 1,
			expectedTotal: 1,
			verifyResponse: func(t *testing.T, resp *bookmark.GetBookmarksResponse) {
				assert.NotNil(t, resp)
				assert.Len(t, resp.Data, 1)
				assert.Equal(t, int64(1), resp.Total)
				assert.Equal(t, "Facebook", resp.Data[0].Description)
			},
		},
		{
			name:   "success - cache miss, fetch from service and cache",
			userID: mockUserID,
			offset: 0,
			limit:  10,
			setupCache: func(t *testing.T, ctx context.Context, cacheGroupKey, cacheKey string, shouldHit bool, shouldUnmarshalFail bool, shouldSetFail bool, shouldMarshalFail bool) *cacheMocks.DB {
				cache := cacheMocks.NewDB(t)
				cache.On("GetCacheData", ctx, cacheGroupKey, cacheKey).Return([]byte(nil), testErrCache).Once()
				response := &bookmark.GetBookmarksResponse{
					Data: []*models.Bookmark{
						{
							Base: models.Base{
								ID: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
							},
							Description: "Google",
							URL:         "https://www.google.com",
							Code:        "def5678",
							UserID:      mockUserID,
						},
					},
					Total: 1,
				}
				data, _ := json.Marshal(response)
				cache.On("SetCacheData", ctx, cacheGroupKey, cacheKey, data, time.Hour).Return(nil).Once()
				return cache
			},
			setupService: func(t *testing.T, ctx context.Context, userID string, offset, limit int) *serviceMocks.Service {
				service := serviceMocks.NewService(t)
				response := &bookmark.GetBookmarksResponse{
					Data: []*models.Bookmark{
						{
							Base: models.Base{
								ID: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
							},
							Description: "Google",
							URL:         "https://www.google.com",
							Code:        "def5678",
							UserID:      userID,
						},
					},
					Total: 1,
				}
				service.On("GetBookmarks", ctx, userID, offset, limit).Return(response, nil).Once()
				return service
			},
			expectedError: nil,
			expectedCount: 1,
			expectedTotal: 1,
			verifyResponse: func(t *testing.T, resp *bookmark.GetBookmarksResponse) {
				assert.NotNil(t, resp)
				assert.Len(t, resp.Data, 1)
				assert.Equal(t, int64(1), resp.Total)
				assert.Equal(t, "Google", resp.Data[0].Description)
			},
		},
		{
			name:   "success - cache hit with empty data, fetch from service",
			userID: mockUserID,
			offset: 0,
			limit:  10,
			setupCache: func(t *testing.T, ctx context.Context, cacheGroupKey, cacheKey string, shouldHit bool, shouldUnmarshalFail bool, shouldSetFail bool, shouldMarshalFail bool) *cacheMocks.DB {
				cache := cacheMocks.NewDB(t)
				cache.On("GetCacheData", ctx, cacheGroupKey, cacheKey).Return([]byte{}, nil).Once()
				response := &bookmark.GetBookmarksResponse{
					Data: []*models.Bookmark{
						{
							Base: models.Base{
								ID: "c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f",
							},
							Description: "GitHub",
							URL:         "https://www.github.com",
							Code:        "ghi9012",
							UserID:      mockUserID,
						},
					},
					Total: 1,
				}
				data, _ := json.Marshal(response)
				cache.On("SetCacheData", ctx, cacheGroupKey, cacheKey, data, time.Hour).Return(nil).Once()
				return cache
			},
			setupService: func(t *testing.T, ctx context.Context, userID string, offset, limit int) *serviceMocks.Service {
				service := serviceMocks.NewService(t)
				response := &bookmark.GetBookmarksResponse{
					Data: []*models.Bookmark{
						{
							Base: models.Base{
								ID: "c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f",
							},
							Description: "GitHub",
							URL:         "https://www.github.com",
							Code:        "ghi9012",
							UserID:      userID,
						},
					},
					Total: 1,
				}
				service.On("GetBookmarks", ctx, userID, offset, limit).Return(response, nil).Once()
				return service
			},
			expectedError: nil,
			expectedCount: 1,
			expectedTotal: 1,
			verifyResponse: func(t *testing.T, resp *bookmark.GetBookmarksResponse) {
				assert.NotNil(t, resp)
				assert.Len(t, resp.Data, 1)
				assert.Equal(t, "GitHub", resp.Data[0].Description)
			},
		},
		{
			name:   "success - cache hit but unmarshal fails, fetch from service",
			userID: mockUserID,
			offset: 0,
			limit:  10,
			setupCache: func(t *testing.T, ctx context.Context, cacheGroupKey, cacheKey string, shouldHit bool, shouldUnmarshalFail bool, shouldSetFail bool, shouldMarshalFail bool) *cacheMocks.DB {
				cache := cacheMocks.NewDB(t)
				cache.On("GetCacheData", ctx, cacheGroupKey, cacheKey).Return([]byte("invalid json"), nil).Once()
				response := &bookmark.GetBookmarksResponse{
					Data: []*models.Bookmark{
						{
							Base: models.Base{
								ID: "d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a",
							},
							Description: "YouTube",
							URL:         "https://www.youtube.com",
							Code:        "jkl3456",
							UserID:      mockUserID,
						},
					},
					Total: 1,
				}
				data, _ := json.Marshal(response)
				cache.On("SetCacheData", ctx, cacheGroupKey, cacheKey, data, time.Hour).Return(nil).Once()
				return cache
			},
			setupService: func(t *testing.T, ctx context.Context, userID string, offset, limit int) *serviceMocks.Service {
				service := serviceMocks.NewService(t)
				response := &bookmark.GetBookmarksResponse{
					Data: []*models.Bookmark{
						{
							Base: models.Base{
								ID: "d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a",
							},
							Description: "YouTube",
							URL:         "https://www.youtube.com",
							Code:        "jkl3456",
							UserID:      userID,
						},
					},
					Total: 1,
				}
				service.On("GetBookmarks", ctx, userID, offset, limit).Return(response, nil).Once()
				return service
			},
			expectedError: nil,
			expectedCount: 1,
			expectedTotal: 1,
			verifyResponse: func(t *testing.T, resp *bookmark.GetBookmarksResponse) {
				assert.NotNil(t, resp)
				assert.Len(t, resp.Data, 1)
				assert.Equal(t, "YouTube", resp.Data[0].Description)
			},
		},
		{
			name:   "error - service returns error",
			userID: mockUserID,
			offset: 0,
			limit:  10,
			setupCache: func(t *testing.T, ctx context.Context, cacheGroupKey, cacheKey string, shouldHit bool, shouldUnmarshalFail bool, shouldSetFail bool, shouldMarshalFail bool) *cacheMocks.DB {
				cache := cacheMocks.NewDB(t)
				cache.On("GetCacheData", ctx, cacheGroupKey, cacheKey).Return([]byte(nil), testErrCache).Once()
				return cache
			},
			setupService: func(t *testing.T, ctx context.Context, userID string, offset, limit int) *serviceMocks.Service {
				service := serviceMocks.NewService(t)
				service.On("GetBookmarks", ctx, userID, offset, limit).Return(nil, testErrService).Once()
				return service
			},
			expectedError:  testErrService,
			expectedCount:  0,
			expectedTotal:  0,
			verifyResponse: nil,
		},
		{
			name:   "success - cache set fails but still return result",
			userID: mockUserID,
			offset: 0,
			limit:  10,
			setupCache: func(t *testing.T, ctx context.Context, cacheGroupKey, cacheKey string, shouldHit bool, shouldUnmarshalFail bool, shouldSetFail bool, shouldMarshalFail bool) *cacheMocks.DB {
				cache := cacheMocks.NewDB(t)
				cache.On("GetCacheData", ctx, cacheGroupKey, cacheKey).Return([]byte(nil), testErrCache).Once()
				response := &bookmark.GetBookmarksResponse{
					Data: []*models.Bookmark{
						{
							Base: models.Base{
								ID: "e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b",
							},
							Description: "Twitter",
							URL:         "https://www.twitter.com",
							Code:        "mno7890",
							UserID:      mockUserID,
						},
					},
					Total: 1,
				}
				data, _ := json.Marshal(response)
				cache.On("SetCacheData", ctx, cacheGroupKey, cacheKey, data, time.Hour).Return(testErrCache).Once()
				return cache
			},
			setupService: func(t *testing.T, ctx context.Context, userID string, offset, limit int) *serviceMocks.Service {
				service := serviceMocks.NewService(t)
				response := &bookmark.GetBookmarksResponse{
					Data: []*models.Bookmark{
						{
							Base: models.Base{
								ID: "e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b",
							},
							Description: "Twitter",
							URL:         "https://www.twitter.com",
							Code:        "mno7890",
							UserID:      userID,
						},
					},
					Total: 1,
				}
				service.On("GetBookmarks", ctx, userID, offset, limit).Return(response, nil).Once()
				return service
			},
			expectedError: nil,
			expectedCount: 1,
			expectedTotal: 1,
			verifyResponse: func(t *testing.T, resp *bookmark.GetBookmarksResponse) {
				assert.NotNil(t, resp)
				assert.Len(t, resp.Data, 1)
				assert.Equal(t, "Twitter", resp.Data[0].Description)
			},
		},
		{
			name:   "success - with different offset and limit",
			userID: mockUserID,
			offset: 10,
			limit:  5,
			setupCache: func(t *testing.T, ctx context.Context, cacheGroupKey, cacheKey string, shouldHit bool, shouldUnmarshalFail bool, shouldSetFail bool, shouldMarshalFail bool) *cacheMocks.DB {
				cache := cacheMocks.NewDB(t)
				cache.On("GetCacheData", ctx, cacheGroupKey, cacheKey).Return([]byte(nil), testErrCache).Once()
				response := &bookmark.GetBookmarksResponse{
					Data:  []*models.Bookmark{},
					Total: 15,
				}
				data, _ := json.Marshal(response)
				cache.On("SetCacheData", ctx, cacheGroupKey, cacheKey, data, time.Hour).Return(nil).Once()
				return cache
			},
			setupService: func(t *testing.T, ctx context.Context, userID string, offset, limit int) *serviceMocks.Service {
				service := serviceMocks.NewService(t)
				response := &bookmark.GetBookmarksResponse{
					Data:  []*models.Bookmark{},
					Total: 15,
				}
				service.On("GetBookmarks", ctx, userID, offset, limit).Return(response, nil).Once()
				return service
			},
			expectedError: nil,
			expectedCount: 0,
			expectedTotal: 15,
			verifyResponse: func(t *testing.T, resp *bookmark.GetBookmarksResponse) {
				assert.NotNil(t, resp)
				assert.Empty(t, resp.Data)
				assert.Equal(t, int64(15), resp.Total)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			cacheGroupKey := fmt.Sprintf("get_bookmarks_%s", tc.userID)
			cacheKey := fmt.Sprintf("%d_%d", tc.offset, tc.limit)

			service := tc.setupService(t, ctx, tc.userID, tc.offset, tc.limit)
			cache := tc.setupCache(t, ctx, cacheGroupKey, cacheKey, false, false, false, false)

			cacheService := bookmark.NewBookmarkCache(service, cache)

			result, err := cacheService.GetBookmarks(ctx, tc.userID, tc.offset, tc.limit)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Data, tc.expectedCount)
				assert.Equal(t, tc.expectedTotal, result.Total)
				if tc.verifyResponse != nil {
					tc.verifyResponse(t, result)
				}
			}
		})
	}
}

func TestBookmarkCache_Create(t *testing.T) {
	t.Parallel()

	var (
		testErrService = errors.New("service error")
		testErrCache   = errors.New("cache error")
		mockUserID     = "550e8400-e29b-41d4-a716-446655440000"
	)

	testCases := []struct {
		name           string
		description    string
		url            string
		userID         string
		setupService   func(t *testing.T, ctx context.Context, description, url, userID string) *serviceMocks.Service
		setupCache     func(t *testing.T, ctx context.Context, userID string, shouldDeleteFail bool) *cacheMocks.DB
		expectedError  error
		verifyBookmark func(t *testing.T, bookmark *models.Bookmark)
	}{
		{
			name:        "success - create bookmark and invalidate cache",
			description: "My blog",
			url:         "https://truonglq.com",
			userID:      mockUserID,
			setupCache: func(t *testing.T, ctx context.Context, userID string, shouldDeleteFail bool) *cacheMocks.DB {
				cache := cacheMocks.NewDB(t)
				cacheGroupKey := fmt.Sprintf("get_bookmarks_%s", userID)
				if shouldDeleteFail {
					cache.On("DeleteCacheData", ctx, cacheGroupKey).Return(testErrCache).Once()
				} else {
					cache.On("DeleteCacheData", ctx, cacheGroupKey).Return(nil).Once()
				}
				return cache
			},
			setupService: func(t *testing.T, ctx context.Context, description, url, userID string) *serviceMocks.Service {
				service := serviceMocks.NewService(t)
				bookmark := &models.Bookmark{
					Base: models.Base{
						ID: "11111111-2222-3333-4444-555555555555",
					},
					Description: description,
					URL:         url,
					Code:        "abcd1234",
					UserID:      userID,
				}
				service.On("Create", ctx, description, url, userID).Return(bookmark, nil).Once()
				return service
			},
			expectedError: nil,
			verifyBookmark: func(t *testing.T, bookmark *models.Bookmark) {
				assert.Equal(t, "11111111-2222-3333-4444-555555555555", bookmark.ID)
				assert.Equal(t, "My blog", bookmark.Description)
				assert.Equal(t, "https://truonglq.com", bookmark.URL)
				assert.Equal(t, mockUserID, bookmark.UserID)
			},
		},
		{
			name:        "success - cache invalidation fails but still create bookmark",
			description: "My blog",
			url:         "https://truonglq.com",
			userID:      mockUserID,
			setupCache: func(t *testing.T, ctx context.Context, userID string, shouldDeleteFail bool) *cacheMocks.DB {
				cache := cacheMocks.NewDB(t)
				cacheGroupKey := fmt.Sprintf("get_bookmarks_%s", userID)
				cache.On("DeleteCacheData", ctx, cacheGroupKey).Return(testErrCache).Once()
				return cache
			},
			setupService: func(t *testing.T, ctx context.Context, description, url, userID string) *serviceMocks.Service {
				service := serviceMocks.NewService(t)
				bookmark := &models.Bookmark{
					Base: models.Base{
						ID: "22222222-3333-4444-5555-666666666666",
					},
					Description: description,
					URL:         url,
					Code:        "xyz98765",
					UserID:      userID,
				}
				service.On("Create", ctx, description, url, userID).Return(bookmark, nil).Once()
				return service
			},
			expectedError: nil,
			verifyBookmark: func(t *testing.T, bookmark *models.Bookmark) {
				assert.Equal(t, "22222222-3333-4444-5555-666666666666", bookmark.ID)
				assert.Equal(t, "My blog", bookmark.Description)
			},
		},
		{
			name:        "error - service returns error",
			description: "My blog",
			url:         "https://truonglq.com",
			userID:      mockUserID,
			setupCache: func(t *testing.T, ctx context.Context, userID string, shouldDeleteFail bool) *cacheMocks.DB {
				cache := cacheMocks.NewDB(t)
				cacheGroupKey := fmt.Sprintf("get_bookmarks_%s", userID)
				cache.On("DeleteCacheData", ctx, cacheGroupKey).Return(nil).Once()
				return cache
			},
			setupService: func(t *testing.T, ctx context.Context, description, url, userID string) *serviceMocks.Service {
				service := serviceMocks.NewService(t)
				service.On("Create", ctx, description, url, userID).Return(nil, testErrService).Once()
				return service
			},
			expectedError:  testErrService,
			verifyBookmark: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			service := tc.setupService(t, ctx, tc.description, tc.url, tc.userID)
			cache := tc.setupCache(t, ctx, tc.userID, false)

			cacheService := bookmark.NewBookmarkCache(service, cache)

			result, err := cacheService.Create(ctx, tc.description, tc.url, tc.userID)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tc.verifyBookmark != nil {
					tc.verifyBookmark(t, result)
				}
			}
		})
	}
}

func TestBookmarkCache_Update(t *testing.T) {
	t.Parallel()

	var (
		testErrService = errors.New("service error")
		testErrCache   = errors.New("cache error")
		mockUserID     = "550e8400-e29b-41d4-a716-446655440000"
		mockBookmarkID = "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"
	)

	testCases := []struct {
		name           string
		bookmarkID     string
		userID         string
		description    string
		url            string
		setupService   func(t *testing.T, ctx context.Context, bookmarkID, userID, description, url string) *serviceMocks.Service
		setupCache     func(t *testing.T, ctx context.Context, userID string, shouldDeleteFail bool) *cacheMocks.DB
		expectedError  error
		verifyBookmark func(t *testing.T, bookmark *models.Bookmark)
	}{
		{
			name:        "success - update bookmark and invalidate cache",
			bookmarkID:  mockBookmarkID,
			userID:      mockUserID,
			description: "Updated Facebook",
			url:         "https://www.facebook.com/updated",
			setupCache: func(t *testing.T, ctx context.Context, userID string, shouldDeleteFail bool) *cacheMocks.DB {
				cache := cacheMocks.NewDB(t)
				cacheGroupKey := fmt.Sprintf("get_bookmarks_%s", userID)
				if shouldDeleteFail {
					cache.On("DeleteCacheData", ctx, cacheGroupKey).Return(testErrCache).Once()
				} else {
					cache.On("DeleteCacheData", ctx, cacheGroupKey).Return(nil).Once()
				}
				return cache
			},
			setupService: func(t *testing.T, ctx context.Context, bookmarkID, userID, description, url string) *serviceMocks.Service {
				service := serviceMocks.NewService(t)
				bookmark := &models.Bookmark{
					Base: models.Base{
						ID: bookmarkID,
					},
					Description: description,
					URL:         url,
					Code:        "abc1234",
					UserID:      userID,
				}
				service.On("Update", ctx, bookmarkID, userID, description, url).Return(bookmark, nil).Once()
				return service
			},
			expectedError: nil,
			verifyBookmark: func(t *testing.T, bookmark *models.Bookmark) {
				assert.Equal(t, mockBookmarkID, bookmark.ID)
				assert.Equal(t, "Updated Facebook", bookmark.Description)
				assert.Equal(t, "https://www.facebook.com/updated", bookmark.URL)
				assert.Equal(t, mockUserID, bookmark.UserID)
			},
		},
		{
			name:        "success - cache invalidation fails but still update bookmark",
			bookmarkID:  mockBookmarkID,
			userID:      mockUserID,
			description: "Updated Google",
			url:         "https://www.google.com/updated",
			setupCache: func(t *testing.T, ctx context.Context, userID string, shouldDeleteFail bool) *cacheMocks.DB {
				cache := cacheMocks.NewDB(t)
				cacheGroupKey := fmt.Sprintf("get_bookmarks_%s", userID)
				cache.On("DeleteCacheData", ctx, cacheGroupKey).Return(testErrCache).Once()
				return cache
			},
			setupService: func(t *testing.T, ctx context.Context, bookmarkID, userID, description, url string) *serviceMocks.Service {
				service := serviceMocks.NewService(t)
				bookmark := &models.Bookmark{
					Base: models.Base{
						ID: bookmarkID,
					},
					Description: description,
					URL:         url,
					Code:        "def5678",
					UserID:      userID,
				}
				service.On("Update", ctx, bookmarkID, userID, description, url).Return(bookmark, nil).Once()
				return service
			},
			expectedError: nil,
			verifyBookmark: func(t *testing.T, bookmark *models.Bookmark) {
				assert.Equal(t, mockBookmarkID, bookmark.ID)
				assert.Equal(t, "Updated Google", bookmark.Description)
			},
		},
		{
			name:        "error - service returns error",
			bookmarkID:  mockBookmarkID,
			userID:      mockUserID,
			description: "Updated Facebook",
			url:         "https://www.facebook.com/updated",
			setupCache: func(t *testing.T, ctx context.Context, userID string, shouldDeleteFail bool) *cacheMocks.DB {
				cache := cacheMocks.NewDB(t)
				cacheGroupKey := fmt.Sprintf("get_bookmarks_%s", userID)
				cache.On("DeleteCacheData", ctx, cacheGroupKey).Return(nil).Once()
				return cache
			},
			setupService: func(t *testing.T, ctx context.Context, bookmarkID, userID, description, url string) *serviceMocks.Service {
				service := serviceMocks.NewService(t)
				service.On("Update", ctx, bookmarkID, userID, description, url).Return(nil, testErrService).Once()
				return service
			},
			expectedError:  testErrService,
			verifyBookmark: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			service := tc.setupService(t, ctx, tc.bookmarkID, tc.userID, tc.description, tc.url)
			cache := tc.setupCache(t, ctx, tc.userID, false)

			cacheService := bookmark.NewBookmarkCache(service, cache)

			result, err := cacheService.Update(ctx, tc.bookmarkID, tc.userID, tc.description, tc.url)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tc.verifyBookmark != nil {
					tc.verifyBookmark(t, result)
				}
			}
		})
	}
}

func TestBookmarkCache_Delete(t *testing.T) {
	t.Parallel()

	var (
		testErrService = errors.New("service error")
		testErrCache   = errors.New("cache error")
		mockUserID     = "550e8400-e29b-41d4-a716-446655440000"
		mockBookmarkID = "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"
	)

	testCases := []struct {
		name          string
		bookmarkID    string
		userID        string
		setupService  func(t *testing.T, ctx context.Context, bookmarkID, userID string) *serviceMocks.Service
		setupCache    func(t *testing.T, ctx context.Context, userID string, shouldDeleteFail bool) *cacheMocks.DB
		expectedError error
	}{
		{
			name:       "success - delete bookmark and invalidate cache",
			bookmarkID: mockBookmarkID,
			userID:     mockUserID,
			setupCache: func(t *testing.T, ctx context.Context, userID string, shouldDeleteFail bool) *cacheMocks.DB {
				cache := cacheMocks.NewDB(t)
				cacheGroupKey := fmt.Sprintf("get_bookmarks_%s", userID)
				if shouldDeleteFail {
					cache.On("DeleteCacheData", ctx, cacheGroupKey).Return(testErrCache).Once()
				} else {
					cache.On("DeleteCacheData", ctx, cacheGroupKey).Return(nil).Once()
				}
				return cache
			},
			setupService: func(t *testing.T, ctx context.Context, bookmarkID, userID string) *serviceMocks.Service {
				service := serviceMocks.NewService(t)
				service.On("Delete", ctx, bookmarkID, userID).Return(nil).Once()
				return service
			},
			expectedError: nil,
		},
		{
			name:       "success - cache invalidation fails but still delete bookmark",
			bookmarkID: mockBookmarkID,
			userID:     mockUserID,
			setupCache: func(t *testing.T, ctx context.Context, userID string, shouldDeleteFail bool) *cacheMocks.DB {
				cache := cacheMocks.NewDB(t)
				cacheGroupKey := fmt.Sprintf("get_bookmarks_%s", userID)
				cache.On("DeleteCacheData", ctx, cacheGroupKey).Return(testErrCache).Once()
				return cache
			},
			setupService: func(t *testing.T, ctx context.Context, bookmarkID, userID string) *serviceMocks.Service {
				service := serviceMocks.NewService(t)
				service.On("Delete", ctx, bookmarkID, userID).Return(nil).Once()
				return service
			},
			expectedError: nil,
		},
		{
			name:       "error - service returns error",
			bookmarkID: mockBookmarkID,
			userID:     mockUserID,
			setupCache: func(t *testing.T, ctx context.Context, userID string, shouldDeleteFail bool) *cacheMocks.DB {
				cache := cacheMocks.NewDB(t)
				cacheGroupKey := fmt.Sprintf("get_bookmarks_%s", userID)
				cache.On("DeleteCacheData", ctx, cacheGroupKey).Return(nil).Once()
				return cache
			},
			setupService: func(t *testing.T, ctx context.Context, bookmarkID, userID string) *serviceMocks.Service {
				service := serviceMocks.NewService(t)
				service.On("Delete", ctx, bookmarkID, userID).Return(testErrService).Once()
				return service
			},
			expectedError: testErrService,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			service := tc.setupService(t, ctx, tc.bookmarkID, tc.userID)
			cache := tc.setupCache(t, ctx, tc.userID, false)

			cacheService := bookmark.NewBookmarkCache(service, cache)

			err := cacheService.Delete(ctx, tc.bookmarkID, tc.userID)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
