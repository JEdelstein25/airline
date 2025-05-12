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
		port = "8005"
	}

	r := mux.NewRouter()

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}).Methods("GET")

	// Basic booking endpoints
	r.HandleFunc("/api/bookings", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id":1,"passenger_name":"John Doe","flight_id":123}]`))
	}).Methods("GET")

	log.Printf("Booking Service starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}