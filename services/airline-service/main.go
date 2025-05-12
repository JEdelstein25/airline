package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	r := chi.NewRouter()

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	// Basic airline endpoints
	r.Get("/airlines", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id":1,"name":"Airline Sample","code":"ALS"}]`))
	})

	// Aircraft types endpoint
	r.Get("/aircraft-types", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id":1,"name":"Boeing 737","code":"B737"},{"id":2,"name":"Airbus A320","code":"A320"}]`))
	})

	log.Printf("Airline Service starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}