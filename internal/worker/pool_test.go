package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/luongtruong20201/bookmark-management/internal/worker/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errHandler = errors.New("handler error")

func Test_Pool(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		run  func(t *testing.T, ctx context.Context, cancel context.CancelFunc)
	}{
		{
			name: "NewPool_returns_non_nil",
			run: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) {
				mockHandler := mocks.NewHandler(t)
				pool := NewPool(ctx, mockHandler, 2)
				require.NotNil(t, pool)
				pool.Close()
			},
		},
		{
			name: "Consume_delivers_single_message_to_handler",
			run: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) {
				mockHandler := mocks.NewHandler(t)
				mockHandler.On("Handle", ctx, []byte("msg1")).Return(nil).Once()
				pool := NewPool(ctx, mockHandler, 2)
				pool.Consume([]byte("msg1"))
				time.Sleep(200 * time.Millisecond)
				pool.Close()
			},
		},
		{
			name: "Consume_delivers_multiple_messages_to_handler",
			run: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) {
				mockHandler := mocks.NewHandler(t)
				for _, msg := range [][]byte{[]byte("a"), []byte("b"), []byte("c")} {
					mockHandler.On("Handle", ctx, msg).Return(nil).Once()
				}
				pool := NewPool(ctx, mockHandler, 2)
				pool.Consume([]byte("a"))
				pool.Consume([]byte("b"))
				pool.Consume([]byte("c"))
				time.Sleep(200 * time.Millisecond)
				pool.Close()
			},
		},
		{
			name: "Close_does_not_panic",
			run: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) {
				mockHandler := mocks.NewHandler(t)
				pool := NewPool(ctx, mockHandler, 2)
				require.NotPanics(t, func() { pool.Close() })
			},
		},
		{
			name: "Consume_handler_returns_error_pool_continues",
			run: func(t *testing.T, ctx context.Context, cancel context.CancelFunc) {
				mockHandler := mocks.NewHandler(t)
				mockHandler.On("Handle", ctx, mock.Anything).Return(errHandler).Once()
				mockHandler.On("Handle", ctx, mock.Anything).Return(nil).Once()
				pool := NewPool(ctx, mockHandler, 2)
				pool.Consume([]byte("err-msg"))
				pool.Consume([]byte("ok-msg"))
				time.Sleep(200 * time.Millisecond)
				pool.Close()
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			tc.run(t, ctx, cancel)
		})
	}
}
