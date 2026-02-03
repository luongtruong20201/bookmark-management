package endpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/luongtruong20201/bookmark-management/internal/api"
	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/internal/test/fixture"
	jwtPkg "github.com/luongtruong20201/bookmark-management/pkg/jwt"
	jwtMocks "github.com/luongtruong20201/bookmark-management/pkg/jwt/mocks"
	sqldb "github.com/luongtruong20201/bookmark-management/pkg/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestUserEndpoint_RegisterUser(t *testing.T) {
	t.Parallel()

	cfg := &api.Config{
		AppPort:     "8080",
		ServiceName: "bookmark-service",
		InstanceId:  "instance-1",
	}

	testCases := []struct {
		name           string
		setupHTTP      func(api.Engine) *httptest.ResponseRecorder
		expectedStatus int
		verifyBody     func(t *testing.T, body map[string]any)
		verifyUser     func(t *testing.T, db *gorm.DB)
	}{
		{
			name: "success",
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"username":     "johndoe",
					"password":     "password123",
					"display_name": "John Doe",
					"email":        "john.doe@example.com",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, body["message"], "Register an user successfully!")
				data, ok := body["data"].(map[string]any)
				assert.True(t, ok, "data should be a map")
				assert.NotEmpty(t, data["id"])
				assert.Equal(t, data["username"], "johndoe")
				assert.Equal(t, data["display_name"], "John Doe")
				assert.Equal(t, data["email"], "john.doe@example.com")
				_, exists := data["password"]
				assert.False(t, exists)
			},
			verifyUser: func(t *testing.T, db *gorm.DB) {
				var user model.User
				err := db.Where("username = ?", "johndoe").First(&user).Error
				assert.NoError(t, err)
				assert.Equal(t, "johndoe", user.Username)
				assert.Equal(t, "John Doe", user.DisplayName)
				assert.Equal(t, "john.doe@example.com", user.Email)
				assert.NotEmpty(t, user.ID)
				assert.NotEqual(t, "password123", user.Password)
				assert.Contains(t, user.Password, "$2a$")
			},
		},
		{
			name: "invalid request body - missing required fields",
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"username": "johndoe",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, body["message"], "Input error")
			},
			verifyUser: nil,
		},
		{
			name: "invalid request body - invalid email format",
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"username":     "johndoe",
					"password":     "password123",
					"display_name": "John Doe",
					"email":        "invalid-email",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, body["message"], "Input error")
			},
			verifyUser: nil,
		},
		{
			name: "invalid request body - password too short",
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"username":     "johndoe",
					"password":     "12345",
					"display_name": "John Doe",
					"email":        "john.doe@example.com",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, body["message"], "Input error")
			},
			verifyUser: nil,
		},
		{
			name: "invalid JSON",
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewReader([]byte("invalid json")))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, body["message"], "Input Error")
			},
			verifyUser: nil,
		},
		{
			name: "duplicate username",
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"username":     "existinguser",
					"password":     "password123",
					"display_name": "Existing User",
					"email":        "existinguser@example.com",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			expectedStatus: http.StatusInternalServerError,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, body["message"], "Processing Error")
			},
			verifyUser: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := sqldb.InitMockDB(t)
			if err := db.AutoMigrate(&model.User{}); err != nil {
				t.Fatalf("failed to migrate user schema: %v", err)
			}
			if tc.name == "duplicate username" {
				existingUser := &model.User{
					Base: model.Base{
						ID: "550e8400-e29b-41d4-a716-446655440000",
					},
					Username:    "existinguser",
					Password:    "hashedpassword",
					DisplayName: "Existing User",
					Email:       "existinguser@example.com",
				}
				db.Create(existingUser)
			}
			app := api.New(&api.EngineOpts{
				Engine: gin.New(),
				Cfg:    cfg,
				DB:     db,
			})
			rec := tc.setupHTTP(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var body map[string]any
			if len(rec.Body.Bytes()) > 0 {
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				if err != nil {
					t.Logf("Failed to unmarshal body: %v, body: %s", err, rec.Body.String())
					body = make(map[string]any)
				}
			} else {
				body = make(map[string]any)
			}

			if tc.verifyBody != nil {
				tc.verifyBody(t, body)
			}

			if tc.verifyUser != nil {
				tc.verifyUser(t, db)
			}
		})
	}
}

