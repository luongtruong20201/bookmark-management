package main

import "github.com/luongtruong20201/bookmark-management/internal/api"

func main() {
	app := api.New()
	app.Start()
}
