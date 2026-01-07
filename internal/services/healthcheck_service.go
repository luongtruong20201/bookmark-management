package service

//go:generate mockery --name Healthcheck --filename health_check_service.go
type Healthcheck interface {
	Check() (string, string, string)
}

type healthcheckService struct {
	serviceName string
	instanceId  string
}

// NewHealthcheck creates a new healthcheck service instance with the provided
// service name and instance ID.
func NewHealthcheck(serviceName string, instanceId string) Healthcheck {
	return &healthcheckService{
		serviceName: serviceName,
		instanceId:  instanceId,
	}
}

// Check performs a health check and returns the status message, service name, and instance ID.
func (s healthcheckService) Check() (string, string, string) {
	return "OK", s.serviceName, s.instanceId
}
