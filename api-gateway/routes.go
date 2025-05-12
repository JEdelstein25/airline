package main

import (
	"log"
	"net/http"
	"os"
	"strings"
)

// extendedSetupRoutes will replace the basic setupRoutes function
// as we implement more microservices
func extendedSetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()
	monolithURL := getEnvOrDefault("MONOLITH_URL", "http://localhost:8081")

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	})

	// Routes for Airport Service (when implemented)
	// Example: mux.HandleFunc("GET /airports{$}", forwardToService("airport-service"))
	
	// Routes for Airline Service (when implemented)
	// Example: mux.HandleFunc("GET /airlines{$}", forwardToService("airline-service"))

	// Routes for Fleet Service (when implemented)
	// Example: mux.HandleFunc("GET /aircraft{$}", forwardToService("fleet-service"))
	// Example: mux.HandleFunc("GET /aircraft-types{$}", forwardToService("fleet-service"))

	// Default route - forward to monolith
	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		forwardRequest(monolithURL, w, r)
	})

	return mux
}

// Creates a handler function that forwards requests to the specified service
func forwardToService(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceURL := getServiceURL(serviceName)
		log.Printf("Forwarding request to %s: %s %s", serviceName, r.Method, r.URL.Path)
		forwardRequest(serviceURL, w, r)
	}
}

// Get the URL for a service from environment variables or use default
func getServiceURL(serviceName string) string {
	envKey := strings.ToUpper(strings.ReplaceAll(serviceName, "-", "_")) + "_URL"
	defaultURL := "http://" + serviceName + ":8080"
	return getEnvOrDefault(envKey, defaultURL)
}

// Helper function to get environment variable with default fallback
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}