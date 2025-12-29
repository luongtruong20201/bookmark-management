package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/luongtruong20201/bookmark-management/services/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHealthcheckHandler_Check(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupRequest   func(*gin.Context)
		setupMockSvc   func(*testing.T) *mocks.Healthcheck
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "success",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},
			setupMockSvc: func(t *testing.T) *mocks.Healthcheck {
				mockSvc := mocks.NewHealthcheck(t)
				mockSvc.On("Check").Return("OK", "bookmark_service", "instance_id")
				return mockSvc
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"instance_id":"instance_id","message":"OK","service_name":"bookmark_service"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ctx)
			svc := tc.setupMockSvc(t)
			handler := NewHealthcheck(svc)

			handler.Check(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, rec.Body.String(), tc.expectedBody)
		})
	}
}
