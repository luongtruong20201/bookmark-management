package infrastructure

import (
	"github.com/luongtruong20201/bookmark-management/pkg/common"
	sqldb "github.com/luongtruong20201/bookmark-management/pkg/sql"
	"gorm.io/gorm"
)

func CreateSqlDBAndMigrate() *gorm.DB {
	db, err := sqldb.NewClient("")
	common.HandleError(err)

	err = MigrateSQLDB(db, "file://./migrations", "up", 0)
	common.HandleError(err)

	return db
}
