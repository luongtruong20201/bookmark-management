package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/luongtruong20201/bookmark-management/services"
)

type Healthcheck interface {
	Check(*gin.Context)
}

type healthcheckHandler struct {
	healthcheckSvc service.Healthcheck
}

// NewHealthcheck creates a new healthcheck handler with the provided healthcheck service.
func NewHealthcheck(svc service.Healthcheck) Healthcheck {
	return &healthcheckHandler{
		healthcheckSvc: svc,
	}
}

// Check handles the healthcheck endpoint request. It calls the healthcheck service
// and returns a JSON response with the status message, service name, and instance ID.
func (h *healthcheckHandler) Check(c *gin.Context) {
	message, serviceName, instanceId := h.healthcheckSvc.Check()
	c.JSON(http.StatusOK, gin.H{
		"message":      message,
		"service_name": serviceName,
		"instance_id":  instanceId,
	})
}
