package user

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/luongtruong20201/bookmark-management/internal/services/user"
	"github.com/luongtruong20201/bookmark-management/pkg/dbutils"
	"github.com/luongtruong20201/bookmark-management/pkg/request"
	"github.com/luongtruong20201/bookmark-management/pkg/response"
)

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
	body, err := request.BindInputFromRequest[loginRequestBody](c)
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
