package sqldb

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewClient creates a new GORM database client connection to a PostgreSQL database.
// It reads database configuration from environment variables using the specified prefix,
// constructs a DSN, and establishes a connection. Returns the GORM DB instance or an error
// if configuration loading or connection establishment fails.
func NewClient(prefix string) (*gorm.DB, error) {
	cfg, err := newConfig(prefix)
	if err != nil {
		return nil, err
	}

	dsn := cfg.GetDSN()
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	return db, nil
}
