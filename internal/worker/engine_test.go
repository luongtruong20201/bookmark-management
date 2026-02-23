package worker

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/luongtruong20201/bookmark-management/internal/repositories/queue"
	"github.com/luongtruong20201/bookmark-management/internal/worker/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewEngine(t *testing.T) {
	t.Parallel()

	mockQueue := mocks.NewQueue(t)
	mockHandler := mocks.NewHandler(t)

	eng := NewEngine(mockQueue, mockHandler, 2)

	require.NotNil(t, eng)
	assert.Implements(t, (*Engine)(nil), eng)
}

func TestEngine_Start(t *testing.T) {
	t.Parallel()

	popErrOther := errors.New("queue unavailable")

	testCases := []struct {
		name              string
		setupMocks        func(t *testing.T, ctx context.Context, mockQueue *mocks.Queue, mockHandler *mocks.Handler)
		cancelImmediately bool
		runBeforeCancel   time.Duration
		assertQueue       bool
		wantHandlerCalled bool
	}{
		{
			name:              "exits when context already cancelled",
			setupMocks:        func(t *testing.T, ctx context.Context, mockQueue *mocks.Queue, mockHandler *mocks.Handler) {},
			cancelImmediately: true,
			runBeforeCancel:   0,
			assertQueue:       false,
			wantHandlerCalled: false,
		},
		{
			name: "exits when context cancelled after some polls with ErrNoMessage",
			setupMocks: func(t *testing.T, ctx context.Context, mockQueue *mocks.Queue, mockHandler *mocks.Handler) {
				mockQueue.On("PopMessage", ctx).Return(nil, queue.ErrNoMessage)
			},
			cancelImmediately: false,
			runBeforeCancel:   50 * time.Millisecond,
			assertQueue:       true,
			wantHandlerCalled: false,
		},
		{
			name: "consumes message and calls handler",
			setupMocks: func(t *testing.T, ctx context.Context, mockQueue *mocks.Queue, mockHandler *mocks.Handler) {
				msg := []byte(`{"test": "data"}`)
				mockQueue.On("PopMessage", ctx).Return(msg, nil).Once()
				mockQueue.On("PopMessage", ctx).Return(nil, queue.ErrNoMessage)
				mockHandler.On("Handle", ctx, msg).Return(nil).Once()
			},
			cancelImmediately: false,
			runBeforeCancel:   200 * time.Millisecond,
			assertQueue:       true,
			wantHandlerCalled: true,
		},
		{
			name: "on pop error other than ErrNoMessage retries and does not call handler",
			setupMocks: func(t *testing.T, ctx context.Context, mockQueue *mocks.Queue, mockHandler *mocks.Handler) {
				mockQueue.On("PopMessage", ctx).Return(nil, popErrOther)
			},
			cancelImmediately: false,
			runBeforeCancel:   100 * time.Millisecond,
			assertQueue:       true,
			wantHandlerCalled: false,
		},
		{
			name: "on ErrNoMessage sleeps and does not call handler",
			setupMocks: func(t *testing.T, ctx context.Context, mockQueue *mocks.Queue, mockHandler *mocks.Handler) {
				mockQueue.On("PopMessage", ctx).Return(nil, queue.ErrNoMessage)
			},
			cancelImmediately: false,
			runBeforeCancel:   80 * time.Millisecond,
			assertQueue:       true,
			wantHandlerCalled: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx, cancel := context.WithCancel(t.Context())
			defer cancel()

			mockQueue := mocks.NewQueue(t)
			mockHandler := mocks.NewHandler(t)
			tc.setupMocks(t, ctx, mockQueue, mockHandler)

			eng := NewEngine(mockQueue, mockHandler, 1)

			if tc.cancelImmediately {
				cancel()
			}

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				eng.Start(ctx)
			}()

			if !tc.cancelImmediately {
				time.Sleep(tc.runBeforeCancel)
				cancel()
			}

			done := make(chan struct{})
			go func() {
				wg.Wait()
				close(done)
			}()

			select {
			case <-done:
			case <-time.After(3 * time.Second):
				t.Fatal("Start did not exit within timeout")
			}

			if tc.assertQueue {
				mockQueue.AssertExpectations(t)
			}
			if tc.wantHandlerCalled {
				mockHandler.AssertExpectations(t)
			} else {
				mockHandler.AssertNotCalled(t, "Handle", mock.Anything, mock.Anything)
			}
		})
	}
}
