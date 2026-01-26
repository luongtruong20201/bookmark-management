package sqldb

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// InitMockDB creates an in-memory SQLite database instance for testing purposes.
// It generates a unique database file name using UUID and configures GORM with
// silent logging. This is useful for unit tests that require database operations
// without setting up a real PostgreSQL database. The database is automatically
// cleaned up when the test completes.
func InitMockDB(t *testing.T) *gorm.DB {
	cdn := fmt.Sprintf("file:%s?mode=memory&cache=shared", uuid.New().String())
	db, err := gorm.Open(sqlite.Open(cdn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		t.Fatal("fail to create db: ", err)
	}

	return db
}
