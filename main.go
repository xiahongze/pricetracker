package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/xiahongze/pricetracker/gutils"
	"github.com/xiahongze/pricetracker/handlers"
	"github.com/xiahongze/pricetracker/pushover"
)

func build(client *pushover.Client) *echo.Echo {
	e := echo.New()
	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} [echo] ${short_file}:${line}: ${method} ${uri} ${status} from ${remote_ip} latency=${latency_human} error=(${error})\n",
	}))
	e.Use(middleware.Recover())

	// Routes
	e.POST("/create", handlers.MakeCreate(client))
	e.POST("/read", handlers.Read)
	e.POST("/update", handlers.MakeUpdate(client))
	e.POST("/delete", handlers.MakeDelete(client))

	return e
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	var (
		appToken   string
		userToken  string
		port       string
		schdlFreq  int
		fetchLimit int
	)
	flag.StringVar(&appToken, "appToken", "", "pushover app token")
	flag.StringVar(&userToken, "userToken", "", "pushover user token")
	flag.StringVar(&port, "port", "8080", "server port")
	flag.IntVar(&schdlFreq, "schdlFreq", 2, "schedule frequency in minutes")
	flag.IntVar(&fetchLimit, "fetchLimit", 10, "fetch limit from google datastore")
	flag.Parse()
	if appToken == "" || userToken == "" {
		log.Fatalln("appToken and userToken must be given")
	}

	client := pushover.Client{AppToken: appToken, User: userToken}
	schdlFreqMin := time.Minute * time.Duration(schdlFreq)

	e := build(&client)

	// scheduling tasks
	schd := gutils.Schedule(func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("INFO: Main: Recovered in Schedule: ", r)
			}
		}()
		log.Println("INFO: Main: Running on schedule")
		gutils.Refresh(&client, fetchLimit)
	}, schdlFreqMin)

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
