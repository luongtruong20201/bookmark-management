package shorten

import (
	"context"
	"errors"
	"testing"

	model "github.com/luongtruong20201/bookmark-management/internal/models"
	mockBookmarkRepo "github.com/luongtruong20201/bookmark-management/internal/repositories/bookmark/mocks"
	mockStorage "github.com/luongtruong20201/bookmark-management/internal/repositories/url/mocks"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		code           string
		setupRepo      func(t *testing.T, ctx context.Context) *mockStorage.URLStorage
		setupBookmarkRepo func(t *testing.T, ctx context.Context) *mockBookmarkRepo.Repository
		expectedResult string
		expectedError  error
	}{
		{
			name: "success - redis code (7 chars)",
			code: "1234567",
			setupRepo: func(t *testing.T, ctx context.Context) *mockStorage.URLStorage {
				repo := mockStorage.NewURLStorage(t)
				repo.On("Get", ctx, "1234567").Return("https://truonglq.com", nil).Once()
				return repo
			},
			setupBookmarkRepo: func(t *testing.T, ctx context.Context) *mockBookmarkRepo.Repository {
				return mockBookmarkRepo.NewRepository(t)
			},
			expectedResult: "https://truonglq.com",
			expectedError:  nil,
		},
		{
			name: "fail - redis code not found",
			code: "1234567",
			setupRepo: func(t *testing.T, ctx context.Context) *mockStorage.URLStorage {
				repo := mockStorage.NewURLStorage(t)
				repo.On("Get", ctx, "1234567").Return("", redis.Nil).Once()
				return repo
			},
			setupBookmarkRepo: func(t *testing.T, ctx context.Context) *mockBookmarkRepo.Repository {
				return mockBookmarkRepo.NewRepository(t)
			},
			expectedResult: "",
			expectedError:  ErrCodeNotFound,
		},
		{
			name: "fail - redis connection error",
			code: "1234567",
			setupRepo: func(t *testing.T, ctx context.Context) *mockStorage.URLStorage {
				repo := mockStorage.NewURLStorage(t)
				repo.On("Get", ctx, "1234567").Return("", redis.ErrClosed).Once()
				return repo
			},
			setupBookmarkRepo: func(t *testing.T, ctx context.Context) *mockBookmarkRepo.Repository {
				return mockBookmarkRepo.NewRepository(t)
			},
			expectedResult: "",
			expectedError:  redis.ErrClosed,
		},
		{
			name: "success - bookmark code (8 chars)",
			code: "12345678",
			setupRepo: func(t *testing.T, ctx context.Context) *mockStorage.URLStorage {
				return mockStorage.NewURLStorage(t)
			},
			setupBookmarkRepo: func(t *testing.T, ctx context.Context) *mockBookmarkRepo.Repository {
				repo := mockBookmarkRepo.NewRepository(t)
				bookmark := &model.Bookmark{
					URL: "https://example.com",
				}
				repo.On("GetBookmarkByCode", ctx, "12345678").Return(bookmark, nil).Once()
				return repo
			},
			expectedResult: "https://example.com",
			expectedError:  nil,
		},
		{
			name: "fail - bookmark code not found",
			code: "12345678",
			setupRepo: func(t *testing.T, ctx context.Context) *mockStorage.URLStorage {
				return mockStorage.NewURLStorage(t)
			},
			setupBookmarkRepo: func(t *testing.T, ctx context.Context) *mockBookmarkRepo.Repository {
				repo := mockBookmarkRepo.NewRepository(t)
				repo.On("GetBookmarkByCode", ctx, "12345678").Return(nil, dbutils.ErrNotFoundType).Once()
				return repo
			},
			expectedResult: "",
			expectedError:  ErrCodeNotFound,
		},
		{
			name: "fail - invalid code length (too short)",
			code: "12345",
			setupRepo: func(t *testing.T, ctx context.Context) *mockStorage.URLStorage {
				return mockStorage.NewURLStorage(t)
			},
			setupBookmarkRepo: func(t *testing.T, ctx context.Context) *mockBookmarkRepo.Repository {
				return mockBookmarkRepo.NewRepository(t)
			},
			expectedResult: "",
			expectedError:  ErrCodeNotFound,
		},
		{
			name: "fail - invalid code length (too long)",
			code: "123456789",
			setupRepo: func(t *testing.T, ctx context.Context) *mockStorage.URLStorage {
				return mockStorage.NewURLStorage(t)
			},
			setupBookmarkRepo: func(t *testing.T, ctx context.Context) *mockBookmarkRepo.Repository {
				return mockBookmarkRepo.NewRepository(t)
			},
			expectedResult: "",
			expectedError:  ErrCodeNotFound,
		},
		{
			name: "fail - empty code",
			code: "",
			setupRepo: func(t *testing.T, ctx context.Context) *mockStorage.URLStorage {
				return mockStorage.NewURLStorage(t)
			},
			setupBookmarkRepo: func(t *testing.T, ctx context.Context) *mockBookmarkRepo.Repository {
				return mockBookmarkRepo.NewRepository(t)
			},
			expectedResult: "",
			expectedError:  ErrCodeNotFound,
		},
		{
			name: "fail - bookmark repository database error",
			code: "12345678",
			setupRepo: func(t *testing.T, ctx context.Context) *mockStorage.URLStorage {
				return mockStorage.NewURLStorage(t)
			},
			setupBookmarkRepo: func(t *testing.T, ctx context.Context) *mockBookmarkRepo.Repository {
				repo := mockBookmarkRepo.NewRepository(t)
				repo.On("GetBookmarkByCode", ctx, "12345678").Return(nil, errors.New("database connection error")).Once()
				return repo
			},
			expectedResult: "",
			expectedError:  errors.New("database connection error"),
		},
		{
			name: "success - bookmark with empty URL field",
			code: "12345678",
			setupRepo: func(t *testing.T, ctx context.Context) *mockStorage.URLStorage {
				return mockStorage.NewURLStorage(t)
			},
			setupBookmarkRepo: func(t *testing.T, ctx context.Context) *mockBookmarkRepo.Repository {
				repo := mockBookmarkRepo.NewRepository(t)
				bookmark := &model.Bookmark{
					URL: "",
				}
				repo.On("GetBookmarkByCode", ctx, "12345678").Return(bookmark, nil).Once()
				return repo
			},
			expectedResult: "",
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			repo := tc.setupRepo(t, ctx)
			bookmarkRepo := tc.setupBookmarkRepo(t, ctx)
			svc := &shortenURL{
				repository:   repo,
				bookmarkRepo: bookmarkRepo,
			}

			res, err := svc.GetURL(ctx, tc.code)
			assert.Equal(t, tc.expectedResult, res)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
