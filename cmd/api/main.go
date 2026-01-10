package main

import (
	"github.com/luongtruong20201/bookmark-management/internal/api"
	"github.com/luongtruong20201/bookmark-management/pkg/logger"
	"github.com/luongtruong20201/bookmark-management/pkg/redis"
)

//	@title			Bookmark API
//	@version		1.0.0
//	@description	API documentation for bookmark service
//	@host			localhost:8080
//	@BasePath		/

// main is the entry point of the application. It initializes the configuration,
// creates a new API instance, and starts the server.
func main() {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	logger.SetLogLevel()

	redis, err := redis.NewClient("")
	if err != nil {
		panic(err)
	}

	app := api.New(cfg, redis)
	app.Start()
}
