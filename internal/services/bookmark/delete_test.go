package bookmark

import (
	"context"
	"errors"
	"testing"

	repoMocks "github.com/luongtruong20201/bookmark-management/internal/repositories/bookmark/mocks"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkService_Delete(t *testing.T) {
	t.Parallel()

	var (
		testErrDatabase = errors.New("database error")
	)

	testCases := []struct {
		name          string
		setupRepo     func(t *testing.T, ctx context.Context, bookmarkID, userID string) *repoMocks.Repository
		bookmarkID    string
		userID        string
		expectedError error
	}{
		{
			name:       "success - delete bookmark",
			bookmarkID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
			userID:     "550e8400-e29b-41d4-a716-446655440000",
			setupRepo: func(t *testing.T, ctx context.Context, bookmarkID, userID string) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				repo.On("DeleteBookmark", ctx, bookmarkID, userID).Return(nil).Once()
				return repo
			},
			expectedError: nil,
		},
		{
			name:       "error - bookmark not found",
			bookmarkID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
			userID:     "550e8400-e29b-41d4-a716-446655440000",
			setupRepo: func(t *testing.T, ctx context.Context, bookmarkID, userID string) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				repo.On("DeleteBookmark", ctx, bookmarkID, userID).Return(dbutils.ErrNotFoundType).Once()
				return repo
			},
			expectedError: dbutils.ErrNotFoundType,
		},
		{
			name:       "error - repository error",
			bookmarkID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
			userID:     "550e8400-e29b-41d4-a716-446655440000",
			setupRepo: func(t *testing.T, ctx context.Context, bookmarkID, userID string) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				repo.On("DeleteBookmark", ctx, bookmarkID, userID).Return(testErrDatabase).Once()
				return repo
			},
			expectedError: testErrDatabase,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			repo := tc.setupRepo(t, ctx, tc.bookmarkID, tc.userID)

			svc := NewBookmarkSvc(repo, nil)

			err := svc.Delete(ctx, tc.bookmarkID, tc.userID)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

