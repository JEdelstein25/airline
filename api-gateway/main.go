package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := setupRoutes()

	log.Printf("API Gateway starting on port %s...", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), mux))
}

func setupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"ok":true}`))
	})

	// Initially route all requests to the monolith
	// This will be replaced with service-specific routes as we migrate
	mux.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		monolithURL := os.Getenv("MONOLITH_URL")
		if monolithURL == "" {
			monolithURL = "http://localhost:8081"
		}
		
		// Forward request to monolith
		forwardRequest(monolithURL, w, r)
	})

	return mux
}

func forwardRequest(targetURL string, w http.ResponseWriter, r *http.Request) {
	// Simple proxy implementation
	// This will be expanded as we implement more sophisticated routing
	client := &http.Client{}
	
	// Create a new request to forward
	req, err := http.NewRequest(r.Method, targetURL+r.URL.Path, r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}
	
	// Copy headers
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}
	
	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error forwarding request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	
	// Copy response headers
	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	
	// Set status code
	w.WriteHeader(resp.StatusCode)
	
	// Copy response body
	buf := make([]byte, 32*1024) // 32k buffer
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			w.Write(buf[:n])
		}
		if err != nil {
			break
		}
	}
}