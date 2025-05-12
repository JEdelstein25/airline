package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8001"
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	// Basic airline endpoints
	http.HandleFunc("/airlines", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id":1,"name":"Airline Sample","code":"ALS"}]`))
	})

	// Aircraft types endpoint
	http.HandleFunc("/aircraft-types", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id":1,"name":"Boeing 737","code":"B737"},{"id":2,"name":"Airbus A320","code":"A320"}]`))
	})

	log.Printf("Airline Service starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}