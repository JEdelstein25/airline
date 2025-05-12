package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/JEdelstein25/airline/services/booking-service/db"
	"github.com/JEdelstein25/airline/services/booking-service/models"
)

// Handler contains all the dependencies needed for the handlers
type Handler struct {
	PassengerRepo     *db.PassengerRepository
	SeatAssignmentRepo *db.SeatAssignmentRepository
	ItineraryRepo     *db.ItineraryRepository
}

// NewHandler creates a new handler instance
func NewHandler(dbConn *sql.DB) *Handler {
	return &Handler{
		PassengerRepo:     db.NewPassengerRepository(dbConn),
		SeatAssignmentRepo: db.NewSeatAssignmentRepository(dbConn),
		ItineraryRepo:     db.NewItineraryRepository(dbConn),
	}
}

// RegisterRoutes registers all API routes
func RegisterRoutes(r *mux.Router, dbConn *sql.DB) {
	h := NewHandler(dbConn)

	// Passenger routes
	r.HandleFunc("/passengers", h.ListPassengers).Methods("GET")
	r.HandleFunc("/passengers/{id}", h.GetPassenger).Methods("GET")
	r.HandleFunc("/passengers", h.CreatePassenger).Methods("POST")
	r.HandleFunc("/passengers/{id}", h.UpdatePassenger).Methods("PUT")
	r.HandleFunc("/passengers/{id}", h.DeletePassenger).Methods("DELETE")

	// Seat assignment routes
	r.HandleFunc("/seat-assignments/{id}", h.GetSeatAssignment).Methods("GET")
	r.HandleFunc("/flights/{id}/seat-assignments", h.GetSeatAssignmentsForFlight).Methods("GET")
	r.HandleFunc("/seat-assignments", h.CreateSeatAssignment).Methods("POST")
	r.HandleFunc("/seat-assignments/{id}", h.UpdateSeatAssignment).Methods("PUT")
	r.HandleFunc("/seat-assignments/{id}", h.DeleteSeatAssignment).Methods("DELETE")

	// Itinerary routes
	r.HandleFunc("/itineraries/{id}", h.GetItinerary).Methods("GET")
	r.HandleFunc("/passengers/{id}/itineraries", h.GetPassengerItineraries).Methods("GET")
	r.HandleFunc("/itineraries", h.CreateItinerary).Methods("POST")
	r.HandleFunc("/itineraries/{id}", h.DeleteItinerary).Methods("DELETE")
}

// Helper functions
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Failed to marshal JSON response"}`)) 
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// Passenger Handlers

// ListPassengers returns all passengers
func (h *Handler) ListPassengers(w http.ResponseWriter, r *http.Request) {
	passengers, err := h.PassengerRepo.ListPassengers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve passengers")
		return
	}

	respondWithJSON(w, http.StatusOK, passengers)
}

// GetPassenger returns a specific passenger
func (h *Handler) GetPassenger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid passenger ID")
		return
	}

	passenger, err := h.PassengerRepo.GetPassenger(id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Passenger not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve passenger")
		return
	}

	respondWithJSON(w, http.StatusOK, passenger)
}

// CreatePassenger creates a new passenger
func (h *Handler) CreatePassenger(w http.ResponseWriter, r *http.Request) {
	var passenger models.Passenger
	if err := json.NewDecoder(r.Body).Decode(&passenger); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.PassengerRepo.CreatePassenger(&passenger); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create passenger")
		return
	}

	respondWithJSON(w, http.StatusCreated, passenger)
}

// UpdatePassenger updates an existing passenger
func (h *Handler) UpdatePassenger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid passenger ID")
		return
	}

	var passenger models.Passenger
	if err := json.NewDecoder(r.Body).Decode(&passenger); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	passenger.ID = id
	if err := h.PassengerRepo.UpdatePassenger(&passenger); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update passenger")
		return
	}

	respondWithJSON(w, http.StatusOK, passenger)
}

// DeletePassenger removes a passenger
func (h *Handler) DeletePassenger(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid passenger ID")
		return
	}

	if err := h.PassengerRepo.DeletePassenger(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete passenger")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// Seat Assignment Handlers

// GetSeatAssignment returns a specific seat assignment
func (h *Handler) GetSeatAssignment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid seat assignment ID")
		return
	}

	assignment, err := h.SeatAssignmentRepo.GetSeatAssignment(id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Seat assignment not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve seat assignment")
		return
	}

	respondWithJSON(w, http.StatusOK, assignment)
}

// GetSeatAssignmentsForFlight returns all seat assignments for a flight
func (h *Handler) GetSeatAssignmentsForFlight(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid flight ID")
		return
	}

	assignments, err := h.SeatAssignmentRepo.GetSeatAssignmentsForFlight(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve seat assignments")
		return
	}

	respondWithJSON(w, http.StatusOK, assignments)
}

// CreateSeatAssignment creates a new seat assignment
func (h *Handler) CreateSeatAssignment(w http.ResponseWriter, r *http.Request) {
	var assignment models.SeatAssignment
	if err := json.NewDecoder(r.Body).Decode(&assignment); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.SeatAssignmentRepo.CreateSeatAssignment(&assignment); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create seat assignment")
		return
	}

	respondWithJSON(w, http.StatusCreated, assignment)
}

// UpdateSeatAssignment updates an existing seat assignment
func (h *Handler) UpdateSeatAssignment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid seat assignment ID")
		return
	}

	var assignment models.SeatAssignment
	if err := json.NewDecoder(r.Body).Decode(&assignment); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	assignment.ID = id
	if err := h.SeatAssignmentRepo.UpdateSeatAssignment(&assignment); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update seat assignment")
		return
	}

	respondWithJSON(w, http.StatusOK, assignment)
}

// DeleteSeatAssignment removes a seat assignment
func (h *Handler) DeleteSeatAssignment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid seat assignment ID")
		return
	}

	if err := h.SeatAssignmentRepo.DeleteSeatAssignment(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete seat assignment")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// Itinerary Handlers

// GetItinerary returns a specific itinerary
func (h *Handler) GetItinerary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid itinerary ID")
		return
	}

	itinerary, err := h.ItineraryRepo.GetItinerary(id)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Itinerary not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve itinerary")
		return
	}

	respondWithJSON(w, http.StatusOK, itinerary)
}

// GetPassengerItineraries returns all itineraries for a passenger
func (h *Handler) GetPassengerItineraries(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid passenger ID")
		return
	}

	itineraries, err := h.ItineraryRepo.GetPassengerItineraries(id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve itineraries")
		return
	}

	respondWithJSON(w, http.StatusOK, itineraries)
}

// CreateItinerary creates a new itinerary
func (h *Handler) CreateItinerary(w http.ResponseWriter, r *http.Request) {
	var itinerary models.Itinerary
	if err := json.NewDecoder(r.Body).Decode(&itinerary); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.ItineraryRepo.CreateItinerary(&itinerary); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create itinerary")
		return
	}

	respondWithJSON(w, http.StatusCreated, itinerary)
}

// DeleteItinerary removes an itinerary
func (h *Handler) DeleteItinerary(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid itinerary ID")
		return
	}

	if err := h.ItineraryRepo.DeleteItinerary(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete itinerary")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}