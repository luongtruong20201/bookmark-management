package endpoint

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/luongtruong20201/bookmark-management/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestHealthcheckEndPoint(t *testing.T) {
	t.Parallel()

	cfg := &api.Config{
		AppPort:     "8080",
		ServiceName: "bookmark-api",
		InstanceId:  uuid.New().String(),
	}

	testCases := []struct {
		name            string
		setupTestHTTP   func(api.Engine) *httptest.ResponseRecorder
		expectedStatus  int
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
			expectedStatus:  http.StatusOK,
			expectedMessage: "OK",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := tc.setupTestHTTP(api.New(cfg))

			assert.Equal(t, tc.expectedStatus, rec.Code)
			resp := map[string]any{}

			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, resp["message"], "OK")
		})
	}
}
