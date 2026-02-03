package request

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type testInput struct {
	Username  string `json:"username" validate:"required"`
	UserID    string `uri:"id" validate:"required"`
	Page      int    `form:"page" validate:"required"`
	RequestID string `header:"X-Request-Id" validate:"required"`
}

func TestBindInputFromRequest_Success(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	rec := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rec)

	r.POST("/users/:id", func(c *gin.Context) {
		input, err := BindInputFromRequest[testInput](c)
		assert.NoError(t, err)
		assert.NotNil(t, input)

		c.JSON(http.StatusOK, gin.H{
			"username":   input.Username,
			"user_id":    input.UserID,
			"page":       input.Page,
			"request_id": input.RequestID,
		})
	})

	bodyBytes, _ := json.Marshal(map[string]any{
		"username": "johndoe",
	})
	req := httptest.NewRequest(http.MethodPost, "/users/123?page=1", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", "req-1")

	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]any
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "johndoe", resp["username"])
	assert.Equal(t, "123", resp["user_id"])
	assert.Equal(t, float64(1), resp["page"])
	assert.Equal(t, "req-1", resp["request_id"])
}

func TestBindInputFromRequest_AbortOnError(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		req            func() *http.Request
		expectedStatus int
	}{
		{
			name: "error - invalid JSON body",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/users/123?page=1", bytes.NewBufferString("invalid json"))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Request-Id", "req-1")
				return req
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "error - missing URI param (route mismatch -> 404 before handler)",
			req: func() *http.Request {
				bodyBytes, _ := json.Marshal(map[string]any{"username": "johndoe"})
				req := httptest.NewRequest(http.MethodPost, "/users?page=1", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Request-Id", "req-1")
				return req
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "error - missing query param (validation fails)",
			req: func() *http.Request {
				bodyBytes, _ := json.Marshal(map[string]any{"username": "johndoe"})
				req := httptest.NewRequest(http.MethodPost, "/users/123", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Request-Id", "req-1")
				return req
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "error - missing header (validation fails)",
			req: func() *http.Request {
				bodyBytes, _ := json.Marshal(map[string]any{"username": "johndoe"})
				req := httptest.NewRequest(http.MethodPost, "/users/123?page=1", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				return req
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "error - missing username (validation fails)",
			req: func() *http.Request {
				bodyBytes, _ := json.Marshal(map[string]any{})
				req := httptest.NewRequest(http.MethodPost, "/users/123?page=1", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Request-Id", "req-1")
				return req
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			_, r := gin.CreateTestContext(rec)

			nextCalled := false
			r.POST("/users/:id",
				func(c *gin.Context) {
					_, _ = BindInputFromRequest[testInput](c)
					c.Next()
				},
				func(c *gin.Context) {
					nextCalled = true
					c.JSON(http.StatusOK, gin.H{"ok": true})
				},
			)

			r.ServeHTTP(rec, tc.req())

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.expectedStatus != http.StatusOK {
				assert.False(t, nextCalled)
			}
		})
	}
}

type queryInput struct {
	Page     int `form:"page" binding:"omitempty,gte=1"`
	PageSize int `form:"pageSize" binding:"omitempty,gte=1,lte=100"`
}

func TestBindInputFromRequestWithAuth_Success(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	const mockUserID = "550e8400-e29b-41d4-a716-446655440000"

	rec := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rec)

	r.POST("/users/:id", func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": mockUserID})
		input, userId, err := BindInputFromRequestWithAuth[testInput](c)
		assert.NoError(t, err)
		assert.NotNil(t, input)
		assert.Equal(t, mockUserID, userId)

		c.JSON(http.StatusOK, gin.H{
			"username": input.Username,
			"user_id":  input.UserID,
			"page":     input.Page,
			"uid":      userId,
		})
	})

	bodyBytes, _ := json.Marshal(map[string]any{
		"username": "johndoe",
	})
	req := httptest.NewRequest(http.MethodPost, "/users/123?page=1", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", "req-1")

	r.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]any
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, "johndoe", resp["username"])
	assert.Equal(t, "123", resp["user_id"])
	assert.Equal(t, float64(1), resp["page"])
	assert.Equal(t, mockUserID, resp["uid"])
}

func TestBindInputFromRequestWithAuth_Error(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	const mockUserID = "550e8400-e29b-41d4-a716-446655440000"

	testCases := []struct {
		name           string
		req            func() *http.Request
		setupContext   func(c *gin.Context)
		expectedStatus int
		shouldAbort    bool
	}{
		{
			name: "error - invalid JSON body",
			req: func() *http.Request {
				req := httptest.NewRequest(http.MethodPost, "/users/123?page=1", bytes.NewBufferString("invalid json"))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Request-Id", "req-1")
				return req
			},
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
		{
			name: "error - missing user ID in claims",
			req: func() *http.Request {
				bodyBytes, _ := json.Marshal(map[string]any{"username": "johndoe"})
				req := httptest.NewRequest(http.MethodPost, "/users/123?page=1", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Request-Id", "req-1")
				return req
			},
			setupContext: func(c *gin.Context) {
			},
			expectedStatus: http.StatusUnauthorized,
			shouldAbort:    true,
		},
		{
			name: "error - empty user ID in claims",
			req: func() *http.Request {
				bodyBytes, _ := json.Marshal(map[string]any{"username": "johndoe"})
				req := httptest.NewRequest(http.MethodPost, "/users/123?page=1", bytes.NewReader(bodyBytes))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Request-Id", "req-1")
				return req
			},
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": ""})
			},
			expectedStatus: http.StatusUnauthorized,
			shouldAbort:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			_, r := gin.CreateTestContext(rec)

			nextCalled := false
			r.POST("/users/:id",
				func(c *gin.Context) {
					tc.setupContext(c)
					_, _, err := BindInputFromRequestWithAuth[testInput](c)
					if err != nil {
						return
					}
					c.Next()
				},
				func(c *gin.Context) {
					nextCalled = true
					c.JSON(http.StatusOK, gin.H{"ok": true})
				},
			)

			r.ServeHTTP(rec, tc.req())

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.shouldAbort {
				assert.False(t, nextCalled)
			}
		})
	}
}

func TestBindInputFromQueryWithAuth_Success(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	const mockUserID = "550e8400-e29b-41d4-a716-446655440000"

	rec := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rec)

	r.GET("/bookmarks", func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": mockUserID})
		input, userId, err := BindInputFromQueryWithAuth[queryInput](c)
		assert.NoError(t, err)
		assert.NotNil(t, input)
		assert.Equal(t, mockUserID, userId)

		c.JSON(http.StatusOK, gin.H{
			"page":     input.Page,
			"pageSize": input.PageSize,
			"uid":      userId,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/bookmarks?page=2&pageSize=20", nil)
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]any
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, float64(2), resp["page"])
	assert.Equal(t, float64(20), resp["pageSize"])
	assert.Equal(t, mockUserID, resp["uid"])
}

func TestBindInputFromQueryWithAuth_WithDefaults(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	const mockUserID = "550e8400-e29b-41d4-a716-446655440000"

	rec := httptest.NewRecorder()
	_, r := gin.CreateTestContext(rec)

	r.GET("/bookmarks", func(c *gin.Context) {
		c.Set("claims", jwt.MapClaims{"sub": mockUserID})
		input, userId, err := BindInputFromQueryWithAuth[queryInput](c)
		assert.NoError(t, err)
		assert.NotNil(t, input)
		assert.Equal(t, mockUserID, userId)

		c.JSON(http.StatusOK, gin.H{
			"page":     input.Page,
			"pageSize": input.PageSize,
			"uid":      userId,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/bookmarks", nil)
	r.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]any
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	assert.Equal(t, float64(0), resp["page"])
	assert.Equal(t, float64(0), resp["pageSize"])
	assert.Equal(t, mockUserID, resp["uid"])
}

func TestBindInputFromQueryWithAuth_Error(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	const mockUserID = "550e8400-e29b-41d4-a716-446655440000"

	testCases := []struct {
		name           string
		req            func() *http.Request
		setupContext   func(c *gin.Context)
		expectedStatus int
		shouldAbort    bool
	}{
		{
			name: "error - invalid page value (validation fails)",
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/bookmarks?page=-1&pageSize=10", nil)
			},
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
		{
			name: "error - invalid pageSize value (exceeds max)",
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/bookmarks?page=1&pageSize=200", nil)
			},
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			expectedStatus: http.StatusBadRequest,
			shouldAbort:    true,
		},
		{
			name: "error - missing user ID in claims",
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/bookmarks?page=1&pageSize=10", nil)
			},
			setupContext: func(c *gin.Context) {

			},
			expectedStatus: http.StatusUnauthorized,
			shouldAbort:    true,
		},
		{
			name: "error - empty user ID in claims",
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/bookmarks?page=1&pageSize=10", nil)
			},
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": ""})
			},
			expectedStatus: http.StatusUnauthorized,
			shouldAbort:    true,
		},
		{
			name: "error - user ID is not a string",
			req: func() *http.Request {
				return httptest.NewRequest(http.MethodGet, "/bookmarks?page=1&pageSize=10", nil)
			},
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": 12345})
			},
			expectedStatus: http.StatusUnauthorized,
			shouldAbort:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			_, r := gin.CreateTestContext(rec)

			nextCalled := false
			r.GET("/bookmarks",
				func(c *gin.Context) {
					tc.setupContext(c)
					_, _, err := BindInputFromQueryWithAuth[queryInput](c)
					if err != nil {

						return
					}
					c.Next()
				},
				func(c *gin.Context) {
					nextCalled = true
					c.JSON(http.StatusOK, gin.H{"ok": true})
				},
			)

			r.ServeHTTP(rec, tc.req())

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.shouldAbort {
				assert.False(t, nextCalled)
			}
		})
	}
}
