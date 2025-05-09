package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.learning/api/auth"
	"go.learning/api/user"
	"go.learning/config"
	"go.learning/middlewares"
	"go.learning/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	"github.com/go-redis/redis/v8"
)

var conf config.Config

func init() {
	var err error
	conf, err = config.LoadConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	e := echo.New()
	// Set Cors origin and methods
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Set up middleware for logging
	e.Use(middleware.Logger())

	// Set up database connection
	dbPG := initDBPortgre(conf.Databasepostgres)

	// Set up Redis connection
	redisClient := initRedis(conf.Redis)

	// Automigrate the database
	migrate(dbPG)

	// Register routes
	go registerRoutes(e, dbPG, redisClient, conf)

	// Set up graceful shutdown
	waitForGracefulShutdown(e)
}

func registerRoutes(e *echo.Echo, dbPG *gorm.DB, redisClient *redis.Client, cfg config.Config) {

	userRepository := user.NewRepository(dbPG)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	authService := auth.NewService(userRepository, redisClient, cfg)
	authHandler := auth.NewHandler(authService)

	// Auth routes
	e.POST("/login", authHandler.Login)
	e.POST("/refresh-token", authHandler.RefreshToken)
	e.POST("/logout", authHandler.Logout)

	// User routes
	e.POST("/register", userHandler.Register)

	user_routes := e.Group("/user")

	user_routes.GET("/:id", middlewares.TokenAuthMiddleware(userHandler.Get, redisClient, cfg.JWT.SecretKey))
	user_routes.PUT("", middlewares.TokenAuthMiddleware(userHandler.Update, redisClient, cfg.JWT.SecretKey))
	user_routes.DELETE("/:id", middlewares.TokenAuthMiddleware(userHandler.Delete, redisClient, cfg.JWT.SecretKey))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Server.Port)))
}

func initDBPortgre(c config.Databasepostgres) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		c.Host,
		c.Port,
		c.DBName,
		c.Username,
		c.Password,
		c.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicf("error connecting to DBPortgre: %v", err)
	}
	log.Infof("connected to database Portgre %s:%d", c.Host, c.Port)
	return db
}

func initRedis(c config.Redis) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Host, c.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		log.Panicf("error connecting to Redis: %v", err)
	}
	log.Infof("connected to Redis %s:%d", c.Host, c.Port)
	return redisClient
}

func waitForGracefulShutdown(e *echo.Echo) {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}, &models.Logger{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	fmt.Println("Database migration completed!")
}
