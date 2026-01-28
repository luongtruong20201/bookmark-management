package main

import (
	"github.com/luongtruong20201/bookmark-management/internal/api"
	model "github.com/luongtruong20201/bookmark-management/internal/models"
	"github.com/luongtruong20201/bookmark-management/pkg/common"
	"github.com/luongtruong20201/bookmark-management/pkg/jwt"
	"github.com/luongtruong20201/bookmark-management/pkg/logger"
	"github.com/luongtruong20201/bookmark-management/pkg/redis"
	sqldb "github.com/luongtruong20201/bookmark-management/pkg/sql"
)

//	@title			Bookmark API
//	@version		1.0.0
//	@description	API documentation for bookmark service
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@host			localhost:8080
//	@BasePath		/

// main is the entry point of the application. It initializes the configuration,
// creates a new API instance, and starts the server.
func main() {
	cfg, err := api.NewConfig()
	common.HandleError(err)

	logger.SetLogLevel()

	redis, err := redis.NewClient("")
	common.HandleError(err)

	db, err := sqldb.NewClient("")
	common.HandleError(err)
	common.HandleError(db.AutoMigrate(&model.User{}))
	jwtGenerator, err := jwt.NewJWTGenerator("./private.pem")
	common.HandleError(err)

	jwtValidator, err := jwt.NewJWTValidator("./public.pem")
	common.HandleError(err)

	app := api.New(cfg, redis, db, jwtGenerator, jwtValidator)
	common.HandleError(app.Start())
}
