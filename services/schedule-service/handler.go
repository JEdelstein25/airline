package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/stellora/airline/api-server/api"
	"github.com/stellora/airline/api-server/localtime"
	"github.com/stellora/airline/shared/db"
)

// Handler handles HTTP requests for the schedule service
type Handler struct {
	db      *sql.DB
	queries *db.Queries
}

// NewHandler creates a new handler with the given database connection and queries
func NewHandler(db *sql.DB, queries *db.Queries) *Handler {
	return &Handler{
		db:      db,
		queries: queries,
	}
}

// listSchedules handles GET /schedules
func (h *Handler) listSchedules(w http.ResponseWriter, r *http.Request) {
	schedules, err := h.queries.ListSchedules(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := mapSlice(fromDBSchedule, schedules)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getSchedule handles GET /schedules/{id}
func (h *Handler) getSchedule(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
		return
	}

	schedule, err := h.queries.GetSchedule(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Schedule not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := fromDBSchedule(schedule)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// createSchedule handles POST /schedules
func (h *Handler) createSchedule(w http.ResponseWriter, r *http.Request) {
	var body api.CreateScheduleJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if body.Number == "" {
		http.Error(w, "number must not be empty", http.StatusBadRequest)
		return
	}

	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	queriesTx := h.queries.WithTx(tx)

	airline, err := getAirlineBySpec(r.Context(), queriesTx, body.Airline)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, fmt.Sprintf("airline %q not found", body.Airline), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	originAirport, err := getOrCreateAirportBySpec(r.Context(), tx, queriesTx, body.OriginAirport)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, fmt.Sprintf("originAirport %q not found", body.OriginAirport), http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("looking up originAirport: %v", err), http.StatusInternalServerError)
		return
	}
	
	destinationAirport, err := getOrCreateAirportBySpec(r.Context(), tx, queriesTx, body.DestinationAirport)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, fmt.Sprintf("destinationAirport %q not found", body.DestinationAirport), http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("looking up destinationAirport: %v", err), http.StatusInternalServerError)
		return
	}

	fleet, err := getFleetBySpec(r.Context(), queriesTx, airline.ID, body.Fleet)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, fmt.Sprintf("fleet %q not found", body.Fleet), http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("looking up fleet: %v", err), http.StatusInternalServerError)
		return
	}

	startDate, err := localtime.ParseLocalDate(body.StartDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("parsing startDate: %v", err), http.StatusBadRequest)
		return
	}
	
	endDate, err := localtime.ParseLocalDate(body.EndDate)
	if err != nil {
		http.Error(w, fmt.Sprintf("parsing endDate: %v", err), http.StatusBadRequest)
		return
	}

	departureTime, err := localtime.ParseTimeOfDay(body.DepartureTime)
	if err != nil {
		http.Error(w, fmt.Sprintf("parsing departureTime: %v", err), http.StatusBadRequest)
		return
	}

	created, err := queriesTx.CreateSchedule(r.Context(), db.CreateScheduleParams{
		AirlineID:            airline.ID,
		Number:               body.Number,
		OriginAirportID:      originAirport.ID,
		DestinationAirportID: destinationAirport.ID,
		FleetID:              fleet.ID,
		StartLocaldate:       &startDate,
		EndLocaldate:         &endDate,
		DaysOfWeek:           toDBDaysOfWeek(body.DaysOfWeek),
		DepartureLocaltime:   &departureTime,
		DurationSec:          int64(body.DurationSec),
		Published:            body.Published != nil && *body.Published,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	schedule, err := queriesTx.GetSchedule(r.Context(), created)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := syncScheduleFlightInstances(r.Context(), queriesTx, schedule); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := fromDBSchedule(schedule)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// updateSchedule handles PATCH /schedules/{id}
func (h *Handler) updateSchedule(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
		return
	}

	var body api.UpdateScheduleJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	queriesTx := h.queries.WithTx(tx)

	existing, err := queriesTx.GetSchedule(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Schedule not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := db.UpdateScheduleParams{
		ID: id,
	}

	if body.Number != nil {
		params.Number = sql.NullString{String: *body.Number, Valid: true}
	}
	
	if body.OriginAirport != nil {
		originAirport, err := getOrCreateAirportBySpec(r.Context(), tx, queriesTx, *body.OriginAirport)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		params.OriginAirportID = sql.NullInt64{Int64: originAirport.ID, Valid: true}
	}
	
	if body.DestinationAirport != nil {
		destinationAirport, err := getOrCreateAirportBySpec(r.Context(), tx, queriesTx, *body.DestinationAirport)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		params.DestinationAirportID = sql.NullInt64{Int64: destinationAirport.ID, Valid: true}
	}
	
	if body.Fleet != nil {
		fleet, err := getFleetBySpec(r.Context(), queriesTx, existing.AirlineID, *body.Fleet)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, fmt.Sprintf("fleet %q not found", *body.Fleet), http.StatusBadRequest)
				return
			}
			http.Error(w, fmt.Sprintf("looking up fleet: %v", err), http.StatusInternalServerError)
			return
		}
		params.FleetID = sql.NullInt64{Int64: fleet.ID, Valid: true}
	}
	
	if body.StartDate != nil {
		startDate, err := localtime.ParseLocalDate(*body.StartDate)
		if err != nil {
			http.Error(w, fmt.Sprintf("parsing startDate: %v", err), http.StatusBadRequest)
			return
		}
		params.StartLocaldate = &startDate
	}
	
	if body.EndDate != nil {
		endDate, err := localtime.ParseLocalDate(*body.EndDate)
		if err != nil {
			http.Error(w, fmt.Sprintf("parsing endDate: %v", err), http.StatusBadRequest)
			return
		}
		params.EndLocaldate = &endDate
	}
	
	if body.DaysOfWeek != nil {
		params.DaysOfWeek = sql.NullString{String: toDBDaysOfWeek(*body.DaysOfWeek), Valid: true}
	}
	
	if body.DepartureTime != nil {
		departureTime, err := localtime.ParseTimeOfDay(*body.DepartureTime)
		if err != nil {
			http.Error(w, fmt.Sprintf("parsing departureTime: %v", err), http.StatusBadRequest)
			return
		}
		params.DepartureLocaltime = &departureTime
	}
	
	if body.DurationSec != nil {
		params.DurationSec = sql.NullInt64{Int64: int64(*body.DurationSec), Valid: true}
	}
	
	if body.Published != nil {
		params.Published = sql.NullBool{Bool: *body.Published, Valid: true}
	}

	if _, err := queriesTx.UpdateSchedule(r.Context(), params); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Schedule not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	schedule, err := queriesTx.GetSchedule(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Schedule not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := syncScheduleFlightInstances(r.Context(), queriesTx, schedule); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := fromDBSchedule(schedule)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// deleteSchedule handles DELETE /schedules/{id}
func (h *Handler) deleteSchedule(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
		return
	}

	if err := h.queries.DeleteSchedule(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// deleteAllSchedules handles DELETE /schedules
func (h *Handler) deleteAllSchedules(w http.ResponseWriter, r *http.Request) {
	if err := h.queries.DeleteAllSchedules(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// listFlightsForSchedule handles GET /schedules/{id}/flights
func (h *Handler) listFlightsForSchedule(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid schedule ID", http.StatusBadRequest)
		return
	}

	rows, err := h.queries.ListFlightsForSchedule(r.Context(), sql.NullInt64{Valid: true, Int64: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := mapSlice(fromDBFlight, rows)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}