func TestUserEndpoint_Login(t *testing.T) {
	t.Parallel()

	cfg := &api.Config{
		AppPort:     "8080",
		ServiceName: "bookmark-service",
		InstanceId:  "instance-1",
	}

	testCases := []struct {
		name           string
		setupHTTP      func(api.Engine) *httptest.ResponseRecorder
		setupJWT       func(t *testing.T) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator)
		expectedStatus int
		verifyBody     func(t *testing.T, body map[string]any)
	}{
		{
			name: "success - valid credentials",
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"username": "johndoe",
					"password": "P@ssw0rd11",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator) {
				generator := jwtMocks.NewJWTGenerator(t)
				generator.On("GenerateToken", mock.Anything).Return("mock-token", nil).Once()
				validator := jwtMocks.NewJWTValidator(t)
				return generator, validator
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body map[string]any) {
				token, ok := body["token"].(string)
				assert.True(t, ok, "token should be a string")
				assert.NotEmpty(t, token, "token should not be empty")
			},
		},
		{
			name: "invalid request body - missing username",
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"password": "password123",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				return generator, validator
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, body["message"], "Input error")
			},
		},
		{
			name: "invalid request body - missing password",
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"username": "johndoe",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				return generator, validator
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, body["message"], "Input error")
			},
		},
		{
			name: "invalid JSON",
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewReader([]byte("invalid json")))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				return generator, validator
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, body["message"], "Input Error")
			},
		},
		{
			name: "error - invalid password",
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"username": "johndoe",
					"password": "wrongpassword",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				return generator, validator
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody: func(t *testing.T, body map[string]any) {
				errorMsg, ok := body["error"].(string)
				if !ok {
					message, ok := body["message"].(string)
					assert.True(t, ok)
					assert.Contains(t, message, "invalid username or password")
					return
				}
				assert.Contains(t, errorMsg, "invalid username or password")
			},
		},
		{
			name: "error - user not found",
			setupHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := map[string]any{
					"username": "nonexistent",
					"password": "password123",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				return generator, validator
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody: func(t *testing.T, body map[string]any) {
				message, ok := body["message"].(string)
				assert.True(t, ok)
				assert.Equal(t, "invalid username or password", message)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			jwtGen, jwtVal := tc.setupJWT(t)
			app := api.New(&api.EngineOpts{
				Engine:       gin.New(),
				DB:           db,
				JWTGenerator: jwtGen,
				JWTValidator: jwtVal,
				Cfg:          cfg,
			})
			// app := api.New(cfg, nil, db, jwtGen, jwtVal)
			rec := tc.setupHTTP(app)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var body map[string]any
			if len(rec.Body.Bytes()) > 0 {
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				if err != nil {
					t.Logf("Failed to unmarshal body: %v, body: %s", err, rec.Body.String())
					body = make(map[string]any)
				}
			} else {
				body = make(map[string]any)
			}

			if tc.verifyBody != nil {
				tc.verifyBody(t, body)
			}
		})
	}
}

func TestUserEndpoint_GetProfile(t *testing.T) {
	t.Parallel()

	cfg := &api.Config{
		AppPort:     "8080",
		ServiceName: "bookmark-service",
		InstanceId:  "instance-1",
	}

	const mockUserID = "550e8400-e29b-41d4-a716-446655440000"

	testCases := []struct {
		name           string
		setupHTTP      func(api.Engine, string) *httptest.ResponseRecorder
		setupJWT       func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string)
		expectedStatus int
		verifyBody     func(t *testing.T, body map[string]any)
	}{
		{
			name: "success - valid token",
			setupHTTP: func(api api.Engine, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "valid-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": userID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, body["id"], mockUserID)
				assert.Equal(t, body["username"], "johndoe")
				assert.Equal(t, body["display_name"], "John Doe")
				assert.Equal(t, body["email"], "john.doe@example.com")
				_, exists := body["password"]
				assert.False(t, exists, "password should not be in response")
			},
		},
		{
			name: "error - missing Authorization header",
			setupHTTP: func(api api.Engine, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				return generator, validator, ""
			},
			expectedStatus: http.StatusUnauthorized,
			verifyBody: func(t *testing.T, body map[string]any) {
				errorMsg, ok := body["error"].(string)
				assert.True(t, ok)
				assert.Equal(t, "Authorization header is required", errorMsg)
			},
		},
		{
			name: "error - invalid token format",
			setupHTTP: func(api api.Engine, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				req.Header.Set("Authorization", "InvalidFormat token")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				return generator, validator, ""
			},
			expectedStatus: http.StatusUnauthorized,
			verifyBody: func(t *testing.T, body map[string]any) {
				errorMsg, ok := body["error"].(string)
				assert.True(t, ok)
				assert.Equal(t, "Authorization header format is wrong", errorMsg)
			},
		},
		{
			name: "error - invalid token",
			setupHTTP: func(api api.Engine, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				req.Header.Set("Authorization", "Bearer invalid.token.here")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				validator.On("ValidateToken", "invalid.token.here").
					Return(nil, errors.New("invalid token")).Once()
				return generator, validator, ""
			},
			expectedStatus: http.StatusUnauthorized,
			verifyBody: func(t *testing.T, body map[string]any) {
				errorMsg, ok := body["error"].(string)
				assert.True(t, ok)
				assert.Equal(t, "Invalid token", errorMsg)
			},
		},
		{
			name: "error - user not found",
			setupHTTP: func(api api.Engine, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "valid-nonexistent-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": "non-existent-user-id",
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusInternalServerError,
			verifyBody: func(t *testing.T, body map[string]any) {
				message, ok := body["message"].(string)
				assert.True(t, ok)
				assert.Equal(t, "Processing Error", message)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			jwtGen, jwtVal, token := tc.setupJWT(t, mockUserID)
			// app := api.New(cfg, nil, db, jwtGen, jwtVal)
			app := api.New(&api.EngineOpts{
				Engine:       gin.New(),
				DB:           db,
				JWTGenerator: jwtGen,
				JWTValidator: jwtVal,
				Cfg:          cfg,
			})
			rec := tc.setupHTTP(app, token)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var body map[string]any
			if len(rec.Body.Bytes()) > 0 {
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				if err != nil {
					t.Logf("Failed to unmarshal body: %v, body: %s", err, rec.Body.String())
					body = make(map[string]any)
				}
			} else {
				body = make(map[string]any)
			}

			if tc.verifyBody != nil {
				tc.verifyBody(t, body)
			}
		})
	}
}

