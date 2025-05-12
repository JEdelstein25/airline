package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite" // SQLite driver
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./airport-service.db"
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize handler
	handler := NewAirportHandler(db)

	// Set up routes
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	})

	// Airport endpoints
	mux.HandleFunc("GET /airports", handler.ListAirports)
	mux.HandleFunc("POST /airports", handler.CreateAirport)
	mux.HandleFunc("GET /airports/{airportSpec}", handler.GetAirport)
	mux.HandleFunc("PATCH /airports/{airportSpec}", handler.UpdateAirport)
	mux.HandleFunc("DELETE /airports/{airportSpec}", handler.DeleteAirport)

	// Start server
	log.Printf("Airport Service starting on port %s...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}

// AirportHandler handles airport-related requests
type AirportHandler struct {
	db *sql.DB
}

// NewAirportHandler creates a new AirportHandler
func NewAirportHandler(db *sql.DB) *AirportHandler {
	return &AirportHandler{db: db}
}

// ListAirports returns a list of all airports
func (h *AirportHandler) ListAirports(w http.ResponseWriter, r *http.Request) {
	// Implementation will be added when extracting actual business logic
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "ListAirports endpoint - to be implemented"}`))
}

// GetAirport returns a specific airport
func (h *AirportHandler) GetAirport(w http.ResponseWriter, r *http.Request) {
	// Implementation will be added when extracting actual business logic
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "GetAirport endpoint - to be implemented"}`))
}

// CreateAirport creates a new airport
func (h *AirportHandler) CreateAirport(w http.ResponseWriter, r *http.Request) {
	// Implementation will be added when extracting actual business logic
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "CreateAirport endpoint - to be implemented"}`))
}

// UpdateAirport updates an existing airport
func (h *AirportHandler) UpdateAirport(w http.ResponseWriter, r *http.Request) {
	// Implementation will be added when extracting actual business logic
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "UpdateAirport endpoint - to be implemented"}`))
}

// DeleteAirport deletes an airport
func (h *AirportHandler) DeleteAirport(w http.ResponseWriter, r *http.Request) {
	// Implementation will be added when extracting actual business logic
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "DeleteAirport endpoint - to be implemented"}`))
}