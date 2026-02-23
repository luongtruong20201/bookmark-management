package infrastructure

import (
	"github.com/luongtruong20201/bookmark-management/pkg/common"
	sqldb "github.com/luongtruong20201/bookmark-management/pkg/sql"
	"gorm.io/gorm"
)

// CreateSqlDB creates a database connection without running migrations.
// Use this when the schema is already applied (e.g. worker after API has migrated).
func CreateSqlDB() *gorm.DB {
	db, err := sqldb.NewClient("")
	common.HandleError(err)
	return db
}

func CreateSqlDBAndMigrate() *gorm.DB {
	db := CreateSqlDB()

	err := MigrateSQLDB(db, "file://./migrations", "up", 0)
	common.HandleError(err)

	return db
}