func TestUserEndpoint_UpdateProfile(t *testing.T) {
	t.Parallel()

	cfg := &api.Config{
		AppPort:     "8080",
		ServiceName: "bookmark-service",
		InstanceId:  "instance-1",
	}

	const mockUserID = "550e8400-e29b-41d4-a716-446655440000"

	testCases := []struct {
		name           string
		setupHTTP      func(api.Engine, string) *httptest.ResponseRecorder
		setupJWT       func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string)
		expectedStatus int
		verifyBody     func(t *testing.T, body map[string]any)
	}{
		{
			name: "success - valid token and body",
			setupHTTP: func(api api.Engine, token string) *httptest.ResponseRecorder {
				body := map[string]any{
					"display_name": "John Updated",
					"email":        "john.updated@example.com",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPut, "/v1/self/info", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "valid-update-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": userID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, body["id"], mockUserID)
				assert.Equal(t, body["username"], "johndoe")
				assert.Equal(t, body["display_name"], "John Updated")
				assert.Equal(t, body["email"], "john.updated@example.com")
			},
		},
		{
			name: "error - missing Authorization header",
			setupHTTP: func(api api.Engine, token string) *httptest.ResponseRecorder {
				body := map[string]any{
					"display_name": "John Updated",
					"email":        "john.updated@example.com",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPut, "/v1/self/info", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				return generator, validator, ""
			},
			expectedStatus: http.StatusUnauthorized,
			verifyBody: func(t *testing.T, body map[string]any) {
				errorMsg, ok := body["error"].(string)
				assert.True(t, ok)
				assert.Equal(t, "Authorization header is required", errorMsg)
			},
		},
		{
			name: "error - invalid request body",
			setupHTTP: func(api api.Engine, token string) *httptest.ResponseRecorder {
				body := map[string]any{
					"display_name": "",
					"email":        "invalid-email",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPut, "/v1/self/info", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				api.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "invalid-body-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": userID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, body["message"], "Input error")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := fixture.NewFixture(t, &fixture.UserCommonTestDB{})
			jwtGen, jwtVal, token := tc.setupJWT(t, mockUserID)
			// app := api.New(cfg, nil, db, jwtGen, jwtVal)
			app := api.New(&api.EngineOpts{
				Engine:       gin.New(),
				DB:           db,
				JWTGenerator: jwtGen,
				JWTValidator: jwtVal,
				Cfg:          cfg,
			})
			rec := tc.setupHTTP(app, token)

			assert.Equal(t, tc.expectedStatus, rec.Code)

			var body map[string]any
			if len(rec.Body.Bytes()) > 0 {
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				if err != nil {
					t.Logf("Failed to unmarshal body: %v, body: %s", err, rec.Body.String())
					body = make(map[string]any)
				}
			} else {
				body = make(map[string]any)
			}

			if tc.verifyBody != nil {
				tc.verifyBody(t, body)
			}
		})
	}
}
