package main

import (
	"fmt"
	"net/http"

	"go.learning/config"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	e.Use(middleware.Logger())
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", conf.Server.Port)))
}
