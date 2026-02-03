package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/internal/services/user/mocks"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_UpdateProfile(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	const mockUserID = "550e8400-e29b-41d4-a716-446655440000"

	var (
		testErrDatabase = errors.New("database error")
	)

	testCases := []struct {
		name           string
		setupContext   func(c *gin.Context)
		requestBody    interface{}
		setupMockSvc   func(t *testing.T, ctx context.Context, userID string) *mocks.User
		expectedStatus int
	}{
		{
			name: "success - valid user ID and body",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			requestBody: updateProfileRequestBody{
				DisplayName: "John Updated",
				Email:       "john.updated@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("UpdateUserProfile", ctx, userID, "John Updated", "john.updated@example.com").
					Return(&model.User{
						Base: model.Base{
							ID: mockUserID,
						},
						Username:    "johndoe",
						DisplayName: "John Updated",
						Email:       "john.updated@example.com",
					}, nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:         "error - userID not in context",
			setupContext: func(c *gin.Context) {},
			requestBody: updateProfileRequestBody{
				DisplayName: "John Updated",
				Email:       "john.updated@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "error - userID is not a string",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": 12345})
			},
			requestBody: updateProfileRequestBody{
				DisplayName: "John Updated",
				Email:       "john.updated@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "error - invalid request body",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			requestBody: updateProfileRequestBody{
				DisplayName: "",
				Email:       "invalid-email",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "error - service returns duplication error",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			requestBody: updateProfileRequestBody{
				DisplayName: "John Updated",
				Email:       "existing@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("UpdateUserProfile", ctx, userID, "John Updated", "existing@example.com").
					Return(nil, dbutils.ErrDuplicationType).Once()
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "error - service returns internal error",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			requestBody: updateProfileRequestBody{
				DisplayName: "John Updated",
				Email:       "john.updated@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("UpdateUserProfile", ctx, userID, "John Updated", "john.updated@example.com").
					Return(nil, testErrDatabase).Once()
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "error - claims is not jwt.MapClaims type",
			setupContext: func(c *gin.Context) {
				c.Set("claims", "invalid-claims")
			},
			requestBody: updateProfileRequestBody{
				DisplayName: "John Updated",
				Email:       "john.updated@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "error - claims missing sub field",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{})
			},
			requestBody: updateProfileRequestBody{
				DisplayName: "John Updated",
				Email:       "john.updated@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			var reqBody []byte
			var err error
			if str, ok := tc.requestBody.(string); ok {
				reqBody = []byte(str)
			} else {
				reqBody, err = json.Marshal(tc.requestBody)
				assert.NoError(t, err)
			}

			ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", bytes.NewBuffer(reqBody))
			ctx.Request.Header.Set("Content-Type", "application/json")

			tc.setupContext(ctx)
			mockSvc := tc.setupMockSvc(t, ctx, mockUserID)
			testHandler := NewUser(mockSvc)

			testHandler.UpdateProfile(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
		})
	}
}
