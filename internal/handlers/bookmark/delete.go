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

type deleteBookmarkInput struct {
	ID string `uri:"id" binding:"required"`
}

// DeleteBookmark handles the HTTP request to delete a bookmark for the
// authenticated user. It extracts the bookmark ID from the URI, gets the user ID
// from the JWT token, and delegates the deletion to the bookmark service.
//
// @Summary Delete bookmark
// @Description Delete a bookmark for the authenticated user
// @Tags bookmark
// @Accept json
// @Produce json
// @Param id path string true "Bookmark ID"
// @Success 200 {object} response.Message "Successfully deleted bookmark"
// @Failure 400 {object} response.Message "Invalid request or validation error"
// @Failure 401 {object} response.Message "Unauthorized (missing/invalid token)"
// @Failure 404 {object} response.Message "Bookmark not found"
// @Failure 500 {object} response.Message "Internal server error"
// @Router /v1/bookmarks/{id} [delete]
// @Security BearerAuth
func (h *bookmarkHandler) DeleteBookmark(c *gin.Context) {
	input, userId, err := request.BindInputFromUriWithAuth[deleteBookmarkInput](c)
	if err != nil {
		return
	}

	err = h.svc.Delete(c, input.ID, userId)
	if err != nil {
		log.Error().Err(err).Str("uid", userId).Str("bookmark_id", input.ID).Msg("failed to delete bookmark")

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

