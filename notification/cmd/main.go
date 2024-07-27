package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Hello from Chat Service")
	})

	e.GET("/live", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Chat Service is alive")
	})

	e.GET("/ready", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Chat Service is alive")
	})

	e.GET("/start-up", func(c echo.Context) error {
		return c.HTML(http.StatusOK, "Chat Service is start up")
	})

	port, ok := os.LookupEnv("NOTIFICATION_PORT")
	if !ok {
		port = "8084"
	}

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}
