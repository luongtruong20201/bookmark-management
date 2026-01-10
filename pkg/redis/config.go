package redis

import (
	"github.com/kelseyhightower/envconfig"
)

// config holds the Redis connection configuration loaded from environment variables.
type config struct {
	Address  string `default:"localhost:6379" envconfig:"REDIS_ADDRESS"`
	Password string `default:"" envconfig:"REDIS_PASSWORD"`
	DB       int    `default:"0" envconfig:"REDIS_DB"`
}

// newConfig creates a new configuration instance by reading environment variables.
// The prefix parameter is used to prefix environment variable names (currently unused).
func newConfig(prefix string) (*config, error) {
	cfg := &config{}

	if err := envconfig.Process(prefix, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
