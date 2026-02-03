package bookmark

import (
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
)

func TestBookmarkHandler_GetBookmarks(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	const mockUserID = "550e8400-e29b-41d4-a716-446655440000"

	var (
		testErrService = errors.New("service error")
	)

	testCases := []struct {
		name           string
		queryParams    string
		setupContext   func(c *gin.Context)
		setupService   func(t *testing.T, c *gin.Context) service.Service
		expectedStatus int
		verifyResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:        "success - get bookmarks with pagination",
			queryParams: "?page=1&pageSize=10",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				svcMock := serviceMocks.NewService(t)
				bookmarks := []*model.Bookmark{
					{
						Base: model.Base{
							ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
						},
						Description: "Facebook",
						URL:         "https://www.facebook.com",
						Code:        "abc1234",
						UserID:      mockUserID,
					},
					{
						Base: model.Base{
							ID: "b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e",
						},
						Description: "Google",
						URL:         "https://www.google.com",
						Code:        "def5678",
						UserID:      mockUserID,
					},
				}
				svcMock.On("GetBookmarks", c, mockUserID, 0, 10).Return(bookmarks, nil).Once()
				svcMock.On("CountBookmarks", c, mockUserID).Return(int64(25), nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Data       []*model.Bookmark `json:"data"`
					Pagination struct {
						Page     int   `json:"page"`
						Limit    int   `json:"limit"`
						Total    int64 `json:"total"`
					} `json:"pagination"`
				}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Len(t, resp.Data, 2)
				assert.Equal(t, 1, resp.Pagination.Page)
				assert.Equal(t, 10, resp.Pagination.Limit)
				assert.Equal(t, int64(25), resp.Pagination.Total)
			},
		},
		{
			name:        "success - empty result",
			queryParams: "?page=1&pageSize=10",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("GetBookmarks", c, mockUserID, 0, 10).Return([]*model.Bookmark{}, nil).Once()
				svcMock.On("CountBookmarks", c, mockUserID).Return(int64(0), nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Data       []*model.Bookmark `json:"data"`
					Pagination struct {
						Total int64 `json:"total"`
					} `json:"pagination"`
				}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Empty(t, resp.Data)
				assert.Equal(t, int64(0), resp.Pagination.Total)
			},
		},
		{
			name:        "success - with default pagination",
			queryParams: "",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				svcMock := serviceMocks.NewService(t)
				bookmarks := []*model.Bookmark{
					{
						Base: model.Base{
							ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
						},
						Description: "Facebook",
						URL:         "https://www.facebook.com",
						Code:        "abc1234",
						UserID:      mockUserID,
					},
				}
				svcMock.On("GetBookmarks", c, mockUserID, 0, 10).Return(bookmarks, nil).Once()
				svcMock.On("CountBookmarks", c, mockUserID).Return(int64(1), nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Data       []*model.Bookmark `json:"data"`
					Pagination struct {
						Page  int `json:"page"`
						Limit int `json:"limit"`
					} `json:"pagination"`
				}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Len(t, resp.Data, 1)
				assert.Equal(t, 1, resp.Pagination.Page)
				assert.Equal(t, 10, resp.Pagination.Limit)
			},
		},
		{
			name:        "success - page 2 with pageSize 5",
			queryParams: "?page=2&pageSize=5",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				svcMock := serviceMocks.NewService(t)
				bookmarks := []*model.Bookmark{
					{
						Base: model.Base{
							ID: "c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f",
						},
						Description: "GitHub",
						URL:         "https://www.github.com",
						Code:        "ghi9012",
						UserID:      mockUserID,
					},
				}
				svcMock.On("GetBookmarks", c, mockUserID, 5, 5).Return(bookmarks, nil).Once()
				svcMock.On("CountBookmarks", c, mockUserID).Return(int64(10), nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Data       []*model.Bookmark `json:"data"`
					Pagination struct {
						Page  int `json:"page"`
						Limit int `json:"limit"`
					} `json:"pagination"`
				}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Len(t, resp.Data, 1)
				assert.Equal(t, 2, resp.Pagination.Page)
				assert.Equal(t, 5, resp.Pagination.Limit)
			},
		},
		{
			name:        "error - invalid pagination parameters",
			queryParams: "?page=-1&pageSize=10",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				return serviceMocks.NewService(t)
			},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: nil,
		},
		{
			name:        "error - missing user ID in token",
			queryParams: "?page=1&pageSize=10",
			setupContext: func(c *gin.Context) {
			},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				return serviceMocks.NewService(t)
			},
			expectedStatus: http.StatusUnauthorized,
			verifyResponse: nil,
		},
		{
			name:        "error - GetBookmarks service returns error",
			queryParams: "?page=1&pageSize=10",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("GetBookmarks", c, mockUserID, 0, 10).Return(nil, testErrService).Once()
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
		{
			name:        "error - CountBookmarks service returns error",
			queryParams: "?page=1&pageSize=10",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				svcMock := serviceMocks.NewService(t)
				bookmarks := []*model.Bookmark{
					{
						Base: model.Base{
							ID: "a1b2c3d4-e5f6-4a7b-8c9d-0e1f2a3b4c5d",
						},
						Description: "Facebook",
						URL:         "https://www.facebook.com",
						Code:        "abc1234",
						UserID:      mockUserID,
					},
				}
				svcMock.On("GetBookmarks", c, mockUserID, 0, 10).Return(bookmarks, nil).Once()
				svcMock.On("CountBookmarks", c, mockUserID).Return(int64(0), testErrService).Once()
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
		{
			name:        "error - pageSize exceeds max (validation fails)",
			queryParams: "?page=1&pageSize=200",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": mockUserID})
			},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				return serviceMocks.NewService(t)
			},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/bookmarks"+tc.queryParams, nil)

			tc.setupContext(ctx)
			svc := tc.setupService(t, ctx)
			h := NewBookmarkHandler(svc)

			h.GetBookmarks(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.verifyResponse != nil {
				tc.verifyResponse(t, rec)
			}
		})
	}
}

