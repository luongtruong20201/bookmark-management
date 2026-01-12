package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/luongtruong20201/bookmark-management/internal/services"
	"github.com/luongtruong20201/bookmark-management/pkg/response"
)

type User interface {
	RegisterUser(c *gin.Context)
}

type user struct {
	svc service.User
}

func NewUser(svc service.User) User {
	return &user{
		svc: svc,
	}
}

type createUserInputBody struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func (u *user) RegisterUser(c *gin.Context) {
	body := &createUserInputBody{}
	if err := c.ShouldBindJSON(body); err != nil {
		c.JSON(http.StatusBadRequest, response.InputFieldError(err))
		return
	}

	res, err := u.svc.CreateUser(c, body.Username, body.Password, body.DisplayName, body.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.InternalErrResponse)
		return
	}

	c.JSON(http.StatusOK, res)
}
