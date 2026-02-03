package bookmark

import (
	"context"
	"testing"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/internal/test/fixture"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRepository_DeleteBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		bookmarkID    string
		userID        string
		expectedError error
		verifyFunc    func(t *testing.T, db interface{})
	}{
		{
			name:          "success - delete bookmark",
			bookmarkID:    "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
			userID:        "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
			expectedError: nil,
			verifyFunc: func(t *testing.T, gormDB interface{}) {
				db := gormDB.(*gorm.DB)
				var bookmark model.Bookmark
				err := db.Where("id = ?", "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d").First(&bookmark).Error
				assert.Error(t, err, "bookmark should be deleted")
			},
		},
		{
			name:          "error - bookmark not found",
			bookmarkID:    "00000000-0000-0000-0000-000000000000",
			userID:        "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
			expectedError: dbutils.ErrNotFoundType,
			verifyFunc:    nil,
		},
		{
			name:          "error - bookmark belongs to different user",
			bookmarkID:    "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
			userID:        "550e8400-e29b-41d4-a716-446655440000",
			expectedError: dbutils.ErrNotFoundType,
			verifyFunc:    nil,
		},
		{
			name:          "success - delete bookmark and verify other bookmarks remain",
			bookmarkID:    "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
			userID:        "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
			expectedError: nil,
			verifyFunc: func(t *testing.T, gormDB interface{}) {
				db := gormDB.(*gorm.DB)
				var deletedBookmark model.Bookmark
				err := db.Where("id = ?", "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e").First(&deletedBookmark).Error
				assert.Error(t, err, "deleted bookmark should not exist")

				var remainingBookmark model.Bookmark
				err = db.Where("id = ? AND user_id = ?", "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d", "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91").First(&remainingBookmark).Error
				assert.NoError(t, err, "other bookmarks should still exist")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			repo := NewBookmark(db)

			err := repo.DeleteBookmark(ctx, tc.bookmarkID, tc.userID)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
				if tc.verifyFunc != nil {
					tc.verifyFunc(t, db)
				}
			}
		})
	}
}

