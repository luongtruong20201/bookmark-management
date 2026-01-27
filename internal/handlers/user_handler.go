package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/luongtruong20201/bookmark-management/internal/services"
	"github.com/luongtruong20201/bookmark-management/pkg/response"
	"github.com/rs/zerolog/log"
)

// User defines the interface for user handlers.
// It provides methods to handle user-related requests.
type User interface {
	RegisterUser(c *gin.Context)
}

type user struct {
	svc service.User
}

// NewUser creates a new user handler with the provided user service.
func NewUser(svc service.User) User {
	return &user{
		svc: svc,
	}
}

// createUserInputBody represents the request body for user registration.
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
	body := &createUserInputBody{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	res, err := u.svc.CreateUser(c, body.Username, body.Password, body.DisplayName, body.Email)
	if err != nil {
		log.Error().Err(err).Msg("error when generating password")
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Register an user successfully!",
		"data":    res,
	})
}
