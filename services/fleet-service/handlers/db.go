package handlers

import (
	"context"
	"database/sql"

	"github.com/stellora/airline/services/fleet-service/models"
)

// SchemaSQL is the SQL schema for the fleet service
const SchemaSQL = `
CREATE TABLE IF NOT EXISTS aircraft (
    id INTEGER PRIMARY KEY,
    registration TEXT NOT NULL,
    aircraft_type TEXT NOT NULL,
    airline_id INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS fleets (
    id INTEGER PRIMARY KEY,
    airline_id INTEGER NOT NULL,
    code TEXT NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS fleet_aircraft (
    fleet_id INTEGER NOT NULL,
    aircraft_id INTEGER NOT NULL,
    PRIMARY KEY (fleet_id, aircraft_id),
    FOREIGN KEY (fleet_id) REFERENCES fleets(id) ON DELETE CASCADE,
    FOREIGN KEY (aircraft_id) REFERENCES aircraft(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_aircraft_registration ON aircraft(registration);
CREATE INDEX IF NOT EXISTS idx_aircraft_airline_id ON aircraft(airline_id);
CREATE INDEX IF NOT EXISTS idx_fleets_airline_id ON fleets(airline_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_fleets_airline_code ON fleets(airline_id, code);
`

// Queries struct holds all database queries
type Queries struct {
	db *sql.DB
}

// NewQueries creates a new Queries instance
func NewQueries(db *sql.DB) *Queries {
	return &Queries{db: db}
}

// WithTx creates a new Queries instance with a transaction
func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{db: tx}
}

// Aircraft queries

// GetAircraft gets an aircraft by ID
func (q *Queries) GetAircraft(ctx context.Context, id int64) (models.AircraftView, error) {
	const query = `
	SELECT 
		a.id, a.registration, a.aircraft_type, 
		a.airline_id, al.iata_code AS airline_iata_code, al.name AS airline_name
	FROM aircraft a
	LEFT JOIN airlines al ON a.airline_id = al.id
	WHERE a.id = ?
	`

	var aircraft models.AircraftView
	row := q.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&aircraft.ID, &aircraft.Registration, &aircraft.AircraftType,
		&aircraft.AirlineID, &aircraft.AirlineIataCode, &aircraft.AirlineName,
	)
	return aircraft, err
}

// GetAircraftByRegistration gets an aircraft by registration
func (q *Queries) GetAircraftByRegistration(ctx context.Context, registration string) (models.AircraftView, error) {
	const query = `
	SELECT 
		a.id, a.registration, a.aircraft_type, 
		a.airline_id, al.iata_code AS airline_iata_code, al.name AS airline_name
	FROM aircraft a
	LEFT JOIN airlines al ON a.airline_id = al.id
	WHERE a.registration = ?
	`

	var aircraft models.AircraftView
	row := q.db.QueryRowContext(ctx, query, registration)
	err := row.Scan(
		&aircraft.ID, &aircraft.Registration, &aircraft.AircraftType,
		&aircraft.AirlineID, &aircraft.AirlineIataCode, &aircraft.AirlineName,
	)
	return aircraft, err
}

// ListAircraft lists all aircraft
func (q *Queries) ListAircraft(ctx context.Context) ([]models.AircraftView, error) {
	const query = `
	SELECT 
		a.id, a.registration, a.aircraft_type, 
		a.airline_id, al.iata_code AS airline_iata_code, al.name AS airline_name
	FROM aircraft a
	LEFT JOIN airlines al ON a.airline_id = al.id
	ORDER BY a.id
	`

	rows, err := q.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aircraft []models.AircraftView
	for rows.Next() {
		var a models.AircraftView
		if err := rows.Scan(
			&a.ID, &a.Registration, &a.AircraftType,
			&a.AirlineID, &a.AirlineIataCode, &a.AirlineName,
		); err != nil {
			return nil, err
		}
		aircraft = append(aircraft, a)
	}

	return aircraft, rows.Err()
}

// CreateAircraftParams holds parameters for creating an aircraft
type CreateAircraftParams struct {
	Registration string
	AircraftType string
	AirlineID    int64
}

// CreateAircraftResult holds the result of creating an aircraft
type CreateAircraftResult struct {
	ID int64
}

