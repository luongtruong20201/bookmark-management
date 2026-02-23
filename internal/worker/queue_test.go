package worker

import (
	"testing"

	"github.com/luongtruong20201/bookmark-management/internal/repositories/queue"
	"github.com/luongtruong20201/bookmark-management/internal/worker/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_Queue(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		setupMock func(t *testing.T) Queue
		wantMsg   []byte
		wantErr   error
	}{
		{
			name: "PopMessage_returns_message_when_available",
			setupMock: func(t *testing.T) Queue {
				mockQueue := mocks.NewQueue(t)
				msg := []byte("queue-msg")
				mockQueue.On("PopMessage", mock.Anything).Return(msg, nil).Once()
				return mockQueue
			},
			wantMsg: []byte("queue-msg"),
			wantErr: nil,
		},
		{
			name: "PopMessage_returns_ErrNoMessage_when_empty",
			setupMock: func(t *testing.T) Queue {
				mockQueue := mocks.NewQueue(t)
				mockQueue.On("PopMessage", mock.Anything).Return(nil, queue.ErrNoMessage).Once()
				return mockQueue
			},
			wantMsg: nil,
			wantErr: queue.ErrNoMessage,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			q := tc.setupMock(t)
			got, err := q.PopMessage(ctx)
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.wantErr)
				assert.Nil(t, got)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.wantMsg, got)
		})
	}
}
