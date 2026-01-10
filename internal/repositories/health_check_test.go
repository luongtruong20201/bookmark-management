package repository

import (
	"testing"

	redisPkg "github.com/luongtruong20201/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_Ping(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		expectedError error
		setupRedis    func(t *testing.T) *redis.Client
	}{
		{
			name:          "success",
			expectedError: nil,
			setupRedis: func(t *testing.T) *redis.Client {
				redis := redisPkg.InitMockRedis(t)

				return redis
			},
		},
		{
			name:          "fail",
			expectedError: redis.ErrClosed,
			setupRedis: func(t *testing.T) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				_ = redis.Close()

				return redis
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			redis := tc.setupRedis(t)
			healthCheck := NewHealthCheck(redis)
			ctx := t.Context()

			err := healthCheck.Ping(ctx)
			assert.Equal(t, err, tc.expectedError)
		})
	}
}
