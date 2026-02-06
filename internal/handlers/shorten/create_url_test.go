package shorten

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/luongtruong20201/bookmark-management/internal/services/shorten/mocks"
	"github.com/stretchr/testify/assert"
)

func TestShortenURLHandler_ShortenURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Parallel()

	testCases := []struct {
		name           string
		setupRequest   func(*gin.Context)
		setupMockSvc   func(t *testing.T, ctx context.Context) *mocks.ShortenURL
		expectedStatus int
		expectedResp   map[string]any
	}{
		{
			name: "fail validation",
			setupRequest: func(c *gin.Context) {
				body := map[string]any{
					"url": "https://truonglq.com",
					"exp": "123",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				c.Request = req
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenURL {
				svc := mocks.NewShortenURL(t)

				return svc
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]any{
				"message": "Input Error",
			},
		},
		{
			name: "internal server error",
			setupRequest: func(c *gin.Context) {
				body := map[string]any{
					"url": "https://truonglq.com",
					"exp": 123,
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				c.Request = req
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenURL {
				svc := mocks.NewShortenURL(t)
				svc.On("ShortenURL", ctx, "https://truonglq.com", 123).Return("", errors.New("failed"))
				return svc
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp: map[string]any{
				"message": "internal server error",
			},
		},
		{
			name: "success",
			setupRequest: func(c *gin.Context) {
				body := map[string]any{
					"url": "https://truonglq.com",
					"exp": 123,
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				c.Request = req
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenURL {
				svc := mocks.NewShortenURL(t)
				svc.On("ShortenURL", ctx, "https://truonglq.com", 123).Return("1234567", nil).Once()

				return svc
			},
			expectedStatus: http.StatusOK,
			expectedResp: map[string]any{
				"message": "OK",
				"code":    "1234567",
			},
		},
		{
			name: "missing URL field",
			setupRequest: func(c *gin.Context) {
				body := map[string]any{
					"exp": 123,
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				c.Request = req
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenURL {
				svc := mocks.NewShortenURL(t)
				return svc
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]any{
				"message": "Input error",
			},
		},
		{
			name: "missing exp field",
			setupRequest: func(c *gin.Context) {
				body := map[string]any{
					"url": "https://truonglq.com",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				c.Request = req
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenURL {
				svc := mocks.NewShortenURL(t)
				return svc
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]any{
				"message": "Input error",
			},
		},
		{
			name: "invalid URL format",
			setupRequest: func(c *gin.Context) {
				body := map[string]any{
					"url": "not-a-valid-url",
					"exp": 123,
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				c.Request = req
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenURL {
				svc := mocks.NewShortenURL(t)
				return svc
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]any{
				"message": "Input error",
			},
		},
		{
			name: "negative expiration",
			setupRequest: func(c *gin.Context) {
				body := map[string]any{
					"url": "https://truonglq.com",
					"exp": -1,
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				c.Request = req
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenURL {
				svc := mocks.NewShortenURL(t)
				return svc
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]any{
				"message": "Input error",
			},
		},
		{
			name: "expiration exceeds maximum",
			setupRequest: func(c *gin.Context) {
				body := map[string]any{
					"url": "https://truonglq.com",
					"exp": 604801,
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				c.Request = req
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenURL {
				svc := mocks.NewShortenURL(t)
				return svc
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]any{
				"message": "Input error",
			},
		},
		{
			name: "success with maximum expiration",
			setupRequest: func(c *gin.Context) {
				body := map[string]any{
					"url": "https://truonglq.com",
					"exp": 604800,
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				c.Request = req
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenURL {
				svc := mocks.NewShortenURL(t)
				svc.On("ShortenURL", ctx, "https://truonglq.com", 604800).Return("maxexp01", nil).Once()

				return svc
			},
			expectedStatus: http.StatusOK,
			expectedResp: map[string]any{
				"message": "OK",
				"code":    "maxexp01",
			},
		},
		{
			name: "duplicate key error",
			setupRequest: func(c *gin.Context) {
				body := map[string]any{
					"url": "https://truonglq.com",
					"exp": 123,
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				c.Request = req
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenURL {
				svc := mocks.NewShortenURL(t)
				svc.On("ShortenURL", ctx, "https://truonglq.com", 123).Return("", errors.New("duplicate key")).Once()

				return svc
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp: map[string]any{
				"message": "internal server error",
			},
		},
		{
			name: "empty request body",
			setupRequest: func(c *gin.Context) {
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader([]byte("")))
				req.Header.Set("Content-Type", "application/json")
				c.Request = req
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenURL {
				svc := mocks.NewShortenURL(t)
				return svc
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]any{
				"message": "Input Error",
			},
		},
		{
			name: "invalid JSON",
			setupRequest: func(c *gin.Context) {
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader([]byte("{invalid json}")))
				req.Header.Set("Content-Type", "application/json")
				c.Request = req
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenURL {
				svc := mocks.NewShortenURL(t)
				return svc
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]any{
				"message": "Input Error",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)
			svc := tc.setupMockSvc(t, ctx)
			handler := NewShortenURL(svc)

			handler.ShortenURL(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedResp != nil {
				var actualResp map[string]any
				err := json.Unmarshal(rec.Body.Bytes(), &actualResp)
				assert.NoError(t, err)
				if expectedMsg, ok := tc.expectedResp["message"]; ok {
					assert.Equal(t, expectedMsg, actualResp["message"], "message should match")
					if expectedCode, ok := tc.expectedResp["code"]; ok {
						assert.Equal(t, expectedCode, actualResp["code"], "code should match")
					}
				} else {
					assert.Equal(t, tc.expectedResp, actualResp)
				}
			}
		})
	}
}
