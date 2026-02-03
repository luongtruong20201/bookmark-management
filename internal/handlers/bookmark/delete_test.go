package bookmark

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	service "github.com/luongtruong20201/bookmark-management/internal/services/bookmark"
	serviceMocks "github.com/luongtruong20201/bookmark-management/internal/services/bookmark/mocks"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkHandler_DeleteBookmark(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	var (
		testErrService = errors.New("service error")
	)

	testCases := []struct {
		name           string
		bookmarkID     string
		setupContext   func(c *gin.Context)
		setupService   func(t *testing.T, c *gin.Context) service.Service
		expectedStatus int
		verifyResponse func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name:       "success - delete bookmark",
			bookmarkID: FixtureBookmarkIDFacebook,
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": FixtureUserIDJohnDoe})
			},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("Delete", c, FixtureBookmarkIDFacebook, FixtureUserIDJohnDoe).
					Return(nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Message string `json:"message"`
				}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "Success", resp.Message)
			},
		},
		{
			name:       "error - missing bookmark ID in URI",
			bookmarkID: "",
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": FixtureUserIDJohnDoe})
			},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				return serviceMocks.NewService(t)
			},
			expectedStatus: http.StatusBadRequest,
			verifyResponse: nil,
		},
		{
			name:       "error - missing user ID in token",
			bookmarkID: FixtureBookmarkIDFacebook,
			setupContext: func(c *gin.Context) {},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				return serviceMocks.NewService(t)
			},
			expectedStatus: http.StatusUnauthorized,
			verifyResponse: nil,
		},
		{
			name:       "error - bookmark not found",
			bookmarkID: FixtureBookmarkIDFacebook,
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": FixtureUserIDJohnDoe})
			},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("Delete", c, FixtureBookmarkIDFacebook, FixtureUserIDJohnDoe).
					Return(dbutils.ErrNotFoundType).Once()
				return svcMock
			},
			expectedStatus: http.StatusNotFound,
			verifyResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Message string `json:"message"`
				}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "Bookmark not found", resp.Message)
			},
		},
		{
			name:       "error - service returns internal error",
			bookmarkID: FixtureBookmarkIDFacebook,
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": FixtureUserIDJohnDoe})
			},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("Delete", c, FixtureBookmarkIDFacebook, FixtureUserIDJohnDoe).
					Return(testErrService).Once()
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
			name:       "success - delete bookmark with different ID",
			bookmarkID: FixtureBookmarkIDGoogle,
			setupContext: func(c *gin.Context) {
				c.Set("claims", jwt.MapClaims{"sub": FixtureUserIDJohnDoe})
			},
			setupService: func(t *testing.T, c *gin.Context) service.Service {
				svcMock := serviceMocks.NewService(t)
				svcMock.On("Delete", c, FixtureBookmarkIDGoogle, FixtureUserIDJohnDoe).
					Return(nil).Once()
				return svcMock
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp struct {
					Message string `json:"message"`
				}
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "Success", resp.Message)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			url := "/v1/bookmarks/" + tc.bookmarkID
			ctx.Request = httptest.NewRequest(http.MethodDelete, url, nil)
			ctx.Params = gin.Params{
				{Key: "id", Value: tc.bookmarkID},
			}

			tc.setupContext(ctx)
			svc := tc.setupService(t, ctx)
			h := NewBookmarkHandler(svc)

			h.DeleteBookmark(ctx)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			if tc.verifyResponse != nil {
				tc.verifyResponse(t, rec)
			}
		})
	}
}

