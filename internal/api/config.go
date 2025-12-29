package api

import (
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	AppPort     string `default:"8080" envconfig:"APP_PORT"`
	ServiceName string `default:"bookmark-api" envconfig:"SERVICE_NAME"`
	InstanceId  string `default:"" envconfig:"APP_INSTANCE_ID"`
}

// NewConfig creates a new configuration instance by reading environment variables.
// If APP_INSTANCE_ID is not set, it generates a new UUID for the instance ID.
func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	if cfg.InstanceId == "" {
		cfg.InstanceId = uuid.New().String()
	}

	return cfg, nil
}
