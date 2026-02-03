package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luongtruong20201/bookmark-management/internal/utils"
	"github.com/luongtruong20201/bookmark-management/pkg/response"
)

// GetProfile returns the profile information of the currently authenticated user.
// The user ID is extracted from the JWT token by auth middleware and stored in the Gin context.
// @Summary Get user profile
// @Description Get the currently authenticated user's profile using the Bearer token
// @Tags user
// @Security BearerAuth
// @Produce json
// @Success 200 {object} model.User "User profile"
// @Failure 401 {object} response.Message "Unauthorized (missing/invalid token)"
// @Failure 500 {object} response.Message "Internal server error"
// @Router /v1/self/info [get]
func (u *user) GetProfile(c *gin.Context) {
	userId, err := utils.GetUserIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Message{Message: "Invalid token"})
		return
	}

	res, err := u.svc.GetUserByID(c, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, res)
}
