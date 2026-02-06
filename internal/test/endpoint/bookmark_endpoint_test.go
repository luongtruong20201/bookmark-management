package endpoint

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/luongtruong20201/bookmark-management/internal/api"
	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/internal/test/fixture"
	jwtPkg "github.com/luongtruong20201/bookmark-management/pkg/jwt"
	jwtMocks "github.com/luongtruong20201/bookmark-management/pkg/jwt/mocks"
	redisPkg "github.com/luongtruong20201/bookmark-management/pkg/redis"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookmarkEndpoint_CreateBookmark(t *testing.T) {
	gin.SetMode(gin.TestMode)
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
		verifyDB       func(t *testing.T, db *gorm.DB)
	}{
		{
			name: "success - create bookmark with valid token and body",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				body := map[string]any{
					"description": "My blog",
					"url":         "https://truonglq.com",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/bookmarks", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "valid-bookmark-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": userID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body map[string]any) {
				data, ok := body["data"].(map[string]any)
				assert.True(t, ok, "response should contain data field")

				assert.Equal(t, "My blog", data["description"])
				assert.Equal(t, "https://truonglq.com", data["url"])
				assert.NotEmpty(t, data["id"])
				assert.NotEmpty(t, data["code"])
				_, hasUser := data["user"]
				assert.False(t, hasUser, "user should not be included in response")

				msg, ok := body["message"].(string)
				assert.True(t, ok)
				assert.Equal(t, "Create a bookmark successfully!", msg)
			},
			verifyDB: func(t *testing.T, db *gorm.DB) {
				var bookmark model.Bookmark
				err := db.Where("description = ? AND url = ?", "My blog", "https://truonglq.com").First(&bookmark).Error
				assert.NoError(t, err)
				assert.Equal(t, mockUserID, bookmark.UserID)
				assert.NotEmpty(t, bookmark.Code)
			},
		},
		{
			name: "error - missing Authorization header",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				body := map[string]any{
					"description": "My blog",
					"url":         "https://truonglq.com",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/bookmarks", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

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
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				body := map[string]any{
					"description": "My blog",
					"url":         "invalid-url",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/v1/bookmarks", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "invalid-body-bookmark-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": userID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody:     nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			jwtGen, jwtVal, token := tc.setupJWT(t, mockUserID)
			redis := redisPkg.InitMockRedis(t)

			app := api.New(&api.EngineOpts{
				Engine:       gin.New(),
				DB:           db,
				Redis:        redis,
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

			if tc.verifyDB != nil {
				tc.verifyDB(t, db)
			}
		})
	}
}

func TestBookmarkEndpoint_GetBookmarks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Parallel()

	cfg := &api.Config{
		AppPort:     "8080",
		ServiceName: "bookmark-service",
		InstanceId:  "instance-1",
	}

	const mockUserID = "550e8400-e29b-41d4-a716-446655440000"

	type testCase struct {
		name           string
		setupHTTP      func(api.Engine, string) *httptest.ResponseRecorder
		setupJWT       func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string)
		expectedStatus int
		verifyBody     func(t *testing.T, body map[string]any)
	}

	tests := []testCase{
		{
			name: "success - get bookmarks with pagination",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&pageSize=10", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "valid-get-bookmarks-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": userID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body map[string]any) {
				data, ok := body["data"].([]any)
				assert.True(t, ok, "response should have data field")
				assert.Len(t, data, 2, "should return 2 bookmarks for johndoe")

				pagination, ok := body["pagination"].(map[string]any)
				assert.True(t, ok, "response should have pagination field")
				assert.Equal(t, float64(1), pagination["page"])
				assert.Equal(t, float64(10), pagination["limit"])
				assert.Equal(t, float64(2), pagination["total"])

				if len(data) > 0 {
					bookmark, ok := data[0].(map[string]any)
					assert.True(t, ok)
					assert.NotEmpty(t, bookmark["id"])
					assert.NotEmpty(t, bookmark["description"])
					assert.NotEmpty(t, bookmark["url"])
					assert.NotEmpty(t, bookmark["code"])
				}
			},
		},
		{
			name: "success - get bookmarks with default pagination",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "valid-get-bookmarks-default-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": userID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body map[string]any) {
				pagination, ok := body["pagination"].(map[string]any)
				assert.True(t, ok)
				assert.Equal(t, float64(1), pagination["page"])
				assert.Equal(t, float64(10), pagination["limit"])
			},
		},
		{
			name: "success - get bookmarks page 2 with pageSize 1",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=2&pageSize=1", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "valid-get-bookmarks-page2-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": userID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body map[string]any) {
				data, ok := body["data"].([]any)
				assert.True(t, ok)
				assert.Len(t, data, 1, "should return 1 bookmark on page 2")

				pagination, ok := body["pagination"].(map[string]any)
				assert.True(t, ok)
				assert.Equal(t, float64(2), pagination["page"])
				assert.Equal(t, float64(1), pagination["limit"])
				assert.Equal(t, float64(2), pagination["total"])
			},
		},
		{
			name: "success - get bookmarks for user with no bookmarks",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&pageSize=10", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "valid-get-bookmarks-empty-token"
				emptyUserID := "00000000-0000-0000-0000-000000000000"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": emptyUserID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body map[string]any) {
				data, ok := body["data"].([]any)
				assert.True(t, ok)
				assert.Empty(t, data, "should return empty array for user with no bookmarks")

				pagination, ok := body["pagination"].(map[string]any)
				assert.True(t, ok)
				assert.Equal(t, float64(0), pagination["total"])
			},
		},
		{
			name: "success - get bookmarks for another user",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&pageSize=10", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "valid-get-bookmarks-another-user-token"
				anotherUserID := "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": anotherUserID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body map[string]any) {
				data, ok := body["data"].([]any)
				assert.True(t, ok)
				assert.Len(t, data, 2, "should return 2 bookmarks for an.nguyen")

				pagination, ok := body["pagination"].(map[string]any)
				assert.True(t, ok)
				assert.Equal(t, float64(2), pagination["total"])

				for _, item := range data {
					bookmark, ok := item.(map[string]any)
					assert.True(t, ok)
					assert.NotEmpty(t, bookmark["id"])
					assert.NotEmpty(t, bookmark["description"])
					assert.NotEmpty(t, bookmark["url"])
					assert.NotEmpty(t, bookmark["code"])
				}
			},
		},
		{
			name: "error - missing Authorization header",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&pageSize=10", nil)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

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
			name: "error - invalid pagination parameters (negative page)",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=-1&pageSize=10", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "invalid-pagination-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": userID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody:     nil,
		},
		{
			name: "error - invalid pagination parameters (pageSize exceeds max)",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&pageSize=200", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "invalid-pagesize-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": userID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody:     nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			jwtGen, jwtVal, token := tc.setupJWT(t, mockUserID)
			redis := redisPkg.InitMockRedis(t)

			app := api.New(&api.EngineOpts{
				Engine:       gin.New(),
				DB:           db,
				Redis:        redis,
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

func TestBookmarkEndpoint_UpdateBookmark(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Parallel()

	cfg := &api.Config{
		AppPort:     "8080",
		ServiceName: "bookmark-service",
		InstanceId:  "instance-1",
	}

	const (
		fixtureUserIDJohnDoe      = "550e8400-e29b-41d4-a716-446655440000"
		fixtureUserIDAnNguyen     = "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91"
		fixtureBookmarkIDFacebook = "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"
	)

	testCases := []struct {
		name           string
		setupHTTP      func(api.Engine, string) *httptest.ResponseRecorder
		setupJWT       func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string)
		expectedStatus int
		verifyBody     func(t *testing.T, body map[string]any)
		verifyDB       func(t *testing.T, db *gorm.DB)
	}{
		{
			name: "success - update bookmark with description and URL",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				body := map[string]any{
					"id":          fixtureBookmarkIDFacebook,
					"description": "Updated Facebook Description",
					"url":         "https://www.facebook.com/updated",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+fixtureBookmarkIDFacebook, bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "valid-update-bookmark-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": userID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, "Success", body["message"])
			},
			verifyDB: func(t *testing.T, db *gorm.DB) {
				var bookmark model.Bookmark
				err := db.Where("id = ?", fixtureBookmarkIDFacebook).First(&bookmark).Error
				assert.NoError(t, err)
				assert.Equal(t, "Updated Facebook Description", bookmark.Description)
				assert.Equal(t, "https://www.facebook.com/updated", bookmark.URL)
			},
		},
		{
			name: "success - update bookmark with only description",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				body := map[string]any{
					"id":          fixtureBookmarkIDFacebook,
					"description": "Updated Description Only",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+fixtureBookmarkIDFacebook, bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "valid-update-desc-only-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": userID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, "Success", body["message"])
			},
			verifyDB: func(t *testing.T, db *gorm.DB) {
				var bookmark model.Bookmark
				err := db.Where("id = ?", fixtureBookmarkIDFacebook).First(&bookmark).Error
				assert.NoError(t, err)
				assert.Equal(t, "Updated Description Only", bookmark.Description)
				assert.Equal(t, "https://www.facebook.com", bookmark.URL)
			},
		},
		{
			name: "error - bookmark not found",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				body := map[string]any{
					"id":          "00000000-0000-0000-0000-000000000000",
					"description": "Updated Description",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPut, "/v1/bookmarks/00000000-0000-0000-0000-000000000000", bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "valid-update-notfound-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": userID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusNotFound,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, "Bookmark not found", body["message"])
			},
		},
		{
			name: "error - bookmark belongs to different user",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				body := map[string]any{
					"id":          fixtureBookmarkIDFacebook,
					"description": "Updated Description",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+fixtureBookmarkIDFacebook, bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "valid-update-different-user-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": fixtureUserIDJohnDoe,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusNotFound,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, "Bookmark not found", body["message"])
			},
		},
		{
			name: "error - missing Authorization header",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				body := map[string]any{
					"id":          fixtureBookmarkIDFacebook,
					"description": "Updated Description",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+fixtureBookmarkIDFacebook, bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

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
			name: "error - invalid request body (invalid URL)",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				body := map[string]any{
					"id":  fixtureBookmarkIDFacebook,
					"url": "not-a-valid-url",
				}
				jsBody, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+fixtureBookmarkIDFacebook, bytes.NewReader(jsBody))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "valid-update-invalid-url-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": userID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusBadRequest,
			verifyBody:     nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			jwtGen, jwtVal, token := tc.setupJWT(t, fixtureUserIDAnNguyen)
			redis := redisPkg.InitMockRedis(t)

			app := api.New(&api.EngineOpts{
				Engine:       gin.New(),
				DB:           db,
				Redis:        redis,
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

			if tc.verifyDB != nil {
				tc.verifyDB(t, db)
			}
		})
	}
}

