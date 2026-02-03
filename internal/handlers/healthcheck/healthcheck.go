package healthcheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/luongtruong20201/bookmark-management/internal/services/healthcheck"
)

// Healthcheck defines the interface for healthcheck handlers.
// It provides methods to check the health status of the service.
type Healthcheck interface {
	Check(*gin.Context)
}

// healthcheckHandler implements the Healthcheck interface and provides HTTP handlers
// for health check operations. It encapsulates the healthcheck service dependency
// for business logic execution.
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
// @Summary Health check
// @Description Check the health status of the service
// @Tags health
// @Success 200 {object} map[string]string "Health status response"
// @Router /health-check [get]
func (h *healthcheckHandler) Check(c *gin.Context) {
	message, serviceName, instanceId := h.healthcheckSvc.Check(c)
	c.JSON(http.StatusOK, gin.H{
		"message":      message,
		"service_name": serviceName,
		"instance_id":  instanceId,
	})
}
