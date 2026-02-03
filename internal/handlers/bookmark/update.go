package bookmark

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/luongtruong20201/bookmark-management/pkg/request"
	"github.com/luongtruong20201/bookmark-management/pkg/response"
	"github.com/rs/zerolog/log"
)

// updateBookmarkInput represents the request body for updating a bookmark.
// All fields are optional except ID, which must be provided and match the path param.
type updateBookmarkInput struct {
	ID          string `json:"id" uri:"id" binding:"required"`
	Description string `json:"description" binding:"omitempty,lte=255"`
	URL         string `json:"url" binding:"omitempty,url,lte=2048"`
}

// UpdateBookmark handles the HTTP request to update an existing bookmark for the
// authenticated user. It validates the payload, extracts the user ID from the JWT token,
// and delegates the update to the bookmark service.
//
// @Summary Update bookmark
// @Description Update an existing bookmark for the authenticated user
// @Tags bookmark
// @Accept json
// @Produce json
// @Param id path string true "Bookmark ID"
// @Param request body updateBookmarkInput true "Bookmark update request"
// @Success 200 {object} response.Message "Successfully updated bookmark"
// @Failure 400 {object} response.Message "Invalid request body or validation error"
// @Failure 401 {object} response.Message "Unauthorized (missing/invalid token)"
// @Failure 404 {object} response.Message "Bookmark not found"
// @Failure 500 {object} response.Message "Internal server error"
// @Router /v1/bookmarks/{id} [put]
// @Security BearerAuth
func (h *bookmarkHandler) UpdateBookmark(c *gin.Context) {
	input, userId, err := request.BindInputFromRequestWithAuth[updateBookmarkInput](c)
	if err != nil {
		return
	}

	if input.ID != c.Param("id") {
		c.JSON(http.StatusBadRequest, &response.Message{
			Message: "ID in path and body must match",
		})
		return
	}

	_, err = h.svc.Update(c, input.ID, userId, input.Description, input.URL)
	if err != nil {
		log.Error().Err(err).Str("uid", userId).Str("bookmark_id", input.ID).Msg("failed to update bookmark")

		if errors.Is(err, dbutils.ErrNotFoundType) {
			c.JSON(http.StatusNotFound, &response.Message{
				Message: "Bookmark not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, &response.Message{
		Message: "Success",
	})
}
