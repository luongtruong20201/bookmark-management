package sqldb

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// config holds the database connection configuration parameters.
// It can be populated from environment variables with an optional prefix.
type config struct {
	Host     string `default:"localhost" envconfig:"DB_HOST"`
	Port     string `default:"5433" envconfig:"DB_PORT"`
	User     string `default:"postgres" envconfig:"DB_USERNAME"`
	Password string `default:"postgres" envconfig:"DB_PASSWORD"`
	DBName   string `default:"bookmark_service" envconfig:"DB_NAME"`
}

// newConfig creates a new database configuration by reading environment variables
// with the specified prefix. Returns an error if environment variable processing fails.
func newConfig(prefix string) (*config, error) {
	cfg := &config{}
	if err := envconfig.Process(prefix, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// GetDSN returns a PostgreSQL Data Source Name (DSN) string constructed from
// the configuration values. This DSN can be used to establish a database connection.
func (cfg *config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port)
}
