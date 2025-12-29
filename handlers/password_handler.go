package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/luongtruong20201/bookmark-management/services"
)

type passwordHandler struct {
	svc service.Password
}

type Password interface {
	GenPass(*gin.Context)
}

func NewPassword(svc service.Password) Password {
	return &passwordHandler{
		svc: svc,
	}
}

func (h *passwordHandler) GenPass(c *gin.Context) {
	pass, err := h.svc.GeneratePassword()
	if err != nil {
		c.String(http.StatusInternalServerError, "err")
	}

	c.String(http.StatusOK, pass)
}
