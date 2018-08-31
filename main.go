package main

import (
	"log"

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
	e.Logger.Fatal(e.Start(":8080"))
	schd <- true
	log.Println("INFO: Main: Stopped scheduled tasks")
	log.Println("INFO: Main: Done")
}
