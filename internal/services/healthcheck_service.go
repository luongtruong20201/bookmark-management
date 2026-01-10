package service

import (
	"context"

	repository "github.com/luongtruong20201/bookmark-management/internal/repositories"
)

// Healthcheck defines the interface for healthcheck services.
// It provides methods to check the health status of the service.
//
//go:generate mockery --name Healthcheck --filename health_check_service.go
type Healthcheck interface {
	Check(context.Context) (string, string, string)
}

type healthcheckService struct {
	serviceName string
	instanceId  string
	healthCheck repository.HealthCheck
}

// NewHealthcheck creates a new healthcheck service instance with the provided
// service name and instance ID.
func NewHealthcheck(
	serviceName string,
	instanceId string,
	healthCheck repository.HealthCheck,
) Healthcheck {
	return &healthcheckService{
		serviceName: serviceName,
		instanceId:  instanceId,
		healthCheck: healthCheck,
	}
}

// Check performs a health check and returns the status message, service name, and instance ID.
func (s healthcheckService) Check(ctx context.Context) (string, string, string) {
	if err := s.healthCheck.Ping(ctx); err != nil {
		return "NOT_OK", s.serviceName, s.instanceId
	}

	return "OK", s.serviceName, s.instanceId
}
