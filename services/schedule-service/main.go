package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8004"
	}

	r := mux.NewRouter()

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}).Methods("GET")

	// Basic schedule endpoints
	r.HandleFunc("/api/schedules", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id":1,"flight_number":"AB123","departure":"2025-05-15T10:00:00Z"}]`))
	}).Methods("GET")

	log.Printf("Schedule Service starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}