// CreateAircraft creates a new aircraft
func (q *Queries) CreateAircraft(ctx context.Context, params CreateAircraftParams) (CreateAircraftResult, error) {
	const query = `
	INSERT INTO aircraft (registration, aircraft_type, airline_id)
	VALUES (?, ?, ?)
	`

	result, err := q.db.ExecContext(ctx, query, params.Registration, params.AircraftType, params.AirlineID)
	if err != nil {
		return CreateAircraftResult{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return CreateAircraftResult{}, err
	}

	return CreateAircraftResult{ID: id}, nil
}

// UpdateAircraftParams holds parameters for updating an aircraft
type UpdateAircraftParams struct {
	ID           int64
	Registration sql.NullString
	AircraftType sql.NullString
	AirlineID    sql.NullInt64
}

// UpdateAircraft updates an aircraft
func (q *Queries) UpdateAircraft(ctx context.Context, params UpdateAircraftParams) (sql.Result, error) {
	const query = `
	UPDATE aircraft
	SET
		registration = COALESCE(?, registration),
		aircraft_type = COALESCE(?, aircraft_type),
		airline_id = COALESCE(?, airline_id)
	WHERE id = ?
	`

	return q.db.ExecContext(ctx, query, 
		params.Registration.String, params.Registration.Valid,
		params.AircraftType.String, params.AircraftType.Valid,
		params.AirlineID.Int64, params.AirlineID.Valid,
		params.ID)
}

// DeleteAircraft deletes an aircraft
func (q *Queries) DeleteAircraft(ctx context.Context, id int64) error {
	const query = `DELETE FROM aircraft WHERE id = ?`
	_, err := q.db.ExecContext(ctx, query, id)
	return err
}

// DeleteAllAircraft deletes all aircraft
func (q *Queries) DeleteAllAircraft(ctx context.Context) error {
	const query = `DELETE FROM aircraft`
	_, err := q.db.ExecContext(ctx, query)
	return err
}

// Fleet queries

// GetFleet gets a fleet by ID
func (q *Queries) GetFleet(ctx context.Context, id int64) (models.FleetsView, error) {
	const query = `
	SELECT 
		f.id, f.airline_id, al.iata_code AS airline_iata_code, al.name AS airline_name,
		f.code, f.description
	FROM fleets f
	LEFT JOIN airlines al ON f.airline_id = al.id
	WHERE f.id = ?
	`

	var fleet models.FleetsView
	row := q.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&fleet.ID, &fleet.AirlineID, &fleet.AirlineIataCode, &fleet.AirlineName,
		&fleet.Code, &fleet.Description,
	)
	return fleet, err
}

// GetFleetByCodeParams holds parameters for getting a fleet by code
type GetFleetByCodeParams struct {
	AirlineID int64
	Code      string
}

// GetFleetByCode gets a fleet by code for an airline
func (q *Queries) GetFleetByCode(ctx context.Context, params GetFleetByCodeParams) (models.FleetsView, error) {
	const query = `
	SELECT 
		f.id, f.airline_id, al.iata_code AS airline_iata_code, al.name AS airline_name,
		f.code, f.description
	FROM fleets f
	LEFT JOIN airlines al ON f.airline_id = al.id
	WHERE f.airline_id = ? AND f.code = ?
	`

	var fleet models.FleetsView
	row := q.db.QueryRowContext(ctx, query, params.AirlineID, params.Code)
	err := row.Scan(
		&fleet.ID, &fleet.AirlineID, &fleet.AirlineIataCode, &fleet.AirlineName,
		&fleet.Code, &fleet.Description,
	)
	return fleet, err
}

