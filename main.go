package main

import (
	"log"
	"time"

	"os"
	"os/signal"
	"syscall"

	"github.com/xiahongze/pricetracker/gutils"
	"github.com/xiahongze/pricetracker/server"
)

func main() {
	srv := server.Run()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// scheduling tasks
	schd := gutils.Schedule(func() {
		log.Println("Running on schedule")
		log.Println(gutils.FetchData(10))
	}, time.Second*10)

	// listen to os signals
	<-c
	log.Println("Shutting down server")
	if err := srv.Shutdown(nil); err != nil {
		log.Fatalln("Error shutting down server", err)
	}
	schd <- true
	log.Println("Stopped scheduled tasks")
	log.Println("Done")

}
