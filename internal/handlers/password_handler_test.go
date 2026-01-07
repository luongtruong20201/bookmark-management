package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/luongtruong20201/bookmark-management/internal/services/mocks"
	"github.com/stretchr/testify/assert"
)

func TestPasswordHandler_GenPass(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupRequest   func(*gin.Context)
		setupMockSvc   func() *mocks.Password
		expectedStatus int
		expectedResp   string
	}{
		{
			name: "success",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
			},
			setupMockSvc: func() *mocks.Password {
				svcMock := mocks.NewPassword(t)
				svcMock.On("GeneratePassword").Return("1234567890", nil)
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedResp:   "1234567890",
		},
		{
			name: "internal server error",
			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/gen-pass", nil)
			},
			setupMockSvc: func() *mocks.Password {
				svcMock := mocks.NewPassword(t)
				svcMock.On("GeneratePassword").Return("", errors.New("failed"))
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResp:   "err",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			gc, _ := gin.CreateTestContext(rec)
			tc.setupRequest(gc)
			mockSvc := tc.setupMockSvc()
			testHandler := NewPassword(mockSvc)

			testHandler.GenPass(gc)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResp, rec.Body.String())
		})
	}
}
