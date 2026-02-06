package endpoint

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/luongtruong20201/bookmark-management/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestPasswordEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Parallel()

	cfg := &api.Config{
		AppPort:     "8080",
		ServiceName: "12345",
		InstanceId:  "12345",
	}

	testCases := []struct {
		name            string
		setupTestHTTP   func(api api.Engine) *httptest.ResponseRecorder
		expectedStatus  int
		expectedRespLen int
	}{
		{
			name: "success",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
				respRec := httptest.NewRecorder()
				api.ServeHTTP(respRec, req)
				return respRec
			},
			expectedStatus:  http.StatusOK,
			expectedRespLen: 10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			app := api.New(&api.EngineOpts{
				Engine: gin.New(),
				Cfg:    cfg,
			})
			rec := tc.setupTestHTTP(app)

			assert.Equal(t, tc.expectedStatus, http.StatusOK)
			assert.Equal(t, tc.expectedRespLen, len(rec.Body.Bytes()))
		})
	}
}
