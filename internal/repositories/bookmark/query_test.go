package bookmark

import (
	"context"
	"testing"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/internal/test/fixture"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		userID        string
		offset        int
		limit         int
		expectedError error
		expectedCount int
		expectedIDs   []string
		verifyFunc    func(t *testing.T, bookmarks []*model.Bookmark)
	}{
		{
			name:          "success - get all bookmarks for user with multiple bookmarks",
			userID:        "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
			offset:        0,
			limit:         10,
			expectedError: nil,
			expectedCount: 2,
			expectedIDs: []string{
				"a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
				"b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
			},
			verifyFunc: func(t *testing.T, bookmarks []*model.Bookmark) {
				assert.Len(t, bookmarks, 2)
				assert.Equal(t, "Facebook - Social Media Platform", bookmarks[0].Description)
				assert.Equal(t, "Google - Search Engine", bookmarks[1].Description)
				assert.True(t, bookmarks[0].CreatedAt.Before(bookmarks[1].CreatedAt) || bookmarks[0].CreatedAt.Equal(bookmarks[1].CreatedAt))
			},
		},
		{
			name:          "success - get bookmarks with pagination limit",
			userID:        "2f6c9d14-9e42-4c31-9e6a-0f8c4b2a7c55",
			offset:        0,
			limit:         1,
			expectedError: nil,
			expectedCount: 1,
			expectedIDs: []string{
				"c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f",
			},
			verifyFunc: func(t *testing.T, bookmarks []*model.Bookmark) {
				assert.Len(t, bookmarks, 1)
				assert.Equal(t, "GitHub - Code Repository", bookmarks[0].Description)
			},
		},
		{
			name:          "success - get bookmarks with offset",
			userID:        "2f6c9d14-9e42-4c31-9e6a-0f8c4b2a7c55",
			offset:        1,
			limit:         10,
			expectedError: nil,
			expectedCount: 1,
			expectedIDs: []string{
				"d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a",
			},
			verifyFunc: func(t *testing.T, bookmarks []*model.Bookmark) {
				assert.Len(t, bookmarks, 1)
				assert.Equal(t, "YouTube - Video Platform", bookmarks[0].Description)
			},
		},
		{
			name:          "success - get bookmarks for user with no bookmarks",
			userID:        "00000000-0000-0000-0000-000000000000",
			offset:        0,
			limit:         10,
			expectedError: nil,
			expectedCount: 0,
			expectedIDs:   []string{},
			verifyFunc: func(t *testing.T, bookmarks []*model.Bookmark) {
				assert.Empty(t, bookmarks)
			},
		},
		{
			name:          "success - get bookmarks with offset beyond available data",
			userID:        "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
			offset:        10,
			limit:         10,
			expectedError: nil,
			expectedCount: 0,
			expectedIDs:   []string{},
			verifyFunc: func(t *testing.T, bookmarks []*model.Bookmark) {
				assert.Empty(t, bookmarks)
			},
		},
		{
			name:          "success - get bookmarks for another user",
			userID:        "550e8400-e29b-41d4-a716-446655440000",
			offset:        0,
			limit:         10,
			expectedError: nil,
			expectedCount: 2,
			expectedIDs: []string{
				"c1d2e3f4-a5b6-4c7d-8e9f-0a1b2c3d4e5f",
				"d2e3f4a5-b6c7-4d8e-9f0a-1b2c3d4e5f6a",
			},
			verifyFunc: func(t *testing.T, bookmarks []*model.Bookmark) {
				assert.Len(t, bookmarks, 2)

				for _, bookmark := range bookmarks {
					assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", bookmark.UserID)
				}
			},
		},
		{
			name:          "success - get bookmarks with zero limit",
			userID:        "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
			offset:        0,
			limit:         0,
			expectedError: nil,
			expectedCount: 0,
			expectedIDs:   []string{},
			verifyFunc: func(t *testing.T, bookmarks []*model.Bookmark) {
				assert.Empty(t, bookmarks)
			},
		},
		{
			name:          "success - verify user isolation (user should only see their own bookmarks)",
			userID:        "e3c2a8f1-1d3b-4c62-8e54-6b7f9a2d1c90",
			offset:        0,
			limit:         10,
			expectedError: nil,
			expectedCount: 2,
			expectedIDs: []string{
				"e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b",
				"f6a7b8c9-d0e1-4f2a-3b4c-5d6e7f8a9b0c",
			},
			verifyFunc: func(t *testing.T, bookmarks []*model.Bookmark) {
				assert.Len(t, bookmarks, 2)

				for _, bookmark := range bookmarks {
					assert.Equal(t, "e3c2a8f1-1d3b-4c62-8e54-6b7f9a2d1c90", bookmark.UserID)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			repo := NewBookmark(db)

			bookmarks, err := repo.GetBookmarks(ctx, tc.userID, tc.offset, tc.limit)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, bookmarks)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, bookmarks)
				assert.Len(t, bookmarks, tc.expectedCount)

				if len(tc.expectedIDs) > 0 {
					actualIDs := make([]string, len(bookmarks))
					for i, b := range bookmarks {
						actualIDs[i] = b.ID
					}
					assert.Equal(t, tc.expectedIDs, actualIDs)
				}

				for _, bookmark := range bookmarks {
					assert.Equal(t, tc.userID, bookmark.UserID)
				}

				if tc.verifyFunc != nil {
					tc.verifyFunc(t, bookmarks)
				}
			}
		})
	}
}
