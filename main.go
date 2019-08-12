package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/xiahongze/pricetracker/gutils"
	"github.com/xiahongze/pricetracker/handlers"
)

var port = "8080"

func build() *echo.Echo {
	e := echo.New()
	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} [echo] ${short_file}:${line}: ${method} ${uri} ${status} from ${remote_ip} latency=${latency_human} error=(${error})\n",
	}))
	e.Use(middleware.Recover())

	// Routes
	e.POST("/create", handlers.Create)
	e.POST("/read", handlers.Read)
	e.POST("/update", handlers.Update)
	e.POST("/delete", handlers.Delete)

	return e
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	if val, ok := os.LookupEnv("PORT"); ok {
		port = val
	}
}

func main() {
	e := build()

	// scheduling tasks
	schd := gutils.Schedule(func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("INFO: Main: Recovered in Schedule: ", r)
			}
		}()
		log.Println("INFO: Main: Running on schedule")
		gutils.Refresh()
	}, gutils.SchdFreq)

	// Start server
	go func() {
		if err := e.Start(":" + port); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// 10 seconds count down if serve needs to wait
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	schd <- true
	log.Println("INFO: Main: Stopped scheduled tasks")
	log.Println("INFO: Main: Done")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Printf("ERROR: Main: shutting down with error")
		e.Logger.Fatal(err)
	}
}
