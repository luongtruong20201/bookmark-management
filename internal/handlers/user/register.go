package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/luongtruong20201/bookmark-management/pkg/request"
	"github.com/luongtruong20201/bookmark-management/pkg/response"
	"github.com/rs/zerolog/log"
)

// createUserInputBody represents the request body for user registration.
// It contains all the required fields for creating a new user account.
//
// Fields:
//   - Username: Unique username for the user account (required, must be non-empty)
//   - Password: User password (required, minimum 6 characters)
//   - DisplayName: User's display name shown in the application (required, must be non-empty)
//   - Email: User's email address (required, must be a valid email format)
type createUserInputBody struct {
	Username    string `json:"username" binding:"required" example:"johndoe"`
	Password    string `json:"password" binding:"required,min=6" example:"password123"`
	DisplayName string `json:"display_name" binding:"required" example:"John Doe"`
	Email       string `json:"email" binding:"required,email" example:"john.doe@example.com"`
}

// RegisterUser handles the user registration endpoint request. It validates the input,
// creates a new user account with hashed password, and returns the created user information.
// @Summary Register new user
// @Description Create a new user account with username, password, display name, and email
// @Tags user
// @Accept json
// @Produce json
// @Param request body createUserInputBody true "User registration request"
// @Success 200 {object} model.User "Successfully created user"
// @Failure 400 {object} response.Message "Invalid request body or validation error"
// @Failure 500 {object} response.Message "Internal server error"
// @Router /v1/users/register [post]
func (u *user) RegisterUser(c *gin.Context) {
	body, err := request.BindInputFromRequest[createUserInputBody](c)
	if err != nil {
		return
	}

	res, err := u.svc.CreateUser(c, body.Username, body.Password, body.DisplayName, body.Email)
	if err != nil {
		switch {
		case errors.Is(err, dbutils.ErrDuplicationType):
			c.JSON(http.StatusBadRequest, response.Message{
				Message: "username or email already taken",
			})
			return
		case errors.Is(err, nil):
		default:
			log.Error().Err(err).Msg("error when generating password")
			c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Register an user successfully!",
		"data":    res,
	})
}
