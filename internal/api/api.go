package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/luongtruong20201/bookmark-management/docs"
	handler "github.com/luongtruong20201/bookmark-management/internal/handlers"
	repository "github.com/luongtruong20201/bookmark-management/internal/repositories"
	service "github.com/luongtruong20201/bookmark-management/internal/services"
	"github.com/luongtruong20201/bookmark-management/pkg/stringutils"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Engine defines the interface for the API engine.
// It provides methods to start the server and serve HTTP requests.
type Engine interface {
	Start() error
	http.Handler
}

type api struct {
	redis *redis.Client
	app   *gin.Engine
	cfg   *Config
}

// New creates a new API engine instance with the provided configuration.
// It initializes the Gin router and registers all endpoints.
func New(cfg *Config, redis *redis.Client) Engine {
	a := &api{
		redis: redis,
		app:   gin.New(),
		cfg:   cfg,
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

	healthCheckRepo := repository.NewHealthCheck(a.redis)
	healthcheckSvc := service.NewHealthcheck(a.cfg.ServiceName, a.cfg.InstanceId, healthCheckRepo)
	healthcheckHandler := handler.NewHealthcheck(healthcheckSvc)

	keyGen := stringutils.NewKeyGen()
	shortenRepo := repository.NewURLStorage(a.redis)
	shortenSvc := service.NewShortenURL(keyGen, shortenRepo)
	shortenHandler := handler.NewShortenURL(shortenSvc)

	a.app.GET("/gen-pass", passHandler.GenPass)
	a.app.GET("/health-check", healthcheckHandler.Check)

	v1 := a.app.Group("/v1")
	{
		v1.POST("/links/shorten", shortenHandler.ShortenURL)
		v1.GET("/links/redirect/:code", shortenHandler.GetURL)
	}

	a.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
