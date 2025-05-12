package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/stellora/airline/services/fleet-service/models"
)

// GetFleet handles GET requests for a specific fleet
func (h *Handler) GetFleet(w http.ResponseWriter, r *http.Request) {
	airlineID, err := getAirlineID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	fleetID, err := getFleetID(r)
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

	// Get fleet
	fleet, err := queriesTx.GetFleet(r.Context(), fleetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Fleet not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get fleet")
		return
	}

	// Verify airline ID matches
	if fleet.AirlineID != airlineID {
		respondError(w, http.StatusNotFound, "Fleet not found for specified airline")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	respondJSON(w, http.StatusOK, models.FromDBFleet(fleet))
}

// GetFleetByCode handles GET requests for a fleet by code
func (h *Handler) GetFleetByCode(w http.ResponseWriter, r *http.Request) {
	airlineID, err := getAirlineID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	vars := mux.Vars(r)
	code := vars["code"]

	// Start transaction
	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to start transaction")
		return
	}
	defer tx.Rollback()
	queriesTx := h.queries.WithTx(tx)

	// Get fleet by code
	fleet, err := queriesTx.GetFleetByCode(r.Context(), GetFleetByCodeParams{
		AirlineID: airlineID,
		Code:      code,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Fleet not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get fleet")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	respondJSON(w, http.StatusOK, models.FromDBFleet(fleet))
}

// ListFleetsByAirline handles GET requests for all fleets of an airline
func (h *Handler) ListFleetsByAirline(w http.ResponseWriter, r *http.Request) {
	airlineID, err := getAirlineID(r)
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

	// Get fleets
	fleets, err := queriesTx.ListFleetsByAirline(r.Context(), airlineID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list fleets")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	result := make([]models.Fleet, 0, len(fleets))
	for _, f := range fleets {
		result = append(result, models.FromDBFleet(f))
	}

	respondJSON(w, http.StatusOK, result)
}

// CreateFleetRequest represents the request body for creating a fleet
type CreateFleetRequest struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// CreateFleet handles POST requests for creating a new fleet
func (h *Handler) CreateFleet(w http.ResponseWriter, r *http.Request) {
	airlineID, err := getAirlineID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req CreateFleetRequest
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

	// Create fleet
	created, err := queriesTx.CreateFleet(r.Context(), CreateFleetParams{
		AirlineID:   airlineID,
		Code:        req.Code,
		Description: req.Description,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create fleet")
		return
	}

	// Get created fleet
	fleet, err := queriesTx.GetFleet(r.Context(), created.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve created fleet")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	respondJSON(w, http.StatusCreated, models.FromDBFleet(fleet))
}

// UpdateFleetRequest represents the request body for updating a fleet
type UpdateFleetRequest struct {
	Code        *string `json:"code,omitempty"`
	Description *string `json:"description,omitempty"`
}

// UpdateFleet handles PUT requests for updating a fleet
func (h *Handler) UpdateFleet(w http.ResponseWriter, r *http.Request) {
	airlineID, err := getAirlineID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	fleetID, err := getFleetID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req UpdateFleetRequest
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

	// Verify fleet exists and belongs to the airline
	existing, err := queriesTx.GetFleet(r.Context(), fleetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Fleet not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get fleet")
		return
	}

	if existing.AirlineID != airlineID {
		respondError(w, http.StatusNotFound, "Fleet not found for specified airline")
		return
	}

	// Prepare update parameters
	params := UpdateFleetParams{ID: existing.ID}
	if req.Code != nil {
		params.Code = sql.NullString{String: *req.Code, Valid: true}
	}
	if req.Description != nil {
		params.Description = sql.NullString{String: *req.Description, Valid: true}
	}

	// Update fleet
	if _, err := queriesTx.UpdateFleet(r.Context(), params); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update fleet")
		return
	}

	// Get updated fleet
	fleet, err := queriesTx.GetFleet(r.Context(), existing.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to retrieve updated fleet")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	respondJSON(w, http.StatusOK, models.FromDBFleet(fleet))
}

// DeleteFleet handles DELETE requests for deleting a fleet
func (h *Handler) DeleteFleet(w http.ResponseWriter, r *http.Request) {
	airlineID, err := getAirlineID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	fleetID, err := getFleetID(r)
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

	// Verify fleet exists and belongs to the airline
	fleet, err := queriesTx.GetFleet(r.Context(), fleetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Fleet not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get fleet")
		return
	}

	if fleet.AirlineID != airlineID {
		respondError(w, http.StatusNotFound, "Fleet not found for specified airline")
		return
	}

	// Delete fleet
	if err := queriesTx.DeleteFleet(r.Context(), fleet.ID); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete fleet")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListAircraftByFleet handles GET requests for all aircraft in a fleet
func (h *Handler) ListAircraftByFleet(w http.ResponseWriter, r *http.Request) {
	airlineID, err := getAirlineID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	fleetID, err := getFleetID(r)
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

	// Verify fleet exists and belongs to the airline
	fleet, err := queriesTx.GetFleet(r.Context(), fleetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Fleet not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get fleet")
		return
	}

	if fleet.AirlineID != airlineID {
		respondError(w, http.StatusNotFound, "Fleet not found for specified airline")
		return
	}

	// Get aircraft in fleet
	aircraft, err := queriesTx.ListAircraftByFleet(r.Context(), fleet.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to list aircraft")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	result := make([]models.Aircraft, 0, len(aircraft))
	for _, a := range aircraft {
		result = append(result, models.FromDBAircraft(a))
	}

	respondJSON(w, http.StatusOK, result)
}

// AddAircraftToFleet handles POST requests for adding an aircraft to a fleet
func (h *Handler) AddAircraftToFleet(w http.ResponseWriter, r *http.Request) {
	airlineID, err := getAirlineID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	fleetID, err := getFleetID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

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

	// Verify fleet exists and belongs to the airline
	fleet, err := queriesTx.GetFleet(r.Context(), fleetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Fleet not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get fleet")
		return
	}

	if fleet.AirlineID != airlineID {
		respondError(w, http.StatusNotFound, "Fleet not found for specified airline")
		return
	}

	// Verify aircraft exists and belongs to the airline
	aircraft, err := queriesTx.GetAircraft(r.Context(), aircraftID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Aircraft not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get aircraft")
		return
	}

	if aircraft.AirlineID != airlineID {
		respondError(w, http.StatusBadRequest, "Aircraft does not belong to the specified airline")
		return
	}

	// Add aircraft to fleet
	err = queriesTx.AddAircraftToFleet(r.Context(), AddAircraftToFleetParams{
		FleetID:    fleet.ID,
		AircraftID: aircraft.ID,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to add aircraft to fleet")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// RemoveAircraftFromFleet handles DELETE requests for removing an aircraft from a fleet
func (h *Handler) RemoveAircraftFromFleet(w http.ResponseWriter, r *http.Request) {
	airlineID, err := getAirlineID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	fleetID, err := getFleetID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

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

	// Verify fleet exists and belongs to the airline
	fleet, err := queriesTx.GetFleet(r.Context(), fleetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Fleet not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get fleet")
		return
	}

	if fleet.AirlineID != airlineID {
		respondError(w, http.StatusNotFound, "Fleet not found for specified airline")
		return
	}

	// Verify aircraft exists and belongs to the airline
	aircraft, err := queriesTx.GetAircraft(r.Context(), aircraftID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusNotFound, "Aircraft not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "Failed to get aircraft")
		return
	}

	if aircraft.AirlineID != airlineID {
		respondError(w, http.StatusBadRequest, "Aircraft does not belong to the specified airline")
		return
	}

	// Remove aircraft from fleet
	err = queriesTx.RemoveAircraftFromFleet(r.Context(), RemoveAircraftFromFleetParams{
		FleetID:    fleet.ID,
		AircraftID: aircraft.ID,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to remove aircraft from fleet")
		return
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to commit transaction")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}