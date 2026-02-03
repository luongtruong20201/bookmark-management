package user

import (
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

func TestUserHandler_GetProfile(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	const mockUserID = "550e8400-e29b-41d4-a716-446655440000"

	var (
		testErrDatabase = errors.New("database error")
	)

	testCases := []struct {
		name           string
		setupContext   func(c *gin.Context)
		setupMockSvc   func(t *testing.T, ctx context.Context, userID string) *mocks.User
		expectedStatus int
		expectedBody   *model.User
	}{
		{
			name: "success - valid user ID in context",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("GetUserByID", ctx, userID).
					Return(&model.User{
						Base: model.Base{
							ID: mockUserID,
						},
						Username:    "johndoe",
						DisplayName: "John Doe",
						Email:       "john.doe@example.com",
					}, nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: &model.User{
				Base: model.Base{
					ID: mockUserID,
				},
				Username:    "johndoe",
				DisplayName: "John Doe",
				Email:       "john.doe@example.com",
			},
		},
		{
			name:         "error - userID not in context",
			setupContext: func(c *gin.Context) {},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
		},
		{
			name: "error - userID is not a string",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": 12345})
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
		},
		{
			name: "error - service returns error",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("GetUserByID", ctx, userID).
					Return(nil, testErrDatabase).Once()
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
		{
			name: "error - user not found",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("GetUserByID", ctx, userID).
					Return(nil, dbutils.ErrNotFoundType).Once()
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
		{
			name: "error - empty string userID",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": ""})
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
		},
		{
			name: "error - userID is float64",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": 123.45})
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
		},
		{
			name: "error - userID is bool",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": true})
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
		},
		{
			name: "error - userID is slice",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": []string{"test"}})
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
		},
		{
			name: "error - claims is not jwt.MapClaims type",
			setupContext: func(c *gin.Context) {
				c.Set("claims", "invalid-claims")
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
		},
		{
			name: "error - claims missing sub field",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{})
			},
			setupMockSvc: func(t *testing.T, ctx context.Context, userID string) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)

			tc.setupContext(ctx)
			mockSvc := tc.setupMockSvc(t, ctx, mockUserID)
			testHandler := NewUser(mockSvc)

			testHandler.GetProfile(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedBody != nil {
				var responseBody model.User
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				assert.Equal(t, tc.expectedBody.ID, responseBody.ID)
				assert.Equal(t, tc.expectedBody.Username, responseBody.Username)
				assert.Equal(t, tc.expectedBody.DisplayName, responseBody.DisplayName)
				assert.Equal(t, tc.expectedBody.Email, responseBody.Email)
			} else if tc.expectedStatus == http.StatusUnauthorized {
				var responseBody map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				message, ok := responseBody["message"].(string)
				assert.True(t, ok)
				assert.Equal(t, "Invalid token", message)
			} else if tc.expectedStatus == http.StatusInternalServerError {
				var responseBody map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				message, ok := responseBody["message"].(string)
				assert.True(t, ok)
				assert.Equal(t, "Processing Error", message)
			}
		})
	}
}
