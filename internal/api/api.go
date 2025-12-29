package api

import (
	"fmt"
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
	cfg *Config
}

// New creates a new API engine instance with the provided configuration.
// It initializes the Gin router and registers all endpoints.
func New(cfg *Config) Engine {
	a := &api{
		app: gin.New(),
		cfg: cfg,
	}

	a.registerEndPoint()
	return a
}

// Start starts the HTTP server on the port specified in the configuration.
func (a *api) Start() error {
	return a.app.Run(fmt.Sprintf(":%s", a.cfg.AppPort))
}

// ServeHTTP implements the http.Handler interface, allowing the API to be used
// as a standard HTTP handler.
func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.app.ServeHTTP(w, r)
}

// registerEndPoint registers all API endpoints with their respective handlers.
// It initializes the password and healthcheck services and handlers, then
// registers the routes for password generation and healthcheck.
func (a *api) registerEndPoint() {
	passSvc := service.NewPassword()
	passHandler := handler.NewPassword(passSvc)

	healthcheckSvc := service.NewHealthcheck(a.cfg.ServiceName, a.cfg.InstanceId)
	healthcheckHandler := handler.NewHealthcheck(healthcheckSvc)

	a.app.GET("/gen-pass", passHandler.GenPass)
	a.app.GET("/health-check", healthcheckHandler.Check)
}
