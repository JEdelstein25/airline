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
		port = "8002"
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	// Airports endpoint with sample data
	http.HandleFunc("/airports", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		airports := []map[string]interface{}{
			{
				"id": 1,
				"name": "John F. Kennedy International Airport",
				"iataCode": "JFK",
				"point": map[string]float64{
					"latitude": 40.6413,
					"longitude": -73.7781,
				},
			},
			{
				"id": 2,
				"name": "Los Angeles International Airport",
				"iataCode": "LAX",
				"point": map[string]float64{
					"latitude": 33.9416,
					"longitude": -118.4085,
				},
			},
			{
				"id": 3,
				"name": "San Francisco International Airport",
				"iataCode": "SFO",
				"point": map[string]float64{
					"latitude": 37.7749,
					"longitude": -122.4194,
				},
			},
			{
				"id": 4,
				"name": "O'Hare International Airport",
				"iataCode": "ORD",
				"point": map[string]float64{
					"latitude": 41.9742,
					"longitude": -87.9073,
				},
			},
		}
		json.NewEncoder(w).Encode(airports)
	})

	log.Printf("Airport Service starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}