package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/stellora/airline/services/fleet-service/models"
	"github.com/stellora/airline/shared/db"
)

// Handler handles all aircraft and fleet related requests
type Handler struct {
	db      *sql.DB
	queries *Queries
}

// NewHandler creates a new Handler with dependencies
func NewHandler(db *sql.DB, queries *Queries) *Handler {
	return &Handler{
		db:      db,
		queries: queries,
	}
}

// RegisterRoutes registers all handler routes
func (h *Handler) RegisterRoutes(r *mux.Router) {
	// Aircraft endpoints
	r.HandleFunc("/aircraft", h.ListAircraft).Methods("GET")
	r.HandleFunc("/aircraft", h.CreateAircraft).Methods("POST")
	r.HandleFunc("/aircraft/{id:[0-9]+}", h.GetAircraft).Methods("GET")
	r.HandleFunc("/aircraft/{id:[0-9]+}", h.UpdateAircraft).Methods("PUT")
	r.HandleFunc("/aircraft/{id:[0-9]+}", h.DeleteAircraft).Methods("DELETE")
	r.HandleFunc("/aircraft/registration/{registration}", h.GetAircraftByRegistration).Methods("GET")

	// Aircraft Types endpoints
	r.HandleFunc("/aircraft-types", h.ListAircraftTypes).Methods("GET")

	// Fleet endpoints
	r.HandleFunc("/airlines/{airlineId}/fleets", h.ListFleetsByAirline).Methods("GET")
	r.HandleFunc("/airlines/{airlineId}/fleets", h.CreateFleet).Methods("POST")
	r.HandleFunc("/airlines/{airlineId}/fleets/{fleetId:[0-9]+}", h.GetFleet).Methods("GET")
	r.HandleFunc("/airlines/{airlineId}/fleets/{fleetId:[0-9]+}", h.UpdateFleet).Methods("PUT")
	r.HandleFunc("/airlines/{airlineId}/fleets/{fleetId:[0-9]+}", h.DeleteFleet).Methods("DELETE")
	r.HandleFunc("/airlines/{airlineId}/fleets/code/{code}", h.GetFleetByCode).Methods("GET")

	// Aircraft Fleet Management
	r.HandleFunc("/airlines/{airlineId}/fleets/{fleetId}/aircraft", h.ListAircraftByFleet).Methods("GET")
	r.HandleFunc("/airlines/{airlineId}/fleets/{fleetId}/aircraft/{aircraftId}", h.AddAircraftToFleet).Methods("POST")
	r.HandleFunc("/airlines/{airlineId}/fleets/{fleetId}/aircraft/{aircraftId}", h.RemoveAircraftFromFleet).Methods("DELETE")

	// Health check
	r.HandleFunc("/health", h.HealthCheck).Methods("GET")
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "healthy"})
}

// Helper functions

// getAirlineID extracts and validates airline ID from URL path
func getAirlineID(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	airlineID, err := strconv.ParseInt(vars["airlineId"], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid airline ID: %w", err)
	}
	return airlineID, nil
}

// getFleetID extracts and validates fleet ID from URL path
func getFleetID(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	fleetID, err := strconv.ParseInt(vars["fleetId"], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid fleet ID: %w", err)
	}
	return fleetID, nil
}

// getAircraftID extracts and validates aircraft ID from URL path
func getAircraftID(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	aircraftID, err := strconv.ParseInt(vars["aircraftId"], 10, 64)
	if err != nil {
		// Also check for ID in the regular id param
		aircraftID, err = strconv.ParseInt(vars["id"], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid aircraft ID: %w", err)
		}
	}
	return aircraftID, nil
}

// respondJSON sends a JSON response
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		encodeJSON(w, payload)
	}
}

// encodeJSON encodes payload to JSON
func encodeJSON(w http.ResponseWriter, payload interface{}) {
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// respondError sends an error response
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}