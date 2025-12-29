package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	handler "github.com/luongtruong20201/bookmark-management/handlers"
	service "github.com/luongtruong20201/bookmark-management/services"
)

type Engine interface {
	Start() error
	http.Handler
}

type api struct {
	app *gin.Engine
}

func New() Engine {
	a := &api{
		app: gin.New(),
	}

	a.registerEndPoint()
	return a
}

func (a *api) Start() error {
	return a.app.Run(":8080")
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.app.ServeHTTP(w, r)
}

func (a *api) registerEndPoint() {
	passSvc := service.NewPassword()
	passHandler := handler.NewPassword(passSvc)

	a.app.GET("/gen-pass", passHandler.GenPass)
}