func TestBookmarkEndpoint_DeleteBookmark(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Parallel()

	cfg := &api.Config{
		AppPort:     "8080",
		ServiceName: "bookmark-service",
		InstanceId:  "instance-1",
	}

	const (
		fixtureUserIDJohnDoe      = "550e8400-e29b-41d4-a716-446655440000"
		fixtureUserIDAnNguyen     = "9b5c1e3e-7c3b-4f4e-8e7c-6e7a2f5d3a91"
		fixtureBookmarkIDFacebook = "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d"
	)

	tests := []struct {
		name           string
		setupHTTP      func(api.Engine, string) *httptest.ResponseRecorder
		setupJWT       func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string)
		expectedStatus int
		verifyBody     func(t *testing.T, body map[string]any)
		verifyDB       func(t *testing.T, db *gorm.DB)
	}{
		{
			name: "success - delete bookmark",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodDelete, "/v1/bookmarks/"+fixtureBookmarkIDFacebook, nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "valid-delete-bookmark-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": userID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusOK,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, "Success", body["message"])
			},
			verifyDB: func(t *testing.T, db *gorm.DB) {
				var bookmark model.Bookmark
				err := db.Where("id = ?", fixtureBookmarkIDFacebook).First(&bookmark).Error
				assert.Error(t, err, "bookmark should be deleted")
			},
		},
		{
			name: "error - bookmark not found",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodDelete, "/v1/bookmarks/00000000-0000-0000-0000-000000000000", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "valid-delete-notfound-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": userID,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusNotFound,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, "Bookmark not found", body["message"])
			},
		},
		{
			name: "error - bookmark belongs to different user",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodDelete, "/v1/bookmarks/"+fixtureBookmarkIDFacebook, nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

				return rec
			},
			setupJWT: func(t *testing.T, userID string) (jwtPkg.JWTGenerator, jwtPkg.JWTValidator, string) {
				generator := jwtMocks.NewJWTGenerator(t)
				validator := jwtMocks.NewJWTValidator(t)
				token := "valid-delete-different-user-token"
				validator.On("ValidateToken", token).Return(jwt.MapClaims{
					"sub": fixtureUserIDJohnDoe,
					"iat": 1600000000,
					"exp": 1600086400,
				}, nil).Once()

				return generator, validator, token
			},
			expectedStatus: http.StatusNotFound,
			verifyBody: func(t *testing.T, body map[string]any) {
				assert.Equal(t, "Bookmark not found", body["message"])
			},
		},
		{
			name: "error - missing Authorization header",
			setupHTTP: func(app api.Engine, token string) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodDelete, "/v1/bookmarks/"+fixtureBookmarkIDFacebook, nil)
				rec := httptest.NewRecorder()

				app.ServeHTTP(rec, req)

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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			db := fixture.NewFixture(t, &fixture.BookmarkCommonTestDB{})
			jwtGen, jwtVal, token := tc.setupJWT(t, fixtureUserIDAnNguyen)
			redis := redisPkg.InitMockRedis(t)

			app := api.New(&api.EngineOpts{
				Engine:       gin.New(),
				DB:           db,
				Redis:        redis,
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

			if tc.verifyDB != nil {
				tc.verifyDB(t, db)
			}
		})
	}
}
