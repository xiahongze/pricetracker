package server

import (
	"log"
	"net/http"
	"os"

	"github.com/xiahongze/pricetracker/handlers"
)

var port = "8080"

func init() {
	if val, ok := os.LookupEnv("PORT"); ok {
		port = val
	}
}

// Run serves the PriceTracker service
func Run() *http.Server {
	log.Println("INFO: Welcome to PriceTracker, the server that would save you money!")
	srv := &http.Server{Addr: ":" + port}
	http.Handle("/create", handlers.Validator(handlers.CreateHandler))
	http.Handle("/read", handlers.Validator(handlers.ReadHandler))
	go func() {
		log.Printf("INFO: Listening at http://localhost:%s\n", port)
		if err := srv.ListenAndServe(); err != nil {
			// cannot panic, because this probably is an intentional close
			log.Printf("ERROR: ListenAndServe error: %s", err)
		}
	}()
	return srv
}
