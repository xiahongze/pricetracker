package handlers

import (
	"log"
	"net/http"
	"strings"
)

// Validator verifies common request errors
func Validator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Hit " + r.URL.RequestURI())
		if ctype := r.Header.Get("content-type"); strings.ToLower(ctype) != "application/json" {
			http.Error(w, "only accept content-type: application/json", http.StatusBadRequest)
			return
		}
		if r.Body == nil {
			http.Error(w, "request body is empty!", http.StatusBadRequest)
			return
		}
		next.ServeHTTP(w, r)
	})
}
