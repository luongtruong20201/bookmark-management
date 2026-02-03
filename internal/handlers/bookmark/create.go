package bookmark

import (
	"net/http"

	"github.com/gin-gonic/gin"
	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/pkg/request"
	"github.com/luongtruong20201/bookmark-management/pkg/response"
	"github.com/rs/zerolog/log"
)

// createBookmarkInput represents the request body for creating a bookmark.
// It contains an optional description and a required valid URL.
type createBookmarkInput struct {
	Description string `json:"description" binding:"lte=255"`
	URL         string `json:"url" binding:"required,url,lte=2048"`
}

// createBookmarkResponse represents the response body for a successful bookmark creation.
// It wraps the created bookmark in a data field and includes a human-readable message.
type createBookmarkResponse struct {
	Data    *model.Bookmark `json:"data"`
	Message string          `json:"message"`
}

// Create handles the HTTP request to create a new bookmark for the
// authenticated user. It validates the payload, extracts the user ID
// from the JWT token and delegates the creation to the bookmark service.
//
// @Summary Create bookmark
// @Description Create a new bookmark for the authenticated user
// @Tags bookmark
// @Accept json
// @Produce json
// @Param request body createBookmarkInput true "Bookmark create request"
// @Success 200 {object} createBookmarkResponse "Create a bookmark successfully"
// @Failure 400 {object} response.Message "Invalid request body or validation error"
// @Failure 401 {object} response.Message "Unauthorized (missing/invalid token)"
// @Failure 500 {object} response.Message "Internal server error"
// @Router /v1/bookmarks [post]
// @Security BearerAuth
func (h *bookmarkHandler) Create(c *gin.Context) {
	body, userId, err := request.BindInputFromRequestWithAuth[createBookmarkInput](c)
	if err != nil {
		return
	}

	res, err := h.svc.Create(c, body.Description, body.URL, userId)
	if err != nil {
		log.Error().Err(err).Str("uid", userId).Msg("failed to create bookmark")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, createBookmarkResponse{
		Data:    res,
		Message: "Create a bookmark successfully!",
	})
}
