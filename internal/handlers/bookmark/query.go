package bookmark

import (
	"net/http"

	"github.com/gin-gonic/gin"
	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/pkg/request"
	"github.com/luongtruong20201/bookmark-management/pkg/response"
	"github.com/rs/zerolog/log"
)

// getBookmarksResponse represents the response structure for GetBookmarks endpoint.
type getBookmarksResponse struct {
	Data       []*model.Bookmark          `json:"data"`
	Pagination response.PaginationMetadata `json:"pagination"`
}

// GetBookmarks handles the HTTP request to retrieve bookmarks for the authenticated user.
// It extracts pagination parameters (page, pageSize) from query parameters,
// gets the user ID from the JWT token, and delegates the retrieval to the bookmark service.
//
// @Summary List bookmarks
// @Description Get a paginated list of bookmarks for the authenticated user
// @Tags bookmark
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Items per page"
// @Success 200 {object} getBookmarksResponse "List of bookmarks with pagination"
// @Failure 400 {object} response.Message "Invalid pagination parameters"
// @Failure 401 {object} response.Message "Unauthorized (missing/invalid token)"
// @Failure 500 {object} response.Message "Internal server error"
// @Router /v1/bookmarks [get]
// @Security BearerAuth
func (h *bookmarkHandler) GetBookmarks(c *gin.Context) {
	pagination, userId, err := request.BindInputFromQueryWithAuth[request.PaginationQuery](c)
	if err != nil {
		return
	}

	page, pageSize := pagination.ValidateAndNormalize()
	offset, limit := pagination.ToOffsetLimit()

	result, err := h.svc.GetBookmarks(c, userId, offset, limit)
	if err != nil {
		log.Error().Err(err).Str("uid", userId).Msg("failed to get bookmarks")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, getBookmarksResponse{
		Data:       result.Data,
		Pagination: response.NewPaginationMetadata(page, pageSize, result.Total),
	})
}
