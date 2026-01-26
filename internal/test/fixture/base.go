package fixture

import (
	"testing"

	sqldb "github.com/luongtruong20201/bookmark-management/pkg/sql"
	"gorm.io/gorm"
)

type Fixture interface {
	SetupDB(*gorm.DB)
	Migrate() error
	GenerateData() error
	DB() *gorm.DB
}

func NewFixture(t *testing.T, fix Fixture) *gorm.DB {
	fix.SetupDB(sqldb.InitMockDB(t))
	if err := fix.Migrate(); err != nil {
		t.Fatal("failed to migrate db for testing", err)
	}

	if err := fix.GenerateData(); err != nil {
		t.Fatal("failed to generate data for testing", err)
	}

	return fix.DB()
}
