package queue

import (
	"context"
	"testing"

	redisPkg "github.com/luongtruong20201/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
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

func TestRedisQueue_PopMessage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name  string
		setup func(t *testing.T, repo Repository, ctx context.Context)
		run   func(t *testing.T, repo Repository, ctx context.Context)
	}{
		{
			name:  "empty_queue_returns_ErrNoMessage",
			setup: func(t *testing.T, repo Repository, ctx context.Context) {},
			run: func(t *testing.T, repo Repository, ctx context.Context) {
				msg, err := repo.PopMessage(ctx)
				assert.Nil(t, msg)
				assert.ErrorIs(t, err, ErrNoMessage)
			},
		},
		{
			name: "after_push_returns_messages_in_order",
			setup: func(t *testing.T, repo Repository, ctx context.Context) {
				require.NoError(t, repo.PushMessage(ctx, []byte("a")))
				require.NoError(t, repo.PushMessage(ctx, []byte("b")))
				require.NoError(t, repo.PushMessage(ctx, []byte("c")))
			},
			run: func(t *testing.T, repo Repository, ctx context.Context) {
				msg1, err := repo.PopMessage(ctx)
				require.NoError(t, err)
				assert.Equal(t, "a", string(msg1))

				msg2, err := repo.PopMessage(ctx)
				require.NoError(t, err)
				assert.Equal(t, "b", string(msg2))

				msg3, err := repo.PopMessage(ctx)
				require.NoError(t, err)
				assert.Equal(t, "c", string(msg3))

				msg4, err := repo.PopMessage(ctx)
				assert.Nil(t, msg4)
				assert.ErrorIs(t, err, ErrNoMessage)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			client := redisPkg.InitMockRedis(t)
			repo := NewRedisQueue(client, "pop-queue")
			ctx := t.Context()
			tc.setup(t, repo, ctx)
			tc.run(t, repo, ctx)
		})
	}
}
