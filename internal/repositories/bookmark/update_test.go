package bookmark

import (
	"context"
	"testing"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/internal/test/fixture"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
)

func TestRepository_UpdateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		bookmarkID     string
		userID         string
		updates        *model.Bookmark
		expectedError  error
		verifyFunc     func(t *testing.T, bookmark *model.Bookmark)
	}{
		{
			name:       "success - update bookmark with description and URL",
			bookmarkID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
			userID:     "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
			updates: &model.Bookmark{
				Description: "Updated Facebook Description",
				URL:         "https://www.facebook.com/updated",
			},
			expectedError: nil,
			verifyFunc: func(t *testing.T, bookmark *model.Bookmark) {
				assert.Equal(t, "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d", bookmark.ID)
				assert.Equal(t, "Updated Facebook Description", bookmark.Description)
				assert.Equal(t, "https://www.facebook.com/updated", bookmark.URL)
				assert.Equal(t, "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91", bookmark.UserID)
				assert.Equal(t, "abc12345", bookmark.Code)
			},
		},
		{
			name:       "success - update bookmark with only description",
			bookmarkID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
			userID:     "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
			updates: &model.Bookmark{
				Description: "Updated Description Only",
			},
			expectedError: nil,
			verifyFunc: func(t *testing.T, bookmark *model.Bookmark) {
				assert.Equal(t, "Updated Description Only", bookmark.Description)
				assert.Equal(t, "https://www.facebook.com", bookmark.URL)
			},
		},
		{
			name:       "success - update bookmark with only URL",
			bookmarkID: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
			userID:     "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
			updates: &model.Bookmark{
				URL: "https://www.google.com/updated",
			},
			expectedError: nil,
			verifyFunc: func(t *testing.T, bookmark *model.Bookmark) {
				assert.Equal(t, "Google - Search Engine", bookmark.Description)
				assert.Equal(t, "https://www.google.com/updated", bookmark.URL)
			},
		},
		{
			name:       "error - bookmark not found",
			bookmarkID: "00000000-0000-0000-0000-000000000000",
			userID:     "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
			updates: &model.Bookmark{
				Description: "Updated Description",
			},
			expectedError: dbutils.ErrNotFoundType,
			verifyFunc:    nil,
		},
		{
			name:       "error - bookmark belongs to different user",
			bookmarkID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
			userID:     "550e8400-e29b-41d4-a716-446655440000",
			updates: &model.Bookmark{
				Description: "Updated Description",
			},
			expectedError: dbutils.ErrNotFoundType,
			verifyFunc:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			repo := NewBookmark(db)

			result, err := repo.UpdateBookmark(ctx, tc.bookmarkID, tc.userID, tc.updates)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tc.verifyFunc != nil {
					tc.verifyFunc(t, result)
				}
			}
		})
	}
}

