package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luongtruong20201/bookmark-management/internal/utils"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/luongtruong20201/bookmark-management/pkg/request"
	"github.com/luongtruong20201/bookmark-management/pkg/response"
	"github.com/rs/zerolog/log"
)

// updateProfileRequestBody represents the request body for updating user profile.
type updateProfileRequestBody struct {
	DisplayName string `json:"display_name" binding:"required" example:"John Doe"`
	Email       string `json:"email" binding:"required,email" example:"john.doe@example.com"`
}

// UpdateProfile updates the profile information (display name and email) of the currently authenticated user.
// The user ID is extracted from the JWT token by auth middleware and stored in the Gin context.
// @Summary Update user profile
// @Description Update the currently authenticated user's display name and email using the Bearer token
// @Tags user
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body updateProfileRequestBody true "User profile update request"
// @Success 200 {object} model.User "Updated user profile"
// @Failure 400 {object} response.Message "Invalid request body or validation error"
// @Failure 401 {object} response.Message "Unauthorized (missing/invalid token)"
// @Failure 500 {object} response.Message "Internal server error"
// @Router /v1/self/info [put]
func (u *user) UpdateProfile(c *gin.Context) {
	userId, err := utils.GetUserIDFromRequest(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Message{Message: "Invalid token"})
		return
	}

	body, err := request.BindInputFromRequest[updateProfileRequestBody](c)
	if err != nil {
		return
	}

	res, err := u.svc.UpdateUserProfile(c, userId, body.DisplayName, body.Email)
	if err != nil {
		switch {
		case errors.Is(err, dbutils.ErrDuplicationType):
			c.JSON(http.StatusBadRequest, response.Message{
				Message: "username or email already taken",
			})
			return
		default:
			log.Error().Err(err).Msg("error when updating user profile")
			c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
			return
		}
	}

	c.JSON(http.StatusOK, res)
}
