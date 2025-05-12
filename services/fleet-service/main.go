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
		port = "8003"
	}

	r := mux.NewRouter()

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}).Methods("GET")

	// Basic fleet endpoints
	r.HandleFunc("/api/fleet", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id":1,"aircraft_type":"Boeing 737","registration":"N12345"}]`))
	}).Methods("GET")

	log.Printf("Fleet Service starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}