package request

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type testInput struct {
	Username string `json:"username" validate:"required"`
	UserID string `uri:"id" validate:"required"`
	Page int `form:"page" validate:"required"`
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
