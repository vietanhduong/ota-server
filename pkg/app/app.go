package app

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vietanhduong/ota-server/pkg/apis/v1/profile"
	"github.com/vietanhduong/ota-server/pkg/apis/v1/storage_object"
	"github.com/vietanhduong/ota-server/pkg/apis/v1/user"
	"github.com/vietanhduong/ota-server/pkg/cerrors"
	"github.com/vietanhduong/ota-server/pkg/logger"
	"github.com/vietanhduong/ota-server/pkg/middlewares"
	"github.com/vietanhduong/ota-server/pkg/mysql"
	"github.com/vietanhduong/ota-server/pkg/redis"
	"github.com/vietanhduong/ota-server/pkg/templates"
	"github.com/vietanhduong/ota-server/pkg/utils/env"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"
)

type App struct {
	Echo  *echo.Echo
	MySQL *mysql.DB
	Redis *redis.Client
}

func (a *App) Initialize() {
	a.Echo = echo.New()

	// configure server
	a.Echo.Pre(middleware.RemoveTrailingSlash())
	a.Echo.Use(middleware.Recover())
	a.Echo.Use(middleware.RequestID())
	a.Echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAcceptEncoding},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions, http.MethodHead, http.MethodPatch},
	}))
	// set default timeout
	a.Echo.Use(middlewares.Timeout(10 * time.Minute))

	// serve SPA
	a.Echo.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:   env.GetEnvAsStringOrFallback("STATIC_PATH", "./web"),
		Index:  "index.html",
		Browse: false,
		HTML5:  true,
	}))

	// register error handler
	a.Echo.HTTPErrorHandler = cerrors.HTTPErrorHandler

	// register template
	a.Echo.Renderer = &templates.Template{
		Templates: template.Must(template.ParseGlob("public/templates/*")),
	}

	// customize request log
	format := "\x1b[32mINFO\x1b[0m  | ${time_rfc3339} | ${status} | ${method} ${uri} \n"
	a.Echo.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: format,
		Output: a.Echo.Logger.Output(),
	}))

	// custom logger
	a.Echo.Logger.SetHeader("${level} | ${time_rfc3339} | ${short_file}:${line} | ${message}")
	logger.InitializeLogger()

	// initialize database connection
	// make sure you have injected the database configuration
	// into the environment
	db, err := mysql.InitializeDatabase(mysql.Config{
		Username: env.GetEnvAsStringOrFallback("DB_USERNAME", ""),
		Password: env.GetEnvAsStringOrFallback("DB_PASSWORD", ""),
		Host:     env.GetEnvAsStringOrFallback("DB_HOST", ""),
		Port:     env.GetEnvAsStringOrFallback("DB_PORT", ""),
		Instance: env.GetEnvAsStringOrFallback("DB_INSTANCE", ""),
	})
	if err != nil {
		a.Echo.Logger.Fatalf("initialize database connection failed with error: %+v", err)
	}

	a.MySQL = db

	// initialize redis connection
	redisClient, err := redis.InitializeConnection(redis.Config{
		Host: env.GetEnvAsStringOrFallback("REDIS_HOST", ""),
		Port: env.GetEnvAsStringOrFallback("REDIS_PORT", ""),
		DB:   env.GetEnvAsIntOrFallback("REDIS_DB", 0),
	})
	if err != nil {
		a.Echo.Logger.Fatalf("initialize redis connection failed with error: %+v", err)
	}

	a.Redis = redisClient

	// auto migrate database on startup
	if autoMigrate := env.GetEnvAsIntOrFallback("AUTO_MIGRATE", 0); autoMigrate == 1 {
		if err := a.MySQL.Migration(); err != nil {
			a.Echo.Logger.Fatalf("migrate database was error %+v", err)
		}
	}

	// register routers
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	// make sure you call  `Initialize` before run
	server := &http.Server{
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      a.Echo,
	}
	// run the server in a goroutine so that it doesn't block
	go func() {
		log.Printf("server is starting at addr: %s", addr)
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// wait for interrupt signal to gracefully shutdown the server
	// with a timeout of 10 seconds
	quit := make(chan os.Signal, 1)
	// kill (no param) default send SIGTERM
	// kill -2 is SIGINT (ctrl+c)
	// kill -9 is SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// block until receive signal
	<-quit
	// create a deadline wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// doesn't block if no connections, but will otherwise wait
	// until the timeout deadline
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("error when shutdown server with error: %+v", err)
	}
	log.Println("shutting down")
}

func (a *App) initializeRoutes() {
	g := a.Echo.Group("/api/v1")

	profile.Register(g, a.MySQL, a.Redis)
	storage_object.Register(g, a.MySQL, a.Redis)
	user.Register(g, a.MySQL, a.Redis)
}
