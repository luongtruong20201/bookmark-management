package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/luongtruong20201/bookmark-management/internal/api"
	"github.com/luongtruong20201/bookmark-management/pkg/common"
)

func CreateAPIConfig() *api.Config {
	cfg, err := api.NewConfig()
	common.HandleError(err)
	return cfg
}

func CreateAPI() api.Engine {
	cfg := CreateAPIConfig()
	redis := CreateRedis()
	jwtGennerator, jwtValidator := CreateJWTProvider()
	db := CreateSqlDBAndMigrate()
	engine := gin.New()

	return api.New(&api.EngineOpts{
		Engine:       engine,
		Cfg:          cfg,
		Redis:        redis,
		DB:           db,
		JWTGenerator: jwtGennerator,
		JWTValidator: jwtValidator,
	})
}
