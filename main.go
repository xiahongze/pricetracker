package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/xiahongze/pricetracker/server"
)

func main() {
	srv := server.Run()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// listen to os signals
	<-c
	log.Println("Shutting down server")
	if err := srv.Shutdown(nil); err != nil {
		log.Fatalln("Error shutting down server", err)
	}
	log.Println("Done")
}