// ListFleetsByAirline lists all fleets for an airline
func (q *Queries) ListFleetsByAirline(ctx context.Context, airlineID int64) ([]models.FleetsView, error) {
	const query = `
	SELECT 
		f.id, f.airline_id, al.iata_code AS airline_iata_code, al.name AS airline_name,
		f.code, f.description
	FROM fleets f
	LEFT JOIN airlines al ON f.airline_id = al.id
	WHERE f.airline_id = ?
	ORDER BY f.code
	`

	rows, err := q.db.QueryContext(ctx, query, airlineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fleets []models.FleetsView
	for rows.Next() {
		var f models.FleetsView
		if err := rows.Scan(
			&f.ID, &f.AirlineID, &f.AirlineIataCode, &f.AirlineName,
			&f.Code, &f.Description,
		); err != nil {
			return nil, err
		}
		fleets = append(fleets, f)
	}

	return fleets, rows.Err()
}

// CreateFleetParams holds parameters for creating a fleet
type CreateFleetParams struct {
	AirlineID   int64
	Code        string
	Description string
}

// CreateFleetResult holds the result of creating a fleet
type CreateFleetResult struct {
	ID int64
}

// CreateFleet creates a new fleet
func (q *Queries) CreateFleet(ctx context.Context, params CreateFleetParams) (CreateFleetResult, error) {
	const query = `
	INSERT INTO fleets (airline_id, code, description)
	VALUES (?, ?, ?)
	`

	result, err := q.db.ExecContext(ctx, query, params.AirlineID, params.Code, params.Description)
	if err != nil {
		return CreateFleetResult{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return CreateFleetResult{}, err
	}

	return CreateFleetResult{ID: id}, nil
}

// UpdateFleetParams holds parameters for updating a fleet
type UpdateFleetParams struct {
	ID          int64
	Code        sql.NullString
	Description sql.NullString
}

// UpdateFleet updates a fleet
func (q *Queries) UpdateFleet(ctx context.Context, params UpdateFleetParams) (sql.Result, error) {
	const query = `
	UPDATE fleets
	SET
		code = COALESCE(?, code),
		description = COALESCE(?, description)
	WHERE id = ?
	`

	return q.db.ExecContext(ctx, query, 
		params.Code.String, params.Code.Valid,
		params.Description.String, params.Description.Valid,
		params.ID)
}

// DeleteFleet deletes a fleet
func (q *Queries) DeleteFleet(ctx context.Context, id int64) error {
	const query = `DELETE FROM fleets WHERE id = ?`
	_, err := q.db.ExecContext(ctx, query, id)
	return err
}

// Fleet aircraft management

// ListAircraftByFleet lists all aircraft in a fleet
func (q *Queries) ListAircraftByFleet(ctx context.Context, fleetID int64) ([]models.AircraftView, error) {
	const query = `
	SELECT 
		a.id, a.registration, a.aircraft_type, 
		a.airline_id, al.iata_code AS airline_iata_code, al.name AS airline_name
	FROM aircraft a
	INNER JOIN fleet_aircraft fa ON a.id = fa.aircraft_id
	LEFT JOIN airlines al ON a.airline_id = al.id
	WHERE fa.fleet_id = ?
	ORDER BY a.registration
	`

	rows, err := q.db.QueryContext(ctx, query, fleetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var aircraft []models.AircraftView
	for rows.Next() {
		var a models.AircraftView
		if err := rows.Scan(
			&a.ID, &a.Registration, &a.AircraftType,
			&a.AirlineID, &a.AirlineIataCode, &a.AirlineName,
		); err != nil {
			return nil, err
		}
		aircraft = append(aircraft, a)
	}

	return aircraft, rows.Err()
}

// AddAircraftToFleetParams holds parameters for adding an aircraft to a fleet
type AddAircraftToFleetParams struct {
	FleetID    int64
	AircraftID int64
}

// AddAircraftToFleet adds an aircraft to a fleet
func (q *Queries) AddAircraftToFleet(ctx context.Context, params AddAircraftToFleetParams) error {
	const query = `
	INSERT INTO fleet_aircraft (fleet_id, aircraft_id)
	VALUES (?, ?)
	ON CONFLICT (fleet_id, aircraft_id) DO NOTHING
	`

	_, err := q.db.ExecContext(ctx, query, params.FleetID, params.AircraftID)
	return err
}

// RemoveAircraftFromFleetParams holds parameters for removing an aircraft from a fleet
type RemoveAircraftFromFleetParams struct {
	FleetID    int64
	AircraftID int64
}

// RemoveAircraftFromFleet removes an aircraft from a fleet
func (q *Queries) RemoveAircraftFromFleet(ctx context.Context, params RemoveAircraftFromFleetParams) error {
	const query = `DELETE FROM fleet_aircraft WHERE fleet_id = ? AND aircraft_id = ?`
	_, err := q.db.ExecContext(ctx, query, params.FleetID, params.AircraftID)
	return err
}