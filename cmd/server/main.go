package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.learning/config"
	"go.learning/models"
	"go.learning/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
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

	// Set up middleware for logging
	e.Use(middleware.Logger())

	// Set up database connection
	dbPG := initDBPortgre(conf.Databasepostgres)

	// Automigrate the database
	migrate(dbPG)

	// Register routes
	go registerRoutes(e, dbPG, conf)

	// Set up graceful shutdown
	waitForGracefulShutdown(e)
}

func registerRoutes(e *echo.Echo, dbPG *gorm.DB, cfg config.Config) {

	userRepository := user.NewRepository(dbPG)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)
	e.POST("/register", userHandler.Register)
	e.GET("/user/:id", userHandler.Get)
	e.PUT("/user", userHandler.Update)
	e.DELETE("/user/:id", userHandler.Delete)

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
