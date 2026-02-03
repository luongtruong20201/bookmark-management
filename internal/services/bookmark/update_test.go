package bookmark

import (
	"context"
	"errors"
	"testing"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	repoMocks "github.com/luongtruong20201/bookmark-management/internal/repositories/bookmark/mocks"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkService_Update(t *testing.T) {
	t.Parallel()

	var (
		testErrDatabase = errors.New("database error")
	)

	testCases := []struct {
		name           string
		setupRepo      func(t *testing.T, ctx context.Context, bookmarkID, userID string, updates *model.Bookmark) *repoMocks.Repository
		bookmarkID     string
		userID         string
		description    string
		url            string
		expectedError  error
		verifyBookmark func(t *testing.T, bookmark *model.Bookmark)
	}{
		{
			name:        "success - update bookmark with description and URL",
			bookmarkID:  "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			description: "Updated Facebook",
			url:         "https://www.facebook.com/updated",
			setupRepo: func(t *testing.T, ctx context.Context, bookmarkID, userID string, updates *model.Bookmark) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				output := &model.Bookmark{
					Base: model.Base{
						ID: bookmarkID,
					},
					Description: updates.Description,
					URL:         updates.URL,
					Code:        "abc1234",
					UserID:      userID,
				}
				repo.On("UpdateBookmark", ctx, bookmarkID, userID, updates).Return(output, nil).Once()
				return repo
			},
			expectedError: nil,
			verifyBookmark: func(t *testing.T, bookmark *model.Bookmark) {
				assert.Equal(t, "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d", bookmark.ID)
				assert.Equal(t, "Updated Facebook", bookmark.Description)
				assert.Equal(t, "https://www.facebook.com/updated", bookmark.URL)
				assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", bookmark.UserID)
			},
		},
		{
			name:        "success - update bookmark with only description",
			bookmarkID:  "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			description: "Updated Description Only",
			url:         "",
			setupRepo: func(t *testing.T, ctx context.Context, bookmarkID, userID string, updates *model.Bookmark) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				output := &model.Bookmark{
					Base: model.Base{
						ID: bookmarkID,
					},
					Description: updates.Description,
					URL:         "https://www.facebook.com",
					Code:        "abc1234",
					UserID:      userID,
				}
				repo.On("UpdateBookmark", ctx, bookmarkID, userID, updates).Return(output, nil).Once()
				return repo
			},
			expectedError: nil,
			verifyBookmark: func(t *testing.T, bookmark *model.Bookmark) {
				assert.Equal(t, "Updated Description Only", bookmark.Description)
				assert.Equal(t, "https://www.facebook.com", bookmark.URL)
			},
		},
		{
			name:        "success - update bookmark with only URL",
			bookmarkID:  "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			description: "",
			url:         "https://www.google.com/updated",
			setupRepo: func(t *testing.T, ctx context.Context, bookmarkID, userID string, updates *model.Bookmark) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				output := &model.Bookmark{
					Base: model.Base{
						ID: bookmarkID,
					},
					Description: "Facebook",
					URL:         updates.URL,
					Code:        "abc1234",
					UserID:      userID,
				}
				repo.On("UpdateBookmark", ctx, bookmarkID, userID, updates).Return(output, nil).Once()
				return repo
			},
			expectedError: nil,
			verifyBookmark: func(t *testing.T, bookmark *model.Bookmark) {
				assert.Equal(t, "Facebook", bookmark.Description)
				assert.Equal(t, "https://www.google.com/updated", bookmark.URL)
			},
		},
		{
			name:        "error - bookmark not found",
			bookmarkID:  "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			description: "Updated Facebook",
			url:         "https://www.facebook.com/updated",
			setupRepo: func(t *testing.T, ctx context.Context, bookmarkID, userID string, updates *model.Bookmark) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				repo.On("UpdateBookmark", ctx, bookmarkID, userID, updates).Return(nil, dbutils.ErrNotFoundType).Once()
				return repo
			},
			expectedError:  dbutils.ErrNotFoundType,
			verifyBookmark: nil,
		},
		{
			name:        "error - repository error",
			bookmarkID:  "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
			userID:      "550e8400-e29b-41d4-a716-446655440000",
			description: "Updated Facebook",
			url:         "https://www.facebook.com/updated",
			setupRepo: func(t *testing.T, ctx context.Context, bookmarkID, userID string, updates *model.Bookmark) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				repo.On("UpdateBookmark", ctx, bookmarkID, userID, updates).Return(nil, testErrDatabase).Once()
				return repo
			},
			expectedError:  testErrDatabase,
			verifyBookmark: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			repo := tc.setupRepo(t, ctx, tc.bookmarkID, tc.userID, &model.Bookmark{
				Description: tc.description,
				URL:         tc.url,
			})

			svc := NewBookmarkSvc(repo, nil)

			result, err := svc.Update(ctx, tc.bookmarkID, tc.userID, tc.description, tc.url)

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

