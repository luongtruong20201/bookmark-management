package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/luongtruong20201/bookmark-management/pkg/jwt/mocks"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuth_JWTAuth(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	const (
		mockUserID = "550e8400-e29b-41d4-a716-446655440000"
		mockToken  = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.test"
	)

	var (
		testErrValidation = errors.New("invalid token")
	)

	testCases := []struct {
		name           string
		authHeader     string
		setupMock      func(t *testing.T) *mocks.JWTValidator
		expectedStatus int
		expectedBody   map[string]interface{}
		expectedUserID interface{}
		shouldAbort    bool
	}{
		{
			name:       "success - valid Bearer token with sub claim",
			authHeader: "Bearer " + mockToken,
			setupMock: func(t *testing.T) *mocks.JWTValidator {
				mockValidator := mocks.NewJWTValidator(t)
				mockValidator.On("ValidateToken", mockToken).
					Return(jwt.MapClaims{
						"sub": mockUserID,
						"iat": 1600000000,
						"exp": 1600086400,
					}, nil).Once()
				return mockValidator
			},
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
			expectedUserID: mockUserID,
			shouldAbort:    false,
		},
		{
			name:       "error - missing Authorization header",
			authHeader: "",
			setupMock: func(t *testing.T) *mocks.JWTValidator {
				mockValidator := mocks.NewJWTValidator(t)
				return mockValidator
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Authorization header is required",
			},
			expectedUserID: nil,
			shouldAbort:    true,
		},
		{
			name:       "error - Authorization header without Bearer prefix",
			authHeader: mockToken,
			setupMock: func(t *testing.T) *mocks.JWTValidator {
				mockValidator := mocks.NewJWTValidator(t)
				return mockValidator
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Authorization header format is wrong",
			},
			expectedUserID: nil,
			shouldAbort:    true,
		},
		{
			name:       "error - Authorization header with wrong prefix",
			authHeader: "Basic " + mockToken,
			setupMock: func(t *testing.T) *mocks.JWTValidator {
				mockValidator := mocks.NewJWTValidator(t)
				return mockValidator
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Authorization header format is wrong",
			},
			expectedUserID: nil,
			shouldAbort:    true,
		},
		{
			name:       "error - Authorization header with empty token",
			authHeader: "Bearer ",
			setupMock: func(t *testing.T) *mocks.JWTValidator {
				mockValidator := mocks.NewJWTValidator(t)
				mockValidator.On("ValidateToken", "").
					Return(nil, testErrValidation).Once()
				return mockValidator
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid token",
			},
			expectedUserID: nil,
			shouldAbort:    true,
		},
		{
			name:       "error - Authorization header with multiple spaces",
			authHeader: "Bearer  " + mockToken,
			setupMock: func(t *testing.T) *mocks.JWTValidator {
				mockValidator := mocks.NewJWTValidator(t)
				return mockValidator
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Authorization header format is wrong",
			},
			expectedUserID: nil,
			shouldAbort:    true,
		},
		{
			name:       "error - Authorization header with token containing spaces",
			authHeader: "Bearer " + mockToken + " extra",
			setupMock: func(t *testing.T) *mocks.JWTValidator {
				mockValidator := mocks.NewJWTValidator(t)
				return mockValidator
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Authorization header format is wrong",
			},
			expectedUserID: nil,
			shouldAbort:    true,
		},
		{
			name:       "error - invalid token (validator returns error)",
			authHeader: "Bearer " + mockToken,
			setupMock: func(t *testing.T) *mocks.JWTValidator {
				mockValidator := mocks.NewJWTValidator(t)
				mockValidator.On("ValidateToken", mockToken).
					Return(nil, testErrValidation).Once()
				return mockValidator
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid token",
			},
			expectedUserID: nil,
			shouldAbort:    true,
		},
		{
			name:       "error - valid token but missing sub claim",
			authHeader: "Bearer " + mockToken,
			setupMock: func(t *testing.T) *mocks.JWTValidator {
				mockValidator := mocks.NewJWTValidator(t)
				mockValidator.On("ValidateToken", mockToken).
					Return(jwt.MapClaims{
						"iat": 1600000000,
						"exp": 1600086400,
					}, nil).Once()
				return mockValidator
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Invalid token",
			},
			expectedUserID: nil,
			shouldAbort:    true,
		},
		{
			name:       "error - valid token but sub claim is empty string",
			authHeader: "Bearer " + mockToken,
			setupMock: func(t *testing.T) *mocks.JWTValidator {
				mockValidator := mocks.NewJWTValidator(t)
				mockValidator.On("ValidateToken", mockToken).
					Return(jwt.MapClaims{
						"sub": "",
						"iat": 1600000000,
						"exp": 1600086400,
					}, nil).Once()
				return mockValidator
			},
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
			expectedUserID: "",
			shouldAbort:    false,
		},
		{
			name:       "success - valid token with sub claim as number",
			authHeader: "Bearer " + mockToken,
			setupMock: func(t *testing.T) *mocks.JWTValidator {
				mockValidator := mocks.NewJWTValidator(t)
				mockValidator.On("ValidateToken", mockToken).
					Return(jwt.MapClaims{
						"sub": 12345,
						"iat": 1600000000,
						"exp": 1600086400,
					}, nil).Once()
				return mockValidator
			},
			expectedStatus: http.StatusOK,
			expectedBody:   nil,
			expectedUserID: 12345,
			shouldAbort:    false,
		},
		{
			name:       "error - Authorization header with only Bearer",
			authHeader: "Bearer",
			setupMock: func(t *testing.T) *mocks.JWTValidator {
				mockValidator := mocks.NewJWTValidator(t)
				return mockValidator
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Authorization header format is wrong",
			},
			expectedUserID: nil,
			shouldAbort:    true,
		},
		{
			name:       "error - Authorization header with Bearer and multiple parts",
			authHeader: "Bearer token1 token2 token3",
			setupMock: func(t *testing.T) *mocks.JWTValidator {
				mockValidator := mocks.NewJWTValidator(t)
				return mockValidator
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "Authorization header format is wrong",
			},
			expectedUserID: nil,
			shouldAbort:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			_, engine := gin.CreateTestContext(rec)

			nextCalled := false
			testHandler := func(c *gin.Context) {
				nextCalled = true
				c.JSON(http.StatusOK, gin.H{"status": "ok"})
			}

			mockValidator := tc.setupMock(t)
			middleware := NewJWTAuth(mockValidator)
			engine.Use(middleware.JWTAuth())
			engine.GET("/test", testHandler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tc.authHeader != "" {
				req.Header.Set("Authorization", tc.authHeader)
			}

			engine.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			if tc.expectedBody != nil {
				var responseBody map[string]interface{}
				err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBody, responseBody)
			}

			if tc.shouldAbort {
				assert.False(t, nextCalled, "Next() should not be called when request is aborted")
			} else {
				assert.True(t, nextCalled, "Next() should be called when request is not aborted")
			}
		})
	}
}

