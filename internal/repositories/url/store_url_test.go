package url

import (
	"context"
	"strings"
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
			code:           "1234567",
			url:            "https://truonglq.com",
			expectedResult: false,
			expectedError:  redis.ErrClosed,
		},
		{
			name: "store with custom expiration",
			setupMock: func(t *testing.T, ctx context.Context) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				return redis
			},
			code:           "custom123",
			url:            "https://example.com",
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "store with maximum expiration",
			setupMock: func(t *testing.T, ctx context.Context) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				return redis
			},
			code:           "maxexp01",
			url:            "https://maxexp.com",
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "store with empty code",
			setupMock: func(t *testing.T, ctx context.Context) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				return redis
			},
			code:           "",
			url:            "https://emptycode.com",
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "store with very long URL",
			setupMock: func(t *testing.T, ctx context.Context) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				return redis
			},
			code:           "longurl1",
			url:            "https://example.com/" + strings.Repeat("x", 500) + "end",
			expectedResult: true,
			expectedError:  nil,
		},
		{
			name: "store with special characters in code",
			setupMock: func(t *testing.T, ctx context.Context) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				return redis
			},
			code:           "abc-123_test@key",
			url:            "https://special.com",
			expectedResult: true,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			redis := tc.setupMock(t, ctx)
			repo := NewURLStorage(redis)

			expire := 0
			if tc.name == "store with custom expiration" {
				expire = 3600
			} else if tc.name == "store with maximum expiration" {
				expire = 604800
			}

			ok, err := repo.StoreIfNotExists(ctx, tc.code, tc.url, expire)
			assert.Equal(t, tc.expectedResult, ok)
			assert.Equal(t, tc.expectedError, err)

			if ok && tc.expectedError == nil {
				val, getErr := redis.Get(ctx, tc.code).Result()
				if getErr == nil {
					assert.Equal(t, val, tc.url)
				}
			}
		})
	}
}
