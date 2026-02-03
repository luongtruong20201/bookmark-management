package bookmark

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	model "github.com/luongtruong20201/bookmark-management/internal/models"
	service "github.com/luongtruong20201/bookmark-management/internal/services/bookmark"
	serviceMocks "github.com/luongtruong20201/bookmark-management/internal/services/bookmark/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBookmarkHandler_Create(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	var (
		testErrService = errors.New("service error")
	)

	type requestBody struct {
		Description string `json:"description,omitempty"`
		URL         string `json:"url,omitempty"`
	}

	testCases := []struct {
		name           string
		requestBody    interface{}
		setupContext   func(c *gin.Context)
		setupService   func(t *testing.T, c *gin.Context) service.Service
		expectedStatus int
		verifyResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "success - valid request and token",
			requestBody: requestBody{
				Description: "My blog",
				URL:         "https://truonglq.com",
			},
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": "550e8400-e29b-41d4-a716-446655440000"})
			},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("Create", c, "My blog", "https://truonglq.com", "550e8400-e29b-41d4-a716-446655440000").
					Return(&model.Bookmark{
						Base: model.Base{
							ID: "11111111-2222-3333-4444-555555555555",
						},
						Description: "My blog",
						URL:         "https://truonglq.com",
						Code:        "abcd1234",
						UserID:      "550e8400-e29b-41d4-a716-446655440000",
					}, nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Data    model.Bookmark `json:"data"`
					Message string         `json:"message"`
				}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "Create a bookmark successfully!", resp.Message)
				assert.Equal(t, "11111111-2222-3333-4444-555555555555", resp.Data.ID)
				assert.Equal(t, "My blog", resp.Data.Description)
				assert.Equal(t, "https://truonglq.com", resp.Data.URL)
				assert.Equal(t, "abcd1234", resp.Data.Code)
				assert.Empty(t, resp.Data.User)
			},
		},
		{
			name: "error - invalid request body",
			requestBody: map[string]any{
				"url":  "not-a-valid-url",
				"desc": "invalid field",
			},
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": "550e8400-e29b-41d4-a716-446655440000"})
			},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				return serviceMocks.NewService(t)
			},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: nil,
		},
		{
			name: "error - missing user ID in token",
			requestBody: requestBody{
				Description: "My blog",
				URL:         "https://truonglq.com",
			},
			setupContext: func(c *gin.Context) {},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil, testErrService).Maybe()
				return svcMock
			},
			expectedStatus: http.StatusUnauthorized,
			verifyResponse: nil,
		},
		{
			name: "error - service returns error",
			requestBody: requestBody{
				Description: "My blog",
				URL:         "https://truonglq.com",
			},
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": "550e8400-e29b-41d4-a716-446655440000"})
			},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("Create", c, "My blog", "https://truonglq.com", "550e8400-e29b-41d4-a716-446655440000").
					Return(nil, testErrService).Once()
				return svcMock
			},
			expectedStatus: http.StatusInternalServerError,
			verifyResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var body map[string]any
				err := json.Unmarshal(rec.Body.Bytes(), &body)
				assert.NoError(t, err)
				message, ok := body["message"].(string)
				assert.True(t, ok)
				assert.Equal(t, "Processing Error", message)
			},
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

			ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks", bytes.NewBuffer(reqBody))
			ctx.Request.Header.Set("Content-Type", "application/json")

			tc.setupContext(ctx)
			svc := tc.setupService(t, ctx)
			h := NewBookmarkHandler(svc)

			h.Create(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.verifyResponse != nil {
				tc.verifyResponse(t, rec)
			}
		})
	}
}
