package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luongtruong20201/bookmark-management/internal/api"
	redisPkg "github.com/luongtruong20201/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestShortenURLEndpoint_ShortenURL(t *testing.T) {
	t.Parallel()

	cfg := &api.Config{
		AppPort:     "8080",
		ServiceName: "12345",
		InstanceId:  "12345",
	}

	testCases := []struct {
		name           string
		setupHTTP      func(api.Engine) *httptest.ResponseRecorder
		expectedStatus int
		verifyBody     func(t *testing.T, body map[string]any)
	}{
		{
			name: "success",
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"url": "https://truonglq.com",
					"exp": 3600,
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsBody))
				req.Header.Set("Context-Type", "application/json")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, body["message"], "OK")
				assert.Equal(t, len(body["code"].(string)), 7)
			},
		},
		{
			name: "invalid url",
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"url": "truonglq",
					"exp": 3600,
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, body["message"], "Input error")
			},
		},
		{
			name: "invalid exp",
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"url": "https://truonglq.com",
					"exp": "3600",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, body["message"], "Input Error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			redis := redisPkg.InitMockRedis(t)
			app := api.New(cfg, redis, nil, nil, nil)
			rec := tc.setupHTTP(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var body map[string]any
			_ = json.Unmarshal(rec.Body.Bytes(), &body)

			tc.verifyBody(t, body)
		})
	}
}

func TestShortenURLEndpoint_GetURL(t *testing.T) {
	t.Parallel()

	cfg := &api.Config{
		AppPort:     "8080",
		ServiceName: "12345",
		InstanceId:  "12345",
	}

	testCases := []struct {
		name           string
		setupHTTP      func(api.Engine) *httptest.ResponseRecorder
		setupMockRedis func(t *testing.T, ctx context.Context) *redis.Client
		expectedStatus int
		verifyBody     func(t *testing.T, body map[string]any)
		verifyRedirect func(t *testing.T, location string)
	}{
		{
			name: "not found",
			setupMockRedis: func(t *testing.T, ctx context.Context) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				return redis
			},
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/1234567", nil)
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, body["message"], "url not found")
			},
		},
		{
			name: "internal server error",
			setupMockRedis: func(t *testing.T, ctx context.Context) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				_ = redis.Close()
				return redis
			},
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/1234567", nil)
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			expectedStatus: http.StatusInternalServerError,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, body["message"], "internal server error")
			},
		},
		{
			name: "success",
			setupMockRedis: func(t *testing.T, ctx context.Context) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				redis.Set(ctx, "1234567", "https://truonglq.com", 0)
				return redis
			},
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/1234567", nil)
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			expectedStatus: http.StatusMovedPermanently,
			verifyRedirect: func(t *testing.T, location string) {
				assert.Equal(t, location, "https://truonglq.com")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			app := api.New(cfg, tc.setupMockRedis(t, ctx), nil, nil, nil)
			rec := tc.setupHTTP(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.verifyBody != nil {
				var body map[string]any
				_ = json.Unmarshal(rec.Body.Bytes(), &body)
				tc.verifyBody(t, body)
			}

			if tc.verifyRedirect != nil {
				location := rec.Header().Get("Location")
				tc.verifyRedirect(t, location)
			}
		})
	}
}
