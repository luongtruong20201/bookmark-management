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
	service "github.com/luongtruong20201/bookmark-management/internal/services/user"
	"github.com/luongtruong20201/bookmark-management/internal/services/user/mocks"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_Login(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	const mockToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI1NTBlODQwMC1lMjliLTQxZDQtYTcxNi00NDY2NTU0NDAwMDAiLCJpYXQiOjE2MDAwMDAwMDAsImV4cCI6MTYwMDA4NjQwMH0.test"

	var (
		testErrDatabase = errors.New("database error")
		testErrJWT      = errors.New("jwt generation error")
	)

	testCases := []struct {
		name           string
		requestBody    interface{}
		setupMockSvc   func(t *testing.T, ctx context.Context) *mocks.User
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "success - valid credentials",
			requestBody: loginRequestBody{
				Username: "johndoe",
				Password: "password123",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("Login", ctx, "johndoe", "password123").
					Return(mockToken, nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: loginResponseBody{
				Token: mockToken,
			},
		},
		{
			name: "invalid request body - missing username",
			requestBody: loginRequestBody{
				Password: "password123",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "invalid request body - missing password",
			requestBody: loginRequestBody{
				Username: "johndoe",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:        "invalid JSON",
			requestBody: "invalid json string",
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "error - invalid password",
			requestBody: loginRequestBody{
				Username: "johndoe",
				Password: "wrongpassword",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("Login", ctx, "johndoe", "wrongpassword").
					Return("", service.ErrClientErr).Once()
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "error - user not found",
			requestBody: loginRequestBody{
				Username: "nonexistent",
				Password: "password123",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("Login", ctx, "nonexistent", "password123").
					Return("", dbutils.ErrNotFoundType).Once()
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "error - internal server error",
			requestBody: loginRequestBody{
				Username: "johndoe",
				Password: "password123",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("Login", ctx, "johndoe", "password123").
					Return("", testErrDatabase).Once()
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
		{
			name: "error - JWT generation error",
			requestBody: loginRequestBody{
				Username: "johndoe",
				Password: "password123",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("Login", ctx, "johndoe", "password123").
					Return("", testErrJWT).Once()
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
		{
			name: "invalid request body - empty username and password",
			requestBody: loginRequestBody{
				Username: "",
				Password: "",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:        "invalid request body - empty JSON object",
			requestBody: "{}",
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name:        "invalid request body - null JSON",
			requestBody: "null",
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
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

			ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewBuffer(reqBody))
			ctx.Request.Header.Set("Content-Type", "application/json")
			mockSvc := tc.setupMockSvc(t, ctx)
			testHandler := NewUser(mockSvc)

			testHandler.Login(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedBody != nil {
				var responseBody loginResponseBody
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				expectedResponse := tc.expectedBody.(loginResponseBody)
				assert.Equal(t, expectedResponse.Token, responseBody.Token)
			} else if tc.expectedStatus == http.StatusBadRequest {
				var responseBody map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				if tc.name == "error - invalid password" {
					errorMsg, ok := responseBody["error"].(string)
					assert.True(t, ok)
					assert.Contains(t, errorMsg, "invalid username or password")
				} else if tc.name == "error - user not found" {
					message, ok := responseBody["message"].(string)
					assert.True(t, ok)
					assert.Equal(t, "invalid username or password", message)
				} else {
					_, hasMessage := responseBody["message"]
					_, hasDetails := responseBody["details"]
					assert.True(t, hasMessage || hasDetails, "Bad request should have message or details field")
				}
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
