package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/luongtruong20201/bookmark-management/internal/services"
)

type passwordHandler struct {
	svc service.Password
}

// Password defines the interface for password handlers.
// It provides methods to generate random passwords.
type Password interface {
	GenPass(*gin.Context)
}

// NewPassword creates a new password handler with the provided password service.
func NewPassword(svc service.Password) Password {
	return &passwordHandler{
		svc: svc,
	}
}

// GenPass handles the password generation endpoint request. It generates a new password
// using the password service and returns it as a plain text response.
func (h *passwordHandler) GenPass(c *gin.Context) {
	pass, err := h.svc.GeneratePassword()
	if err != nil {
		c.String(http.StatusInternalServerError, "err")
	}

	c.String(http.StatusOK, pass)
}
