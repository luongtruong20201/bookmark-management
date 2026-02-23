package queue

import (
	"encoding/json"
	"errors"
	"testing"

	repoMocks "github.com/luongtruong20201/bookmark-management/internal/repositories/queue/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSendBookmarkJob(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		uid             string
		bookmarks       []*ImportBookmarkInput
		expectedPushMsg [][]byte
		pushErr         error
		expectedError   string
	}{
		{
			name:      "single chunk - one PushMessage call",
			uid:       "user-1",
			bookmarks: []*ImportBookmarkInput{{Description: "d1", URL: "https://1.com"}, {Description: "d2", URL: "https://2.com"}, {Description: "d3", URL: "https://3.com"}},
			expectedPushMsg: func() [][]byte {
				b, _ := json.Marshal(&ImportMessage{UID: "user-1", Bookmarks: []*ImportBookmarkInput{
					{Description: "d1", URL: "https://1.com"},
					{Description: "d2", URL: "https://2.com"},
					{Description: "d3", URL: "https://3.com"},
				}})
				return [][]byte{b}
			}(),
			expectedError: "",
		},
		{
			name: "multiple chunks - two PushMessage calls",
			uid:  "user-1",
			bookmarks: func() []*ImportBookmarkInput {
				bookmarks := make([]*ImportBookmarkInput, 101)
				for i := 0; i < 101; i++ {
					bookmarks[i] = &ImportBookmarkInput{URL: "https://example.com"}
				}
				return bookmarks
			}(),
			expectedPushMsg: func() [][]byte {
				chunk1 := make([]*ImportBookmarkInput, 100)
				for i := 0; i < 100; i++ {
					chunk1[i] = &ImportBookmarkInput{URL: "https://example.com"}
				}
				chunk2 := []*ImportBookmarkInput{{URL: "https://example.com"}}
				b1, _ := json.Marshal(&ImportMessage{UID: "user-1", Bookmarks: chunk1})
				b2, _ := json.Marshal(&ImportMessage{UID: "user-1", Bookmarks: chunk2})
				return [][]byte{b1, b2}
			}(),
			expectedError: "",
		},
		{
			name:      "error on PushMessage propagates",
			uid:       "user-1",
			bookmarks: []*ImportBookmarkInput{{URL: "https://1.com"}},
			expectedPushMsg: func() [][]byte {
				b, _ := json.Marshal(&ImportMessage{UID: "user-1", Bookmarks: []*ImportBookmarkInput{{URL: "https://1.com"}}})
				return [][]byte{b}
			}(),
			pushErr:       errors.New("queue error"),
			expectedError: "queue error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			mockRepo := repoMocks.NewRepository(t)
			for _, msg := range tc.expectedPushMsg {
				if tc.pushErr != nil {
					mockRepo.On("PushMessage", ctx, msg).Return(tc.pushErr).Once()
					break
				}
				mockRepo.On("PushMessage", ctx, msg).Return(nil).Once()
			}

			svc := NewService(mockRepo)
			err := svc.SendBookmarkJob(ctx, tc.uid, tc.bookmarks)

			if tc.expectedError != "" {
				require.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
				return
			}
			require.NoError(t, err)
		})
	}
}
