package worker

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/luongtruong20201/bookmark-management/internal/worker/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_worker(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		messages [][]byte
	}{
		{
			name:     "Work_delivers_messages_to_handler",
			messages: [][]byte{[]byte("job1"), []byte("job2")},
		},
		{
			name:     "Work_handler_called_with_correct_payload",
			messages: [][]byte{[]byte("payload")},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			mockHandler := mocks.NewHandler(t)
			var mu sync.Mutex
			var received [][]byte
			for _, msg := range tc.messages {
				m := msg
				mockHandler.On("Handle", mock.Anything, m).Run(func(args mock.Arguments) {
					mu.Lock()
					received = append(received, args.Get(1).([]byte))
					mu.Unlock()
				}).Return(nil).Once()
			}

			pool := NewPool(ctx, mockHandler, 2)
			for _, msg := range tc.messages {
				pool.Consume(msg)
			}
			time.Sleep(300 * time.Millisecond)
			pool.Close()
			cancel()

			mu.Lock()
			require.Equal(t, len(tc.messages), len(received), "handler should receive all messages")
			mu.Unlock()
		})
	}
}
