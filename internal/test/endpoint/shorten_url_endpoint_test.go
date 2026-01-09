package endpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luongtruong20201/bookmark-management/internal/api"
	"github.com/luongtruong20201/bookmark-management/pkg/redis"
	"github.com/stretchr/testify/assert"
)

func TestShortenURLEndpoint(t *testing.T) {
	t.Parallel()

	cfg := &api.Config{
		AppPort:     "8080",
		ServiceName: "12345",
		InstanceId:  "12345",
	}

	redis := redis.InitMockRedis(t)

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
				req.Header.Set("Content-Type", "application/json")
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
				assert.Equal(t, body["message"], "unprocessable")
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
				assert.Equal(t, body["message"], "unprocessable")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := api.New(cfg, redis)
			rec := tc.setupHTTP(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var body map[string]any
			_ = json.Unmarshal(rec.Body.Bytes(), &body)

			tc.verifyBody(t, body)
		})
	}
}
