package bookmark

import (
	"context"
	"errors"
	"testing"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	bookmarkRepo "github.com/luongtruong20201/bookmark-management/internal/repositories/bookmark"
	repoMocks "github.com/luongtruong20201/bookmark-management/internal/repositories/bookmark/mocks"
	"github.com/luongtruong20201/bookmark-management/internal/test/fixture"
	mockKeyGen "github.com/luongtruong20201/bookmark-management/pkg/stringutils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkService_GetBookmarks(t *testing.T) {
	t.Parallel()

	var (
		testErrDatabase = errors.New("database error")
		mockUserID      = "550e8400-e29b-41d4-a716-446655440000"
	)

	testCases := []struct {
		name            string
		userID          string
		offset          int
		limit           int
		setupRepo       func(t *testing.T, ctx context.Context, userID string, offset, limit int) *repoMocks.Repository
		expectedError   error
		expectedCount   int
		verifyBookmarks func(t *testing.T, bookmarks []*model.Bookmark)
	}{
		{
			name:   "success - get bookmarks with pagination",
			userID: mockUserID,
			offset: 0,
			limit:  10,
			setupRepo: func(t *testing.T, ctx context.Context, userID string, offset, limit int) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				bookmarks := []*model.Bookmark{
					{
						Base: model.Base{
							ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
						},
						Description: "Facebook",
						URL:         "https://www.facebook.com",
						Code:        "abc1234",
						UserID:      userID,
					},
					{
						Base: model.Base{
							ID: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
						},
						Description: "Google",
						URL:         "https://www.google.com",
						Code:        "def5678",
						UserID:      userID,
					},
				}
				repo.On("GetBookmarks", ctx, userID, offset, limit).Return(bookmarks, nil).Once()
				return repo
			},
			expectedError: nil,
			expectedCount: 2,
			verifyBookmarks: func(t *testing.T, bookmarks []*model.Bookmark) {
				assert.Len(t, bookmarks, 2)
				assert.Equal(t, "Facebook", bookmarks[0].Description)
				assert.Equal(t, "Google", bookmarks[1].Description)
			},
		},
		{
			name:   "success - empty result",
			userID: mockUserID,
			offset: 0,
			limit:  10,
			setupRepo: func(t *testing.T, ctx context.Context, userID string, offset, limit int) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				repo.On("GetBookmarks", ctx, userID, offset, limit).Return([]*model.Bookmark{}, nil).Once()
				return repo
			},
			expectedError: nil,
			expectedCount: 0,
			verifyBookmarks: func(t *testing.T, bookmarks []*model.Bookmark) {
				assert.Empty(t, bookmarks)
			},
		},
		{
			name:   "error - repository returns error",
			userID: mockUserID,
			offset: 0,
			limit:  10,
			setupRepo: func(t *testing.T, ctx context.Context, userID string, offset, limit int) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				repo.On("GetBookmarks", ctx, userID, offset, limit).Return(nil, testErrDatabase).Once()
				return repo
			},
			expectedError:   testErrDatabase,
			expectedCount:   0,
			verifyBookmarks: nil,
		},
		{
			name:   "success - with offset",
			userID: mockUserID,
			offset: 10,
			limit:  5,
			setupRepo: func(t *testing.T, ctx context.Context, userID string, offset, limit int) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				bookmarks := []*model.Bookmark{
					{
						Base: model.Base{
							ID: "c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f",
						},
						Description: "GitHub",
						URL:         "https://www.github.com",
						Code:        "ghi9012",
						UserID:      userID,
					},
				}
				repo.On("GetBookmarks", ctx, userID, offset, limit).Return(bookmarks, nil).Once()
				return repo
			},
			expectedError: nil,
			expectedCount: 1,
			verifyBookmarks: func(t *testing.T, bookmarks []*model.Bookmark) {
				assert.Len(t, bookmarks, 1)
				assert.Equal(t, "GitHub", bookmarks[0].Description)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			repo := tc.setupRepo(t, ctx, tc.userID, tc.offset, tc.limit)
			keyGen := mockKeyGen.NewKeyGenerator(t)
			svc := NewBookmarkSvc(repo, keyGen)

			bookmarks, err := svc.GetBookmarks(ctx, tc.userID, tc.offset, tc.limit)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, bookmarks)
			} else {
				assert.NoError(t, err)
				assert.Len(t, bookmarks, tc.expectedCount)
				if tc.verifyBookmarks != nil {
					tc.verifyBookmarks(t, bookmarks)
				}
			}
		})
	}
}

