package endpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luongtruong20201/bookmark-management/internal/api"
	model "github.com/luongtruong20201/bookmark-management/internal/models"
	sqldb "github.com/luongtruong20201/bookmark-management/pkg/sql"
	"github.com/stretchr/testify/assert"
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
		setupDB        func(t *testing.T) *gorm.DB
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
			setupDB: func(t *testing.T) *gorm.DB {
				db := sqldb.InitMockDB(t)
				db.AutoMigrate(&model.User{})
				return db
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
			setupDB: func(t *testing.T) *gorm.DB {
				db := sqldb.InitMockDB(t)
				db.AutoMigrate(&model.User{})
				return db
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
			setupDB: func(t *testing.T) *gorm.DB {
				db := sqldb.InitMockDB(t)
				db.AutoMigrate(&model.User{})
				return db
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
			setupDB: func(t *testing.T) *gorm.DB {
				db := sqldb.InitMockDB(t)
				db.AutoMigrate(&model.User{})
				return db
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
			setupDB: func(t *testing.T) *gorm.DB {
				db := sqldb.InitMockDB(t)
				db.AutoMigrate(&model.User{})
				return db
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
			setupDB: func(t *testing.T) *gorm.DB {
				db := sqldb.InitMockDB(t)
				db.AutoMigrate(&model.User{})
				existingUser := &model.User{
					ID:          "550e8400-e29b-41d4-a716-446655440000",
					Username:    "existinguser",
					Password:    "hashedpassword",
					DisplayName: "Existing User",
					Email:       "existinguser@example.com",
				}
				db.Create(existingUser)
				return db
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

			db := tc.setupDB(t)
			app := api.New(cfg, nil, db)
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
