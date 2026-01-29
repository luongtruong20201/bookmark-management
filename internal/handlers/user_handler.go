package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/luongtruong20201/bookmark-management/internal/services"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	requtil "github.com/luongtruong20201/bookmark-management/pkg/request"
	"github.com/luongtruong20201/bookmark-management/pkg/response"
	"github.com/rs/zerolog/log"
)

// User defines the interface for user handlers.
// It provides methods to handle user-related requests.
type User interface {
	RegisterUser(c *gin.Context)
	Login(c *gin.Context)
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
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
	body, err := requtil.BindInputFromRequest[createUserInputBody](c)
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

// loginRequestBody represents the request body for user login.
type loginRequestBody struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// loginResponseBody represents the response body for successful login.
type loginResponseBody struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Login handles the user login endpoint request. It validates the credentials,
// authenticates the user, and returns a JWT token upon successful authentication.
// @Summary User login
// @Description Authenticate user with username and password, returns JWT token
// @Tags user
// @Accept json
// @Produce json
// @Param request body loginRequestBody true "User login credentials"
// @Success 200 {object} loginResponseBody "Successfully authenticated, returns JWT token"
// @Failure 400 {object} response.Message "Invalid credentials or validation error"
// @Failure 500 {object} response.Message "Internal server error"
// @Router /v1/users/login [post]
func (u *user) Login(c *gin.Context) {
	body, err := requtil.BindInputFromRequest[loginRequestBody](c)
	if err != nil {
		return
	}

	token, err := u.svc.Login(c, body.Username, body.Password)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrClientErr):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		case errors.Is(err, dbutils.ErrNotFoundType):
			c.JSON(http.StatusBadRequest, response.Message{
				Message: "invalid username or password",
			})
			return
		case errors.Is(err, nil):
		default:
			c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
			return
		}
	}

	c.JSON(http.StatusOK, &loginResponseBody{
		Token: token,
	})
}

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
	userIDValue, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Message{Message: "Invalid token"})
		return
	}

	userId, ok := userIDValue.(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Message{Message: "Invalid token"})
		return
	}

	body, err := requtil.BindInputFromRequest[updateProfileRequestBody](c)
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
	userIDValue, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, response.Message{Message: "Invalid token"})
		return
	}

	userId, ok := userIDValue.(string)
	if !ok {
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
