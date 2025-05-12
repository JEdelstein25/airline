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
		port = "8004"
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	// Schedules endpoint with sample flight data
	http.HandleFunc("/schedules", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		schedules := []map[string]interface{}{
			{
				"id": 1,
				"number": "101",
				"published": true,
				"airline": map[string]interface{}{
					"iataCode": "ALS",
					"name": "Airline Sample",
				},
				"originAirport": map[string]interface{}{
					"iataCode": "JFK",
					"name": "John F. Kennedy International Airport",
					"point": map[string]float64{
						"latitude": 40.6413,
						"longitude": -73.7781,
					},
				},
				"destinationAirport": map[string]interface{}{
					"iataCode": "LAX",
					"name": "Los Angeles International Airport",
					"point": map[string]float64{
						"latitude": 33.9416,
						"longitude": -118.4085,
					},
				},
			},
			{
				"id": 2,
				"number": "202",
				"published": true,
				"airline": map[string]interface{}{
					"iataCode": "ALS",
					"name": "Airline Sample",
				},
				"originAirport": map[string]interface{}{
					"iataCode": "SFO",
					"name": "San Francisco International Airport",
					"point": map[string]float64{
						"latitude": 37.7749,
						"longitude": -122.4194,
					},
				},
				"destinationAirport": map[string]interface{}{
					"iataCode": "ORD",
					"name": "O'Hare International Airport",
					"point": map[string]float64{
						"latitude": 41.9742,
						"longitude": -87.9073,
					},
				},
			},
		}
		json.NewEncoder(w).Encode(schedules)
	})

	log.Printf("Schedule Service starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}