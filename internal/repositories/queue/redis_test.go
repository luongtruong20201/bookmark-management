package queue

import (
	"testing"

	"github.com/redis/go-redis/v9"
	redisPkg "github.com/luongtruong20201/bookmark-management/pkg/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisQueue_PushMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		queueName   string
		msg         []byte
		closeBefore bool
		expectError error
	}{
		{
			name:        "pushes JSON message to redis list",
			queueName:   "test-queue",
			msg:         []byte(`{"hello":"world"}`),
			closeBefore: false,
			expectError: nil,
		},
		{
			name:        "pushes empty message",
			queueName:   "my-queue",
			msg:         []byte{},
			closeBefore: false,
			expectError: nil,
		},
		{
			name:        "pushes binary-like payload",
			queueName:   "import-queue",
			msg:         []byte(`{"uid":"user-1","bookmarks":[]}`),
			closeBefore: false,
			expectError: nil,
		},
		{
			name:        "redis connection closed returns error",
			queueName:   "test-queue",
			msg:         []byte(`{"hello":"world"}`),
			closeBefore: true,
			expectError: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			client := redisPkg.InitMockRedis(t)
			repo := NewRedisQueue(client, tc.queueName)

			if tc.closeBefore {
				client.Close()
			}

			err := repo.PushMessage(ctx, tc.msg)

			if tc.expectError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectError)
				return
			}
			require.NoError(t, err)

			res, err := client.LRange(ctx, tc.queueName, 0, -1).Result()
			require.NoError(t, err)
			assert.Len(t, res, 1)
			assert.Equal(t, string(tc.msg), res[0])
		})
	}
}
