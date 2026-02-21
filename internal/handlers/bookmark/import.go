package bookmark

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luongtruong20201/bookmark-management/internal/services/queue"
	"github.com/luongtruong20201/bookmark-management/internal/utils"
	"github.com/luongtruong20201/bookmark-management/pkg/csv"
	"github.com/luongtruong20201/bookmark-management/pkg/request"
	"github.com/luongtruong20201/bookmark-management/pkg/response"
)

// ImportBookmarks handles importing multiple bookmarks from an uploaded CSV file
// for the authenticated user. The CSV file must contain the columns "description"
// and "url", and is sent as a multipart/form-data field named "file". After the
// file is parsed and validated, the bookmarks are dispatched as background jobs
// to the queue service for asynchronous processing.
//
// @Summary Import bookmarks
// @Description Import multiple bookmarks from a CSV file for the authenticated user
// @Tags bookmark
// @Accept mpfd
// @Produce json
// @Param file formData file true "CSV file containing bookmarks to import"
// @Success 200 {object} response.Message "Processing import"
// @Failure 400 {object} response.Message "Invalid file"
// @Failure 401 {object} response.Message "Unauthorized (missing/invalid token)"
// @Failure 500 {object} response.Message "Internal server error"
// @Router /v1/bookmarks/import [post]
// @Security BearerAuth
func (h *bookmarkHandler) ImportBookmarks(c *gin.Context) {
	uid, err := utils.GetUserIDFromRequest(c)
	if err != nil {
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, &response.Message{
			Message: "Invalid file",
		})
		return
	}

	var importInput []*queue.ImportBookmarkInput
	if err = csv.ParseFromMultipartFile(file, &importInput); err != nil {
		c.JSON(http.StatusBadRequest, &response.Message{
			Message: "Invalid file",
		})
		return
	}

	if err = request.ValidateStruct(importInput); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	if err = h.queue.SendBookmarkJob(c.Request.Context(), uid, importInput); err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}
	
	c.JSON(http.StatusOK, &response.Message{
		Message: "Processing import",
	})
}
