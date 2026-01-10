package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luongtruong20201/bookmark-management/internal/api"
	redisPkg "github.com/luongtruong20201/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestHealthcheckEndPoint(t *testing.T) {
	t.Parallel()

	cfg := &api.Config{
		AppPort:     "8080",
		ServiceName: "12345",
		InstanceId:  "12345",
	}

	testCases := []struct {
		name            string
		setupTestHTTP   func(api.Engine) *httptest.ResponseRecorder
		expectedStatus  int
		setupMockRedis  func(t *testing.T) *redis.Client
		expectedMessage string
	}{
		{
			name: "success",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			setupMockRedis: func(t *testing.T) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				return redis
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "{\"instance_id\":\"12345\",\"message\":\"OK\",\"service_name\":\"12345\"}",
		},
		{
			name: "failed",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
				rec := httptest.NewRecorder()
				api.ServeHTTP(rec, req)
				return rec
			},
			setupMockRedis: func(t *testing.T) *redis.Client {
				redis := redisPkg.InitMockRedis(t)
				_ = redis.Close()
				return redis
			},
			expectedStatus:  http.StatusOK,
			expectedMessage: "{\"instance_id\":\"12345\",\"message\":\"NOT_OK\",\"service_name\":\"12345\"}",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := tc.setupTestHTTP(api.New(cfg, tc.setupMockRedis(t)))

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedMessage, rec.Body.String())
		})
	}
}
