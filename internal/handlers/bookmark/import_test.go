package bookmark

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	bookmarkServiceMocks "github.com/luongtruong20201/bookmark-management/internal/services/bookmark/mocks"
	queueService "github.com/luongtruong20201/bookmark-management/internal/services/queue"
	queueMocks "github.com/luongtruong20201/bookmark-management/internal/services/queue/mocks"
	"github.com/luongtruong20201/bookmark-management/pkg/response"
	"github.com/stretchr/testify/assert"
)

func buildMultipartRequest(t *testing.T, fieldName, filename, content string) (*http.Request, *httptest.ResponseRecorder) {
	t.Helper()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(fieldName, filename)
	assert.NoError(t, err)

	_, err = part.Write([]byte(content))
	assert.NoError(t, err)

	assert.NoError(t, writer.Close())

	req := httptest.NewRequest(http.MethodPost, "/v1/bookmarks/import", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rec := httptest.NewRecorder()

	return req, rec
}

func TestBookmarkHandler_ImportBookmarks(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	const uid = "550e8400-e29b-41d4-a716-446655440000"

	successBookmarks := []*queueService.ImportBookmarkInput{{Description: "My blog", URL: "https://truonglq.com"}}
	nilBookmarks := []*queueService.ImportBookmarkInput(nil)
	tests := []struct {
		name                  string
		buildRequest          func(t *testing.T) (*http.Request, *httptest.ResponseRecorder)
		expectSendBookmarkJob *[]*queueService.ImportBookmarkInput
		expectedStatus        int
		verifyResponse        func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "success",
			buildRequest: func(t *testing.T) (*http.Request, *httptest.ResponseRecorder) {
				return buildMultipartRequest(t, "file", "bookmarks.csv", "description,url\nMy blog,https://truonglq.com\n")
			},
			expectSendBookmarkJob: &successBookmarks,
			expectedStatus:        http.StatusOK,
			verifyResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp response.Message
				assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Equal(t, "Processing import", resp.Message)
			},
		},
		{
			name: "missing file",
			buildRequest: func(t *testing.T) (*http.Request, *httptest.ResponseRecorder) {
				rec := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodPost, "/v1/bookmarks/import", nil)
				return req, rec
			},
			expectSendBookmarkJob: nil,
			expectedStatus:        http.StatusBadRequest,
			verifyResponse:        nil,
		},
		{
			name: "invalid CSV - parser may still yield data, SendBookmarkJob called",
			buildRequest: func(t *testing.T) (*http.Request, *httptest.ResponseRecorder) {
				return buildMultipartRequest(t, "file", "bookmarks.csv", "not valid csv content\n")
			},
			expectSendBookmarkJob: &nilBookmarks,
			expectedStatus:        http.StatusOK,
			verifyResponse:        nil,
		},
		{
			name: "validation error - missing required url",
			buildRequest: func(t *testing.T) (*http.Request, *httptest.ResponseRecorder) {
				return buildMultipartRequest(t, "file", "bookmarks.csv", "description,url\nMy blog,\n")
			},
			expectSendBookmarkJob: nil,
			expectedStatus:        http.StatusBadRequest,
			verifyResponse:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req, rec := tt.buildRequest(t)
			req = req.WithContext(t.Context())
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = req
			ctx.Set("claims", jwt.MapClaims{"sub": uid})

			queueSvc := queueMocks.NewService(t)
			if tt.expectSendBookmarkJob != nil {
				queueSvc.On("SendBookmarkJob", t.Context(), uid, *tt.expectSendBookmarkJob).Return(nil).Once()
			}

			bookmarkSvc := bookmarkServiceMocks.NewService(t)
			h := NewBookmarkHandler(bookmarkSvc, queueSvc)

			h.ImportBookmarks(ctx)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.verifyResponse != nil {
				tt.verifyResponse(t, rec)
			}
		})
	}
}
