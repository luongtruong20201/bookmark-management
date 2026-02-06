package bookmark

import (
	"context"
	"strings"
	"testing"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/internal/test/fixture"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
)

func TestRepository_GetBookmarkByCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		code          string
		expectedError error
		verifyFunc    func(t *testing.T, bookmark *model.Bookmark)
	}{
		{
			name:          "success - get bookmark by code",
			code:          "abc12345",
			expectedError: nil,
			verifyFunc: func(t *testing.T, bookmark *model.Bookmark) {
				assert.NotNil(t, bookmark)
				assert.Equal(t, "abc12345", bookmark.Code)
				assert.Equal(t, "Facebook - Social Media Platform", bookmark.Description)
				assert.Equal(t, "https://www.facebook.com", bookmark.URL)
				assert.Equal(t, "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91", bookmark.UserID)
			},
		},
		{
			name:          "success - get bookmark by another code",
			code:          "def56789",
			expectedError: nil,
			verifyFunc: func(t *testing.T, bookmark *model.Bookmark) {
				assert.NotNil(t, bookmark)
				assert.Equal(t, "def56789", bookmark.Code)
				assert.Equal(t, "Google - Search Engine", bookmark.Description)
				assert.Equal(t, "https://www.google.com", bookmark.URL)
			},
		},
		{
			name:          "error - bookmark not found",
			code:          "nonexist",
			expectedError: dbutils.ErrNotFoundType,
			verifyFunc: func(t *testing.T, bookmark *model.Bookmark) {
				assert.Nil(t, bookmark)
			},
		},
		{
			name:          "success - get bookmark from different user",
			code:          "ghi90123",
			expectedError: nil,
			verifyFunc: func(t *testing.T, bookmark *model.Bookmark) {
				assert.NotNil(t, bookmark)
				assert.Equal(t, "ghi90123", bookmark.Code)
				assert.Equal(t, "GitHub - Code Repository", bookmark.Description)
				assert.Equal(t, "2f6c9d14-9e42-4c31-9e6a-0f8c4b2a7c55", bookmark.UserID)
			},
		},
		{
			name:          "error - empty code",
			code:          "",
			expectedError: dbutils.ErrNotFoundType,
			verifyFunc: func(t *testing.T, bookmark *model.Bookmark) {
				assert.Nil(t, bookmark)
			},
		},
		{
			name:          "error - very long code",
			code:          "a" + strings.Repeat("x", 1000) + "b",
			expectedError: dbutils.ErrNotFoundType,
			verifyFunc: func(t *testing.T, bookmark *model.Bookmark) {
				assert.Nil(t, bookmark)
			},
		},
		{
			name:          "error - code with special characters",
			code:          "abc-123_test@key",
			expectedError: dbutils.ErrNotFoundType,
			verifyFunc: func(t *testing.T, bookmark *model.Bookmark) {
				assert.Nil(t, bookmark)
			},
		},
		{
			name:          "error - code with SQL injection attempt",
			code:          "'; DROP TABLE bookmarks; --",
			expectedError: dbutils.ErrNotFoundType,
			verifyFunc: func(t *testing.T, bookmark *model.Bookmark) {
				assert.Nil(t, bookmark)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			repo := NewBookmark(db)

			bookmark, err := repo.GetBookmarkByCode(ctx, tc.code)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err)
				assert.Nil(t, bookmark)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, bookmark)
				if tc.verifyFunc != nil {
					tc.verifyFunc(t, bookmark)
				}
			}
		})
	}
}

