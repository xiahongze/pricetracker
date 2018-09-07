package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/xiahongze/pricetracker/gutils"
	"github.com/xiahongze/pricetracker/server"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	// build server
	e := server.Build()

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
		if err := e.Start(":" + server.Port); err != nil {
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
		e.Logger.Printf("shutting down with error")
		e.Logger.Fatal(err)
	}
}
