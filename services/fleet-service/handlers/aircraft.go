package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/stellora/airline/services/fleet-service/models"
)

// GetAircraft handles GET requests for a specific aircraft
func (h *Handler) GetAircraft(w http.ResponseWriter, r *http.Request) {
	aircraftID, err := getAircraftID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	aircraft, err := h.queries.GetAircraft(r.Context(), aircraftID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Aircraft not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get aircraft")
		return
	}

	respondJSON(w, http.StatusOK, models.FromDBAircraft(aircraft))
}

// GetAircraftByRegistration handles GET requests for aircraft by registration
func (h *Handler) GetAircraftByRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	registration := vars["registration"]

	aircraft, err := h.queries.GetAircraftByRegistration(r.Context(), registration)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Aircraft not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get aircraft")
		return
	}

	respondJSON(w, http.StatusOK, models.FromDBAircraft(aircraft))
}

// ListAircraft handles GET requests for all aircraft
func (h *Handler) ListAircraft(w http.ResponseWriter, r *http.Request) {
	aircraft, err := h.queries.ListAircraft(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list aircraft")
		return
	}

	result := make([]models.Aircraft, 0, len(aircraft))
	for _, a := range aircraft {
		result = append(result, models.FromDBAircraft(a))
	}

	respondJSON(w, http.StatusOK, result)
}

// CreateAircraftRequest represents the request body for creating an aircraft
type CreateAircraftRequest struct {
	Registration string       `json:"registration"`
	AircraftType string       `json:"aircraftType"`
	Airline      models.AirlineSpec `json:"airline"`
}

// CreateAircraft handles POST requests for creating a new aircraft
func (h *Handler) CreateAircraft(w http.ResponseWriter, r *http.Request) {
	var req CreateAircraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Start transaction
	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}
	defer tx.Rollback()
	queriesTx := h.queries.WithTx(tx)

	// Get airline
	// For now, we assume airline service will be called here
	// This is a placeholder until we implement proper service-to-service communication
	airlineID, err := req.Airline.GetAirlineID()
	if err != nil {
		respondError(w, http.StatusBadRequest, "Invalid airline specification")
		return
	}

	// Create aircraft
	created, err := queriesTx.CreateAircraft(r.Context(), CreateAircraftParams{
		Registration: req.Registration,
		AircraftType: req.AircraftType,
		AirlineID:    airlineID,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create aircraft")
		return
	}

	// Get created aircraft
	aircraft, err := queriesTx.GetAircraft(r.Context(), created.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve created aircraft")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	respondJSON(w, http.StatusCreated, models.FromDBAircraft(aircraft))
}

// UpdateAircraftRequest represents the request body for updating an aircraft
type UpdateAircraftRequest struct {
	Registration *string       `json:"registration,omitempty"`
	AircraftType *string       `json:"aircraftType,omitempty"`
	Airline      *models.AirlineSpec `json:"airline,omitempty"`
}

// UpdateAircraft handles PUT requests for updating an aircraft
func (h *Handler) UpdateAircraft(w http.ResponseWriter, r *http.Request) {
	aircraftID, err := getAircraftID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req UpdateAircraftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Start transaction
	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}
	defer tx.Rollback()
	queriesTx := h.queries.WithTx(tx)

	// Verify aircraft exists
	existing, err := queriesTx.GetAircraft(r.Context(), aircraftID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Aircraft not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get aircraft")
		return
	}

	// Prepare update parameters
	params := UpdateAircraftParams{ID: existing.ID}
	if req.Registration != nil {
		params.Registration = sql.NullString{String: *req.Registration, Valid: true}
	}
	if req.AircraftType != nil {
		params.AircraftType = sql.NullString{String: *req.AircraftType, Valid: true}
	}
	if req.Airline != nil {
		airlineID, err := req.Airline.GetAirlineID()
		if err != nil {
			respondError(w, http.StatusBadRequest, "Invalid airline specification")
			return
		}
		params.AirlineID = sql.NullInt64{Int64: airlineID, Valid: true}
	}

	// Update aircraft
	if _, err := queriesTx.UpdateAircraft(r.Context(), params); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update aircraft")
		return
	}

	// Get updated aircraft
	aircraft, err := queriesTx.GetAircraft(r.Context(), existing.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve updated aircraft")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	respondJSON(w, http.StatusOK, models.FromDBAircraft(aircraft))
}

// DeleteAircraft handles DELETE requests for deleting an aircraft
func (h *Handler) DeleteAircraft(w http.ResponseWriter, r *http.Request) {
	aircraftID, err := getAircraftID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Start transaction
	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}
	defer tx.Rollback()
	queriesTx := h.queries.WithTx(tx)

	// Verify aircraft exists
	_, err = queriesTx.GetAircraft(r.Context(), aircraftID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Aircraft not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get aircraft")
		return
	}

	// Delete aircraft
	if err := queriesTx.DeleteAircraft(r.Context(), aircraftID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete aircraft")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}