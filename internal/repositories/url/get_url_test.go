package url

import (
	"context"
	"strings"
	"testing"

	redisPkg "github.com/luongtruong20201/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestURLStorage_Get(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupRedis     func(t *testing.T, ctx context.Context) *redis.Client
		expectedResult string
		expectedError  error
	}{
		{
			name: "success",
			setupRedis: func(t *testing.T, ctx context.Context) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				redis.Set(ctx, "1234567", "https://truonglq.com", 0)

				return redis
			},
			expectedResult: "https://truonglq.com",
			expectedError:  nil,
		},
		{
			name: "fail - key not exists",
			setupRedis: func(t *testing.T, ctx context.Context) *redis.Client {
				redis := redisPkg.InitMockRedis(t)

				return redis
			},
			expectedResult: "",
			expectedError:  redis.Nil,
		},
		{
			name: "fail - redis connection",
			setupRedis: func(t *testing.T, ctx context.Context) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				_ = redis.Close()

				return redis
			},
			expectedResult: "",
			expectedError:  redis.ErrClosed,
		},
		{
			name: "empty key",
			setupRedis: func(t *testing.T, ctx context.Context) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				return redis
			},
			expectedResult: "",
			expectedError:  redis.Nil,
		},
		{
			name: "very long key",
			setupRedis: func(t *testing.T, ctx context.Context) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				longKey := "a" + strings.Repeat("x", 1000) + "b"
				redis.Set(ctx, longKey, "https://example.com", 0)
				return redis
			},
			expectedResult: "https://example.com",
			expectedError:  nil,
		},
		{
			name: "key with special characters",
			setupRedis: func(t *testing.T, ctx context.Context) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				specialKey := "abc-123_test@key"
				redis.Set(ctx, specialKey, "https://special.com", 0)
				return redis
			},
			expectedResult: "https://special.com",
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			redis := tc.setupRedis(t, ctx)
			repo := NewURLStorage(redis)

			var key string
			switch tc.name {
			case "empty key":
				key = ""
			case "very long key":
				key = "a" + strings.Repeat("x", 1000) + "b"
			case "key with special characters":
				key = "abc-123_test@key"
			default:
				key = "1234567"
			}

			res, err := repo.Get(ctx, key)

			assert.Equal(t, res, tc.expectedResult)
			assert.Equal(t, err, tc.expectedError)
		})
	}
}
