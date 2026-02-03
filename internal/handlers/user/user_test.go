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
	service "github.com/luongtruong20201/bookmark-management/internal/services/user"
	"github.com/luongtruong20201/bookmark-management/internal/services/user/mocks"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_RegisterUser(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	var (
		testErrDatabase = errors.New("database error")
	)

	testCases := []struct {
		name           string
		requestBody    interface{}
		setupMockSvc   func(t *testing.T, ctx context.Context) *mocks.User
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "success",
			requestBody: createUserInputBody{
				Username:    "johndoe",
				Password:    "password123",
				DisplayName: "John Doe",
				Email:       "john.doe@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("CreateUser", ctx, "johndoe", "password123", "John Doe", "john.doe@example.com").
					Return(&model.User{
						Base: model.Base{
							ID: "550e8400-e29b-41d4-a716-446655440000",
						},
						Username:    "johndoe",
						DisplayName: "John Doe",
						Email:       "john.doe@example.com",
					}, nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: model.User{
				Base: model.Base{
					ID: "550e8400-e29b-41d4-a716-446655440000",
				},
				Username:    "johndoe",
				DisplayName: "John Doe",
				Email:       "john.doe@example.com",
			},
		},
		{
			name: "invalid request body - missing required fields",
			requestBody: createUserInputBody{
				Username: "",
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
			name: "invalid request body - invalid email format",
			requestBody: createUserInputBody{
				Username:    "johndoe",
				Password:    "password123",
				DisplayName: "John Doe",
				Email:       "invalid-email",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "invalid request body - password too short",
			requestBody: createUserInputBody{
				Username:    "johndoe",
				Password:    "12345",
				DisplayName: "John Doe",
				Email:       "john.doe@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "internal server error",
			requestBody: createUserInputBody{
				Username:    "johndoe",
				Password:    "password123",
				DisplayName: "John Doe",
				Email:       "john.doe@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("CreateUser", ctx, "johndoe", "password123", "John Doe", "john.doe@example.com").
					Return(nil, testErrDatabase).Once()
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
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
			name: "duplicate username or email",
			requestBody: createUserInputBody{
				Username:    "johndoe",
				Password:    "password123",
				DisplayName: "John Doe",
				Email:       "john.doe@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("CreateUser", ctx, "johndoe", "password123", "John Doe", "john.doe@example.com").
					Return(nil, dbutils.ErrDuplicationType).Once()
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "invalid request body - missing display name",
			requestBody: createUserInputBody{
				Username:    "johndoe",
				Password:    "password123",
				DisplayName: "",
				Email:       "john.doe@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "invalid request body - missing email",
			requestBody: createUserInputBody{
				Username:    "johndoe",
				Password:    "password123",
				DisplayName: "John Doe",
				Email:       "",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				return svcMock
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   nil,
		},
		{
			name: "invalid request body - password exactly 6 characters (boundary)",
			requestBody: createUserInputBody{
				Username:    "johndoe",
				Password:    "123456",
				DisplayName: "John Doe",
				Email:       "john.doe@example.com",
			},
			setupMockSvc: func(t *testing.T, ctx context.Context) *mocks.User {
				svcMock := mocks.NewUser(t)
				svcMock.On("CreateUser", ctx, "johndoe", "123456", "John Doe", "john.doe@example.com").
					Return(&model.User{
						Base: model.Base{
							ID: "550e8400-e29b-41d4-a716-446655440000",
						},
						Username:    "johndoe",
						DisplayName: "John Doe",
						Email:       "john.doe@example.com",
					}, nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusOK,
			expectedBody: model.User{
				Base: model.Base{
					ID: "550e8400-e29b-41d4-a716-446655440000",
				},
				Username:    "johndoe",
				DisplayName: "John Doe",
				Email:       "john.doe@example.com",
			},
		},
		{
			name:        "invalid request body - empty request body",
			requestBody: createUserInputBody{},
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

			ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer(reqBody))
			ctx.Request.Header.Set("Content-Type", "application/json")
			mockSvc := tc.setupMockSvc(t, ctx)
			testHandler := NewUser(mockSvc)

			testHandler.RegisterUser(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedBody != nil {
				var responseBody map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				assert.Equal(t, "Register an user successfully!", responseBody["message"])

				data, ok := responseBody["data"].(map[string]interface{})
				assert.True(t, ok, "data should be a map")

				expectedUser := tc.expectedBody.(model.User)
				assert.Equal(t, expectedUser.ID, data["id"])
				assert.Equal(t, expectedUser.Username, data["username"])
				assert.Equal(t, expectedUser.DisplayName, data["display_name"])
				assert.Equal(t, expectedUser.Email, data["email"])
			} else if tc.expectedStatus == http.StatusBadRequest {
				var responseBody map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)

				if tc.name == "duplicate username or email" {
					message, ok := responseBody["message"].(string)
					assert.True(t, ok)
					assert.Equal(t, "username or email already taken", message)
				} else {
					_, hasMessage := responseBody["message"]
					assert.True(t, hasMessage, "Bad request should have message field")
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
