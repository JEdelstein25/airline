package handlers

import (
	"net/http"

	"github.com/stellora/airline/services/fleet-service/models"
)

// ListAircraftTypes handles GET requests for aircraft types
func (h *Handler) ListAircraftTypes(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, models.AircraftTypes)
}