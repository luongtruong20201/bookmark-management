package repository

import (
	"context"
	"testing"

	redisPkg "github.com/luongtruong20201/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestURLStorage_StoreIfNotExists(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupMock      func(t *testing.T, ctx context.Context) *redis.Client
		code           string
		url            string
		expectedResult bool
		expectedError  error
	}{
		{
			name: "store success",
			setupMock: func(t *testing.T, ctx context.Context) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				return redis
			},
			code:           "1234567",
			url:            "https://truonglq.com",
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "exists key",
			setupMock: func(t *testing.T, ctx context.Context) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				redis.Set(ctx, "1234567", "https://truonglq.com", 0)

				return redis
			},
			code:           "1234567",
			url:            "https://truonglq.com",
			expectedResult: false,
			expectedError:  nil,
		},
		{
			name: "lost connection",
			setupMock: func(t *testing.T, ctx context.Context) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				_ = redis.Close()

				return redis
			},
			expectedResult: false,
			expectedError:  redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			redis := tc.setupMock(t, ctx)
			repo := NewURLStorage(redis)

			ok, err := repo.StoreIfNotExists(ctx, tc.code, tc.url, 0)
			assert.Equal(t, tc.expectedResult, ok)
			assert.Equal(t, tc.expectedError, err)

			if ok {
				val, _ := redis.Get(ctx, tc.code).Result()
				assert.Equal(t, val, tc.url)
			}
		})
	}
}
