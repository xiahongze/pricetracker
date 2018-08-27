package main

import (
	"log"

	"os"
	"os/signal"
	"syscall"

	"github.com/xiahongze/pricetracker/gutils"
	"github.com/xiahongze/pricetracker/server"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	srv := server.Run()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

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

	// listen to os signals
	<-c
	log.Println("INFO: Main: Shutting down server")
	if err := srv.Shutdown(nil); err != nil {
		log.Fatalln("ERROR: Main: Error shutting down server", err)
	}
	schd <- true
	log.Println("INFO: Main: Stopped scheduled tasks")
	log.Println("INFO: Main: Done")

}
