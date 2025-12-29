package main

import (
	"github.com/luongtruong20201/bookmark-management/internal/api"
)

// main is the entry point of the application. It initializes the configuration,
// creates a new API instance, and starts the server.
func main() {
	cfg, err := api.NewConfig()
	if err != nil {
		panic(err)
	}

	app := api.New(cfg)
	app.Start()
}
