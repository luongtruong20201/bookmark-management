package main

import (
	"github.com/luongtruong20201/bookmark-management/internal/infrastructure"
	"github.com/luongtruong20201/bookmark-management/pkg/common"
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
	api := infrastructure.CreateAPI()
	common.HandleError(api.Start())
}
