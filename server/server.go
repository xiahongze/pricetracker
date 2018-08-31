package server

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/xiahongze/pricetracker/handlers"
)

var port = "8080"

func init() {
	if val, ok := os.LookupEnv("PORT"); ok {
		port = val
	}
}

// Build returns a new echo server instance
func Build() *echo.Echo {
	e := echo.New()
	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} echo: ${method} ${uri} ${status} from ${remote_ip} latency=${latency_human} error=${error}\n",
	}))
	e.Use(middleware.Recover())

	// Routes
	e.POST("/create", handlers.CreateHandler)
	e.POST("/read", handlers.ReadHandler)

	return e
}
