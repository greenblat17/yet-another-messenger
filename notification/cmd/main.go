package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Chat Service API
// @version 1.0
// @description This is a sample chat service server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8084
// @BasePath /
func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// @Summary Root endpoint
	// @Description Returns a greeting message
	// @Tags root
	// @Accept  json
	// @Produce  html
	// @Success 200 {string} string "Hello from Chat Service"
	// @Router / [get]
	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello from Chat Service")
	})

	// @Summary Live endpoint
	// @Description Returns service live status
	// @Tags status
	// @Accept  json
	// @Produce  html
	// @Success 200 {string} string "Chat Service is alive"
	// @Router /live [get]
	e.GET("/live", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Chat Service is alive")
	})

	// @Summary Ready endpoint
	// @Description Returns service ready status
	// @Tags status
	// @Accept  json
	// @Produce  html
	// @Success 200 {string} string "Chat Service is alive"
	// @Router /ready [get]
	e.GET("/ready", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Chat Service is alive")
	})

	// @Summary Start-up endpoint
	// @Description Returns service start-up status
	// @Tags status
	// @Accept  json
	// @Produce  html
	// @Success 200 {string} string "Chat Service is start up"
	// @Router /start-up [get]
	e.GET("/start-up", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Chat Service is start up")
	})

	port, ok := os.LookupEnv("NOTIFICATION_PORT")
	if !ok {
		port = "8084"
	}

	// Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	err := e.Start(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
}
