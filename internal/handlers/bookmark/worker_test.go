package bookmark

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/luongtruong20201/bookmark-management/internal/services/bookmark/mocks"
	"github.com/luongtruong20201/bookmark-management/internal/services/queue"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errMock = errors.New("mock service error")

func TestWorker_Handle(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		message     []byte
		setupMock   func(t *testing.T, ctx context.Context) *mocks.Service
		expectedErr bool
	}{
		{
			name: "success_unmarshal_and_calls_ImportBookmarks",
			message: func() []byte {
				msg := queue.ImportMessage{
					UID: "user-123",
					Bookmarks: []*queue.ImportBookmarkInput{
						{Description: "Desc 1", URL: "https://a.com"},
						{Description: "Desc 2", URL: "https://b.com"},
					},
				}
				b, _ := json.Marshal(msg)
				return b
			}(),
			setupMock: func(t *testing.T, ctx context.Context) *mocks.Service {
				svc := mocks.NewService(t)
				svc.On("ImportBookmarks", ctx, "user-123", mock.Anything).Return(nil).Once()
				return svc
			},
			expectedErr: false,
		},
		{
			name:        "invalid_JSON_returns_error",
			message:     []byte(`{invalid json`),
			setupMock:   func(t *testing.T, ctx context.Context) *mocks.Service { return mocks.NewService(t) },
			expectedErr: true,
		},
		{
			name: "ImportBookmarks_error_returns_error",
			message: func() []byte {
				msg := queue.ImportMessage{UID: "user-1", Bookmarks: []*queue.ImportBookmarkInput{{Description: "D", URL: "https://x.com"}}}
				b, _ := json.Marshal(msg)
				return b
			}(),
			setupMock: func(t *testing.T, ctx context.Context) *mocks.Service {
				svc := mocks.NewService(t)
				svc.On("ImportBookmarks", ctx, "user-1", mock.Anything).Return(errMock).Once()
				return svc
			},
			expectedErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockSvc := tc.setupMock(t, ctx)
			h := NewWorkerHandler(mockSvc)

			err := h.Handle(ctx, tc.message)

			if tc.expectedErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
