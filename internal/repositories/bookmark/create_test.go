package bookmark

import (
	"context"
	"strings"
	"testing"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/internal/test/fixture"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRepository_CreateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		inputBookmark  *model.Bookmark
		expectedError  error
		expectedOutput *model.Bookmark
		verifyFunc     func(t *testing.T, db *gorm.DB, bookmark *model.Bookmark)
	}{
		{
			name: "success - create new bookmark",
			inputBookmark: &model.Bookmark{
				Base: model.Base{
					ID: "11111111-2222-3333-4444-555555555555",
				},
				Description: "My personal website",
				URL:         "https://truonglq.com",
				Code:        "mycode01",
				UserID:      "550e8400-e29b-41d4-a716-446655440000",
			},
			expectedError: nil,
			expectedOutput: &model.Bookmark{
				Base: model.Base{
					ID: "11111111-2222-3333-4444-555555555555",
				},
				Description: "My personal website",
				URL:         "https://truonglq.com",
				Code:        "mycode01",
				UserID:      "550e8400-e29b-41d4-a716-446655440000",
			},
			verifyFunc: func(t *testing.T, db *gorm.DB, bookmark *model.Bookmark) {
				var stored model.Bookmark
				err := db.Where("id = ?", bookmark.ID).First(&stored).Error
				assert.NoError(t, err)
				assert.Equal(t, bookmark.ID, stored.ID)
				assert.Equal(t, bookmark.Description, stored.Description)
				assert.Equal(t, bookmark.URL, stored.URL)
				assert.Equal(t, bookmark.Code, stored.Code)
				assert.Equal(t, bookmark.UserID, stored.UserID)
			},
		},
		{
			name: "success - auto-generate bookmark ID",
			inputBookmark: &model.Bookmark{
				Description: "Bookmark without ID",
				URL:         "https://example.com",
				Code:        "code0001",
				UserID:      "550e8400-e29b-41d4-a716-446655440000",
			},
			expectedError:  nil,
			expectedOutput: nil,
			verifyFunc: func(t *testing.T, db *gorm.DB, bookmark *model.Bookmark) {
				var stored model.Bookmark
				err := db.Where("code = ?", bookmark.Code).First(&stored).Error
				assert.NoError(t, err)
				assert.NotEmpty(t, stored.ID)
				assert.Equal(t, bookmark.Description, stored.Description)
				assert.Equal(t, bookmark.URL, stored.URL)
				assert.Equal(t, bookmark.Code, stored.Code)
				assert.Equal(t, bookmark.UserID, stored.UserID)
			},
		},
		{
			name: "error - duplicate bookmark ID",
			inputBookmark: &model.Bookmark{
				Base: model.Base{
					ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
				},
				Description: "Duplicate ID bookmark",
				URL:         "https://duplicate.com",
				Code:        "dupcode1",
				UserID:      "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
			},
			expectedError:  dbutils.ErrDuplicationType,
			expectedOutput: nil,
			verifyFunc:     nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			repo := NewBookmark(db)

			res, err := repo.CreateBookmark(ctx, tc.inputBookmark)

			if tc.expectedError != nil {
				assert.Error(t, err)

				if err != tc.expectedError {
					errStr := strings.ToLower(err.Error())
					assert.True(t, strings.Contains(errStr, "unique constraint") || err == dbutils.ErrDuplicationType,
						"expected duplicate constraint error, got: %v", err)
				}
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				if tc.expectedOutput != nil {
					assert.Equal(t, tc.expectedOutput.ID, res.ID)
					assert.Equal(t, tc.expectedOutput.Description, res.Description)
					assert.Equal(t, tc.expectedOutput.URL, res.URL)
					assert.Equal(t, tc.expectedOutput.Code, res.Code)
					assert.Equal(t, tc.expectedOutput.UserID, res.UserID)
				}
				if tc.verifyFunc != nil {
					tc.verifyFunc(t, db, res)
				}
			}
		})
	}
}