func TestJWTAuth_JWTAuth_Integration(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	const (
		mockUserID = "550e8400-e29b-41d4-a716-446655440000"
		mockToken  = "valid.jwt.token"
	)

	t.Run("success - middleware chain with multiple handlers", func(t *testing.T) {
		t.Parallel()

		rec := httptest.NewRecorder()
		_, engine := gin.CreateTestContext(rec)

		mockValidator := mocks.NewJWTValidator(t)
		mockValidator.On("ValidateToken", mockToken).
			Return(jwt.MapClaims{
				"sub": mockUserID,
			}, nil).Once()

		middleware := NewJWTAuth(mockValidator)
		engine.Use(middleware.JWTAuth())

		handler1Called := false
		handler2Called := false

		engine.GET("/test",
			func(c *gin.Context) {
				handler1Called = true
				c.Next()
			},
			func(c *gin.Context) {
				handler2Called = true
				userID, _ := c.Get("userID")
				c.JSON(http.StatusOK, gin.H{
					"userID": userID,
					"status": "ok",
				})
			},
		)

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("Authorization", "Bearer "+mockToken)

		engine.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, handler1Called)
		assert.True(t, handler2Called)

		var responseBody map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, mockUserID, responseBody["userID"])
	})

	t.Run("error - middleware aborts before handlers", func(t *testing.T) {
		t.Parallel()

		rec := httptest.NewRecorder()
		_, engine := gin.CreateTestContext(rec)

		mockValidator := mocks.NewJWTValidator(t)
		middleware := NewJWTAuth(mockValidator)
		engine.Use(middleware.JWTAuth())

		handlerCalled := false
		engine.GET("/test", func(c *gin.Context) {
			handlerCalled = true
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		engine.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.False(t, handlerCalled, "Handler should not be called when middleware aborts")

		var responseBody map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
		assert.NoError(t, err)
		assert.Equal(t, "Authorization header is required", responseBody["error"])
	})
}
