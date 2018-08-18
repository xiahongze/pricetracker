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
func Run() {
	log.Println("Welcome to PriceTracker, the server that would save you money!")
	http.Handle("/create", handlers.Validator(handlers.CreateHandler))
	http.Handle("/read", handlers.Validator(handlers.ReadHandler))
	log.Printf("Listening at http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
