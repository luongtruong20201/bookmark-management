package main

import "github.com/luongtruong20201/bookmark-management/internal/infrastructure"

func main() {
	_ = infrastructure.CreateSqlDBAndMigrate()
}
