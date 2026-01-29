package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/luongtruong20201/bookmark-management/docs"
	_ "github.com/luongtruong20201/bookmark-management/docs"
	"github.com/luongtruong20201/bookmark-management/internal/api/middlewares"
	handler "github.com/luongtruong20201/bookmark-management/internal/handlers"
	repository "github.com/luongtruong20201/bookmark-management/internal/repositories"
	service "github.com/luongtruong20201/bookmark-management/internal/services"
	jwtPkg "github.com/luongtruong20201/bookmark-management/pkg/jwt"
	"github.com/luongtruong20201/bookmark-management/pkg/stringutils"
	"github.com/luongtruong20201/bookmark-management/pkg/utils"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// Engine defines the interface for the API engine.
// It provides methods to start the server and serve HTTP requests.
type Engine interface {
	Start() error
	http.Handler
}

// handlers holds all HTTP handlers for the API endpoints.
// It groups together handlers for password generation, health checks,
// URL shortening, and user management.
type handlers struct {
	password    handler.Password
	healthCheck handler.Healthcheck
	shorten     handler.ShortenURL
	user        handler.User
}

// api represents the API server instance.
// It contains the Redis client for caching, database connection,
// Gin router engine, and configuration settings.
type api struct {
	redis        *redis.Client
	db           *gorm.DB
	app          *gin.Engine
	cfg          *Config
	jwtGenerator jwtPkg.JWTGenerator
	jwtValidator jwtPkg.JWTValidator
}

// New creates a new API engine instance with the provided configuration.
// It initializes the Gin router and registers all endpoints.
func New(
	cfg *Config,
	redis *redis.Client,
	db *gorm.DB,
	jwtGenerator jwtPkg.JWTGenerator,
	jwtValidator jwtPkg.JWTValidator,
) Engine {
	a := &api{
		redis:        redis,
		db:           db,
		app:          gin.New(),
		cfg:          cfg,
		jwtGenerator: jwtGenerator,
		jwtValidator: jwtValidator,
	}

	a.initRoutes()
	return a
}

// initHandlers initializes and returns all HTTP handlers for the API.
// It creates the necessary services and repositories, then instantiates
// handlers for password generation, health checks, URL shortening, and user management.
func (a *api) initHandlers() *handlers {
	passSvc := service.NewPassword()
	passHandler := handler.NewPassword(passSvc)

	healthCheckRepo := repository.NewHealthCheck(a.redis)
	healthcheckSvc := service.NewHealthcheck(a.cfg.ServiceName, a.cfg.InstanceId, healthCheckRepo)
	healthcheckHandler := handler.NewHealthcheck(healthcheckSvc)

	keyGen := stringutils.NewKeyGen()
	shortenRepo := repository.NewURLStorage(a.redis)
	shortenSvc := service.NewShortenURL(keyGen, shortenRepo)
	shortenHandler := handler.NewShortenURL(shortenSvc)

	hasher := utils.NewHasher()
	userRepo := repository.NewUser(a.db)
	userSvc := service.NewUser(userRepo, hasher, a.jwtGenerator)
	userHandler := handler.NewUser(userSvc)

	return &handlers{
		password:    passHandler,
		healthCheck: healthcheckHandler,
		shorten:     shortenHandler,
		user:        userHandler,
	}
}

// initRoutes registers all API routes with their corresponding handlers.
// It sets up endpoints for password generation, health checks, URL shortening,
// user registration, and Swagger documentation.
func (a *api) initRoutes() {
	handlers := a.initHandlers()

	a.app.GET("/gen-pass", handlers.password.GenPass)
	a.app.GET("/health-check", handlers.healthCheck.Check)

	v1Public := a.app.Group("/v1")
	{
		v1Public.POST("/links/shorten", handlers.shorten.ShortenURL)
		v1Public.GET("/links/redirect/:code", handlers.shorten.GetURL)

		v1Public.POST("/users/register", handlers.user.RegisterUser)
		v1Public.POST("/users/login", handlers.user.Login)
	}

	jwtMiddleware := middlewares.NewJWTAuth(a.jwtValidator)
	v1Private := a.app.Group("/v1")
	v1Private.Use(jwtMiddleware.JWTAuth())
	{
		v1Private.GET("/self/info", handlers.user.GetProfile)
		v1Private.PUT("/self/info", handlers.user.UpdateProfile)
	}

	a.app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	docs.SwaggerInfo.Host = a.cfg.AppHostname
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
