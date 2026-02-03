package shorten

import (
	"context"
	"testing"

	mockStorage "github.com/luongtruong20201/bookmark-management/internal/repositories/url/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupRepo      func(t *testing.T, ctx context.Context) *mockStorage.URLStorage
		expectedResult string
		expectedError  error
	}{
		{
			name: "success",
			setupRepo: func(t *testing.T, ctx context.Context) *mockStorage.URLStorage {
				repo := mockStorage.NewURLStorage(t)
				repo.On("Get", ctx, "1234567").Return("https://truonglq.com", nil).Once()

				return repo
			},
			expectedResult: "https://truonglq.com",
			expectedError:  nil,
		},
		{
			name: "fail - key not exists",
			setupRepo: func(t *testing.T, ctx context.Context) *mockStorage.URLStorage {
				repo := mockStorage.NewURLStorage(t)
				repo.On("Get", ctx, "1234567").Return("", ErrCodeNotFound).Once()

				return repo
			},
			expectedResult: "",
			expectedError:  ErrCodeNotFound,
		},
		{
			name: "fail - redis connection",
			setupRepo: func(t *testing.T, ctx context.Context) *mockStorage.URLStorage {
				repo := mockStorage.NewURLStorage(t)
				repo.On("Get", ctx, "1234567").Return("", redis.ErrClosed).Once()

				return repo
			},
			expectedResult: "",
			expectedError:  redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			repo := tc.setupRepo(t, ctx)

			res, err := repo.Get(ctx, "1234567")
			assert.Equal(t, res, tc.expectedResult)
			assert.Equal(t, err, tc.expectedError)
		})
	}
}