func TestBookmarkService_GetBookmarks_WithFixture(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		userID          string
		offset          int
		limit           int
		expectedCount   int
		verifyBookmarks func(t *testing.T, bookmarks []*model.Bookmark)
	}{
		{
			name:          "success - get bookmarks with fixture",
			userID:        "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
			offset:        0,
			limit:         10,
			expectedCount: 2,
			verifyBookmarks: func(t *testing.T, bookmarks []*model.Bookmark) {
				assert.Len(t, bookmarks, 2)
				assert.Equal(t, "Facebook - Social Media Platform", bookmarks[0].Description)
				assert.Equal(t, "Google - Search Engine", bookmarks[1].Description)
			},
		},
		{
			name:          "success - get bookmarks with pagination limit",
			userID:        "2f6c9d14-9e42-4c31-9e6a-0f8c4b2a7c55",
			offset:        0,
			limit:         1,
			expectedCount: 1,
			verifyBookmarks: func(t *testing.T, bookmarks []*model.Bookmark) {
				assert.Len(t, bookmarks, 1)
				assert.Equal(t, "GitHub - Code Repository", bookmarks[0].Description)
			},
		},
		{
			name:          "success - get bookmarks with offset",
			userID:        "2f6c9d14-9e42-4c31-9e6a-0f8c4b2a7c55",
			offset:        1,
			limit:         10,
			expectedCount: 1,
			verifyBookmarks: func(t *testing.T, bookmarks []*model.Bookmark) {
				assert.Len(t, bookmarks, 1)
				assert.Equal(t, "YouTube - Video Platform", bookmarks[0].Description)
			},
		},
		{
			name:          "success - get bookmarks for user with no bookmarks",
			userID:        "00000000-0000-0000-0000-000000000000",
			offset:        0,
			limit:         10,
			expectedCount: 0,
			verifyBookmarks: func(t *testing.T, bookmarks []*model.Bookmark) {
				assert.Empty(t, bookmarks)
			},
		},
		{
			name:          "success - get bookmarks for johndoe",
			userID:        "550e8400-e29b-41d4-a716-446655440000",
			offset:        0,
			limit:         10,
			expectedCount: 2,
			verifyBookmarks: func(t *testing.T, bookmarks []*model.Bookmark) {
				assert.Len(t, bookmarks, 2)
				for _, bookmark := range bookmarks {
					assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", bookmark.UserID)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			repo := bookmarkRepo.NewBookmark(db)
			keyGen := mockKeyGen.NewKeyGenerator(t)
			svc := NewBookmarkSvc(repo, keyGen)

			bookmarks, err := svc.GetBookmarks(ctx, tc.userID, tc.offset, tc.limit)

			assert.NoError(t, err)
			assert.Len(t, bookmarks, tc.expectedCount)
			if tc.verifyBookmarks != nil {
				tc.verifyBookmarks(t, bookmarks)
			}
		})
	}
}

func TestBookmarkService_CountBookmarks(t *testing.T) {
	t.Parallel()

	var (
		testErrDatabase = errors.New("database error")
		mockUserID      = "550e8400-e29b-41d4-a716-446655440000"
	)

	testCases := []struct {
		name          string
		userID        string
		setupRepo     func(t *testing.T, ctx context.Context, userID string) *repoMocks.Repository
		expectedError error
		expectedTotal int64
	}{
		{
			name:   "success - count bookmarks",
			userID: mockUserID,
			setupRepo: func(t *testing.T, ctx context.Context, userID string) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				repo.On("CountBookmarks", ctx, userID).Return(int64(25), nil).Once()
				return repo
			},
			expectedError: nil,
			expectedTotal: 25,
		},
		{
			name:   "success - zero bookmarks",
			userID: mockUserID,
			setupRepo: func(t *testing.T, ctx context.Context, userID string) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				repo.On("CountBookmarks", ctx, userID).Return(int64(0), nil).Once()
				return repo
			},
			expectedError: nil,
			expectedTotal: 0,
		},
		{
			name:   "error - repository returns error",
			userID: mockUserID,
			setupRepo: func(t *testing.T, ctx context.Context, userID string) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				repo.On("CountBookmarks", ctx, userID).Return(int64(0), testErrDatabase).Once()
				return repo
			},
			expectedError: testErrDatabase,
			expectedTotal: 0,
		},
		{
			name:   "success - large count",
			userID: mockUserID,
			setupRepo: func(t *testing.T, ctx context.Context, userID string) *repoMocks.Repository {
				repo := repoMocks.NewRepository(t)
				repo.On("CountBookmarks", ctx, userID).Return(int64(9999), nil).Once()
				return repo
			},
			expectedError: nil,
			expectedTotal: 9999,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			repo := tc.setupRepo(t, ctx, tc.userID)
			keyGen := mockKeyGen.NewKeyGenerator(t)
			svc := NewBookmarkSvc(repo, keyGen)

			total, err := svc.CountBookmarks(ctx, tc.userID)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Equal(t, int64(0), total)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedTotal, total)
			}
		})
	}
}

func TestBookmarkService_CountBookmarks_WithFixture(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		userID        string
		expectedTotal int64
	}{
		{
			name:          "success - count bookmarks for user with multiple bookmarks",
			userID:        "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91",
			expectedTotal: 2,
		},
		{
			name:          "success - count bookmarks for another user",
			userID:        "2f6c9d14-9e42-4c31-9e6a-0f8c4b2a7c55",
			expectedTotal: 2,
		},
		{
			name:          "success - count bookmarks for johndoe",
			userID:        "550e8400-e29b-41d4-a716-446655440000",
			expectedTotal: 2,
		},
		{
			name:          "success - count bookmarks for user with no bookmarks",
			userID:        "00000000-0000-0000-0000-000000000000",
			expectedTotal: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := context.Background()
			db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			repo := bookmarkRepo.NewBookmark(db)
			keyGen := mockKeyGen.NewKeyGenerator(t)
			svc := NewBookmarkSvc(repo, keyGen)

			total, err := svc.CountBookmarks(ctx, tc.userID)

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedTotal, total)
		})
	}
}
