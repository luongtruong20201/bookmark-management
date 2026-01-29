package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	service "github.com/luongtruong20201/bookmark-management/internal/services"
	"github.com/luongtruong20201/bookmark-management/internal/services/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestShortenURLHandler_ShortenURL(t *testing.T) {
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
				assert.Equal(t, tc.expectedResp, actualResp)
			}
		})
	}
}

func TestShortenURLHandler_GetURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupRequest   func(*gin.Context)
		setupMockSvc   func(t *testing.T, ctx context.Context) *mocks.ShortenURL
		expectedStatus int
		expectedResp   map[string]any
	}{
		{
			name: "unprocessable",
			setupRequest: func(c *gin.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect", nil)
				c.Request = req
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenURL {
				svc := mocks.NewShortenURL(t)

				return svc
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]any{
				"message": "unprocessable",
			},
		},
		{
			name: "not found",
			setupRequest: func(c *gin.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect", nil)
				c.Request = req
				c.Params = gin.Params{gin.Param{Key: "code", Value: "1234567"}}
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenURL {
				svc := mocks.NewShortenURL(t)
				svc.On("GetURL", ctx, "1234567").Return("", service.ErrCodeNotFound)

				return svc
			},
			expectedStatus: http.StatusBadRequest,
			expectedResp: map[string]any{
				"message": "url not found",
			},
		},
		{
			name: "redis connection",
			setupRequest: func(c *gin.Context) {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect", nil)
				c.Request = req
				c.Params = gin.Params{gin.Param{Key: "code", Value: "1234567"}}
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenURL {
				svc := mocks.NewShortenURL(t)
				svc.On("GetURL", ctx, "1234567").Return("", redis.ErrClosed)

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
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect", nil)
				c.Request = req
				c.Params = gin.Params{gin.Param{Key: "code", Value: "1234567"}}
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.ShortenURL {
				svc := mocks.NewShortenURL(t)
				svc.On("GetURL", ctx, "1234567").Return("https://truonglq.com", nil)

				return svc
			},
			expectedStatus: http.StatusMovedPermanently,
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

			handler.GetURL(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedResp != nil {
				var actualResp map[string]any
				err := json.Unmarshal(rec.Body.Bytes(), &actualResp)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResp, actualResp)
			}
		})
	}
}
