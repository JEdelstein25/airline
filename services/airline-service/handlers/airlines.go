package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/stellora/airline/shared/db"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(database *sql.DB) *Handler {
	return &Handler{db: database}
}

// Airline represents the API response structure
type Airline struct {
	Id       int    `json:"id"`
	IATACode string `json:"iata_code"`
	Name     string `json:"name"`
}

// CreateAirlineRequest represents the create airline request structure
type CreateAirlineRequest struct {
	IATACode string `json:"iata_code"`
	Name     string `json:"name"`
}

// UpdateAirlineRequest represents the update airline request structure
type UpdateAirlineRequest struct {
	IATACode *string `json:"iata_code,omitempty"`
	Name     *string `json:"name,omitempty"`
}

// Helper to convert from DB model to API response
func fromDBAirline(a db.Airline) Airline {
	return Airline{
		Id:       int(a.ID),
		IATACode: a.IATACode,
		Name:     a.Name,
	}
}

// Helper to get airline by ID or IATA code
func getAirlineBySpec(ctx context.Context, dbConn *sql.DB, spec string) (db.Airline, error) {
	// Check if spec is numeric (ID) or string (IATA code)
	id, err := strconv.Atoi(spec)
	if err == nil {
		// Spec is an ID
		row := dbConn.QueryRowContext(ctx, "SELECT id, iata_code, name FROM airlines WHERE id = ?", id)
		var airline db.Airline
		err := row.Scan(&airline.ID, &airline.IATACode, &airline.Name)
		if err != nil {
			return db.Airline{}, err
		}
		return airline, nil
	}

	// Spec is an IATA code
	row := dbConn.QueryRowContext(ctx, "SELECT id, iata_code, name FROM airlines WHERE iata_code = ?", spec)
	var airline db.Airline
	err = row.Scan(&airline.ID, &airline.IATACode, &airline.Name)
	if err != nil {
		return db.Airline{}, err
	}
	return airline, nil
}

// GetAirline handles GET /airlines/{airlineSpec}
func (h *Handler) GetAirline(w http.ResponseWriter, r *http.Request) {
	airlineSpec := chi.URLParam(r, "airlineSpec")
	airline, err := getAirlineBySpec(r.Context(), h.db, airlineSpec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fromDBAirline(airline))
}

// ListAirlines handles GET /airlines
func (h *Handler) ListAirlines(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.QueryContext(r.Context(), "SELECT id, iata_code, name FROM airlines")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var airlines []Airline
	for rows.Next() {
		var airline db.Airline
		if err := rows.Scan(&airline.ID, &airline.IATACode, &airline.Name); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		airlines = append(airlines, fromDBAirline(airline))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(airlines)
}

var validAirlineIATACode = regexp.MustCompile(`^[A-Z0-9]{2}$`)

// CreateAirline handles POST /airlines
func (h *Handler) CreateAirline(w http.ResponseWriter, r *http.Request) {
	var req CreateAirlineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !validAirlineIATACode.MatchString(req.IATACode) {
		http.Error(w, "Invalid IATA code format", http.StatusBadRequest)
		return
	}

	result, err := h.db.ExecContext(
		r.Context(),
		"INSERT INTO airlines (iata_code, name) VALUES (?, ?)",
		req.IATACode, req.Name,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	created := Airline{
		Id:       int(id),
		IATACode: req.IATACode,
		Name:     req.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// UpdateAirline handles PUT /airlines/{airlineSpec}
func (h *Handler) UpdateAirline(w http.ResponseWriter, r *http.Request) {
	airlineSpec := chi.URLParam(r, "airlineSpec")

	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	airline, err := getAirlineBySpec(r.Context(), tx, airlineSpec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var req UpdateAirlineRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Build update query dynamically based on provided fields
	query := "UPDATE airlines SET "
	params := []interface{}{}
	setValues := []

	if req.IATACode != nil {
		if !validAirlineIATACode.MatchString(*req.IATACode) {
			http.Error(w, "Invalid IATA code format", http.StatusBadRequest)
			return
		}
		setValues = append(setValues, "iata_code = ?")
		params = append(params, *req.IATACode)
	}

	if req.Name != nil {
		setValues = append(setValues, "name = ?")
		params = append(params, *req.Name)
	}

	if len(setValues) == 0 {
		// No fields to update
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Build the final query
	for i, setValue := range setValues {
		query += setValue
		if i < len(setValues)-1 {
			query += ", "
		}
	}
	query += " WHERE id = ?"
	params = append(params, airline.ID)

	// Execute the update
	_, err = tx.ExecContext(r.Context(), query, params...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the updated airline
	updated, err := getAirlineBySpec(r.Context(), tx, fmt.Sprintf("%d", airline.ID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fromDBAirline(updated))
}

// DeleteAirline handles DELETE /airlines/{airlineSpec}
func (h *Handler) DeleteAirline(w http.ResponseWriter, r *http.Request) {
	airlineSpec := chi.URLParam(r, "airlineSpec")

	tx, err := h.db.BeginTx(r.Context(), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	airline, err := getAirlineBySpec(r.Context(), tx, airlineSpec)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.ExecContext(r.Context(), "DELETE FROM airlines WHERE id = ?", airline.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteAllAirlines handles DELETE /airlines
func (h *Handler) DeleteAllAirlines(w http.ResponseWriter, r *http.Request) {
	_, err := h.db.ExecContext(r.Context(), "DELETE FROM airlines")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}