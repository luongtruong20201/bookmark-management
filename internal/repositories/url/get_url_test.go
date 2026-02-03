package url

import (
	"context"
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
	}

	for _, tc := range testCases {
		ctx := t.Context()
		redis := tc.setupRedis(t, ctx)
		repo := NewURLStorage(redis)

		res, err := repo.Get(ctx, "1234567")

		assert.Equal(t, res, tc.expectedResult)
		assert.Equal(t, err, tc.expectedError)
	}
}
