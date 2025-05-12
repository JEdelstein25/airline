package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	log.Println("Starting API Gateway...")

	// Create router
	r := mux.NewRouter()

	// Add logging middleware
	r.Use(loggingMiddleware)

	// Set up CORS middleware
	r.Use(corsMiddleware)

	// Health check endpoint
	r.HandleFunc("/health", healthCheckHandler).Methods("GET")
	
	// Service health check endpoints
	r.HandleFunc("/services/health", serviceHealthCheckHandler).Methods("GET")

	// Configure service routes
	configureServiceRoutes(r)

	// 404 handler for paths not handled by microservices
	r.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Start server
	// Create a server that explicitly listens on both IPv4 and IPv6
	addr := "0.0.0.0:" + port
	server := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	log.Printf("API Gateway listening on %s", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// Configure all service routes
func configureServiceRoutes(r *mux.Router) {
	// Define service routes with their base URLs
	services := map[string]string{
		"/api/airlines":       getServiceURL("AIRLINE_SERVICE_URL", "http://airline-service:8001"),
		"/api/aircraft-types": getServiceURL("AIRLINE_SERVICE_URL", "http://airline-service:8001") + "/aircraft-types",
		"/api/airports":      getServiceURL("AIRPORT_SERVICE_URL", "http://airport-service:8002"),
		"/api/fleet":         getServiceURL("FLEET_SERVICE_URL", "http://fleet-service:8003"),
		"/api/schedules":     getServiceURL("SCHEDULE_SERVICE_URL", "http://schedule-service:8004"),
		"/api/bookings":      getServiceURL("BOOKING_SERVICE_URL", "http://booking-service:8005"),
	}

	// Create a proxy for each service
	for path, targetURL := range services {
		path := path     // Create a new variable to avoid closure issues
		targetURL := targetURL
		
		proxy := createReverseProxy(targetURL)
		
		// Register handler for this path
		r.PathPrefix(path).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// For debugging
			log.Printf("[DEBUG] Proxying request from %s to %s", r.URL.Path, targetURL)
			
			// Modify path if needed for specific targets
			if strings.HasPrefix(targetURL, "http://airline-service:8001") && strings.HasPrefix(r.URL.Path, "/api/airlines") {
				// Special handling for airline service
				r.URL.Path = strings.TrimPrefix(r.URL.Path, "/api/airlines")
				if r.URL.Path == "" {
					r.URL.Path = "/airlines"
				}
			} else {
				// Standard path rewriting for other services
				r.URL.Path = strings.TrimPrefix(r.URL.Path, path)
				if r.URL.Path == "" {
					r.URL.Path = "/"
				}
			}
			
			log.Printf("[DEBUG] Rewrote URL path to: %s", r.URL.Path)
			proxy.ServeHTTP(w, r)
		})
	}

	// Add testing routes for direct access to microservices
	r.HandleFunc("/api/test/airlines", func(w http.ResponseWriter, r *http.Request) {
		url := getServiceURL("AIRLINE_SERVICE_URL", "http://airline-service:8001") + "/airlines"
		forwardRequest(w, r, url)
	}).Methods("GET")

	r.HandleFunc("/api/test/aircraft-types", func(w http.ResponseWriter, r *http.Request) {
		url := getServiceURL("AIRLINE_SERVICE_URL", "http://airline-service:8001") + "/aircraft-types"
		forwardRequest(w, r, url)
	}).Methods("GET")
}

// forwardRequest forwards a request to a specific URL and returns the response
func forwardRequest(w http.ResponseWriter, r *http.Request, url string) {
	client := &http.Client{Timeout: 5 * time.Second}
	
	// Create a new request
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	
	// Copy headers from original request
	for name, values := range r.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}
	
	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	defer resp.Body.Close()
	
	// Copy all headers from response
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	
	w.WriteHeader(resp.StatusCode)
	// Copy the response body
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error copying response: %v", err)
	}
}

// createReverseProxy creates a reverse proxy to the target URL
func createReverseProxy(target string) *httputil.ReverseProxy {
	url, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Failed to parse target URL %s: %v", target, err)
	}

	return &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = url.Scheme
			req.URL.Host = url.Host
			
			// Path handling is done by the handler that calls this proxy
			
			// Preserve query parameters
			if url.RawQuery == "" || req.URL.RawQuery == "" {
				req.URL.RawQuery = url.RawQuery + req.URL.RawQuery
			} else {
				req.URL.RawQuery = url.RawQuery + "&" + req.URL.RawQuery
			}
			
			// Set host header to match target host
			req.Host = url.Host
			
			// Set User-Agent if not set
			if _, ok := req.Header["User-Agent"]; !ok {
				req.Header.Set("User-Agent", "")
			}
			
			log.Printf("[DEBUG] Forwarding to: %s://%s%s", req.URL.Scheme, req.URL.Host, req.URL.Path)
		},
		ModifyResponse: func(resp *http.Response) error {
			// Add CORS headers to all responses
			resp.Header.Set("Access-Control-Allow-Origin", "*")
			return nil
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("[ERROR] Proxy error: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadGateway)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Service unavailable: %v", err),
				"path":  r.URL.Path,
			})
		},
	}
}

// corsMiddleware adds CORS headers to responses
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")
		
		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs all requests with their path, method, and timing
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("[REQUEST] %s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		
		// Call the next handler
		next.ServeHTTP(w, r)
		
		log.Printf("[RESPONSE] %s %s %s - completed in %v", r.RemoteAddr, r.Method, r.URL.Path, time.Since(start))
	})
}

// healthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

// serviceHealthCheckHandler checks the health of all microservices
func serviceHealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	serviceStatus := map[string]string{
		"gateway": "healthy",
	}
	
	// Check each service's health by making requests to their health endpoints
	checkServiceHealth(serviceStatus, "fleet", getServiceURL("FLEET_SERVICE_URL", "http://fleet-service:8003"))
	checkServiceHealth(serviceStatus, "airlines", getServiceURL("AIRLINE_SERVICE_URL", "http://airline-service:8001"))
	checkServiceHealth(serviceStatus, "airports", getServiceURL("AIRPORT_SERVICE_URL", "http://airport-service:8002"))
	checkServiceHealth(serviceStatus, "schedules", getServiceURL("SCHEDULE_SERVICE_URL", "http://schedule-service:8004"))
	checkServiceHealth(serviceStatus, "bookings", getServiceURL("BOOKING_SERVICE_URL", "http://booking-service:8005"))
	
	// Return the aggregated health status
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(serviceStatus)
}

// getServiceURL gets the URL for a service from environment or uses the default
func getServiceURL(envVar, defaultURL string) string {
	url := os.Getenv(envVar)
	if url == "" {
		return defaultURL
	}
	return url
}

// checkServiceHealth checks if a service is healthy and updates the status map
func checkServiceHealth(statusMap map[string]string, serviceName, serviceURL string) {
	healthURL := fmt.Sprintf("%s/health", serviceURL)
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	
	resp, err := client.Get(healthURL)
	if err != nil {
		statusMap[serviceName] = fmt.Sprintf("unhealthy: %v", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		statusMap[serviceName] = fmt.Sprintf("unhealthy: status code %d", resp.StatusCode)
		return
	}
	
	statusMap[serviceName] = "healthy"
}

// notFoundHandler handles requests for paths that are not mapped to any service
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Resource not found",
		"path":  r.URL.Path,
	})
}