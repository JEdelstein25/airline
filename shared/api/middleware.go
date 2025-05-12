package api

import (
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs request information
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf(
			"Completed %s %s in %v",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}

// ContentTypeMiddleware ensures JSON content type
func ContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}