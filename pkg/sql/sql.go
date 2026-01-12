package sqldb

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Host     string `default:"localhost" envconfig:"DB_HOST"`
	Port     string `default:"5433" envconfig:"DB_PORT"`
	User     string `default:"postgres" envconfig:"DB_USERNAME"`
	Password string `default:"postgres" envconfig:"DB_PASSWORD"`
	DBName   string `default:"bookmark_service" envconfig:"DB_NAME"`
}

func newConfig(prefix string) (*config, error) {
	cfg := &config{}
	if err := envconfig.Process(prefix, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port)
}
