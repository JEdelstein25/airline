package db

// Schema definitions shared across services
// This contains the common database schema that will be used by all microservices
// during the initial migration phase.

// Note: In a later phase, each service will maintain its own schema

// Common types and structures used across services
type AirlineID int64
type AirportID int64
type AircraftID int64
type FleetID int64
type ScheduleID int64
type FlightID int64
type PassengerID int64
type ItineraryID int64

// These structures mirror the database schema
type Airline struct {
	ID       AirlineID `json:"id"`
	IATACode string    `json:"iata_code"`
	Name     string    `json:"name"`
}

type Airport struct {
	ID       AirportID `json:"id"`
	IATACode string    `json:"iata_code"`
	OADBID   *int64    `json:"oadb_id,omitempty"`
}

type Aircraft struct {
	ID           AircraftID `json:"id"`
	Registration string     `json:"registration"`
	AircraftType string     `json:"aircraft_type"`
	AirlineID    AirlineID  `json:"airline_id"`
}

type Fleet struct {
	ID          FleetID   `json:"id"`
	AirlineID   AirlineID `json:"airline_id"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
}

type Schedule struct {
	ID                  ScheduleID `json:"id"`
	AirlineID           AirlineID  `json:"airline_id"`
	Number              string     `json:"number"`
	OriginAirportID     AirportID  `json:"origin_airport_id"`
	DestinationAirportID AirportID  `json:"destination_airport_id"`
	FleetID             FleetID    `json:"fleet_id"`
	StartLocalDate      string     `json:"start_localdate"`
	EndLocalDate        string     `json:"end_localdate"`
	DaysOfWeek          string     `json:"days_of_week"`
	DepartureLocalTime  string     `json:"departure_localtime"`
	DurationSec         int64      `json:"duration_sec"`
	Published           bool       `json:"published"`
}

type Flight struct {
	ID                      FlightID   `json:"id"`
	SourceScheduleID        *ScheduleID `json:"source_schedule_id,omitempty"`
	SourceScheduleInstanceLocalDate *string    `json:"source_schedule_instance_localdate,omitempty"`
	AirlineID               AirlineID  `json:"airline_id"`
	Number                  string     `json:"number"`
	OriginAirportID         AirportID  `json:"origin_airport_id"`
	DestinationAirportID     AirportID  `json:"destination_airport_id"`
	FleetID                 FleetID    `json:"fleet_id"`
	AircraftID              *AircraftID `json:"aircraft_id,omitempty"`
	DepartureDateTime       string     `json:"departure_datetime"`
	ArrivalDateTime         string     `json:"arrival_datetime"`
	DepartureDateTimeUTC    string     `json:"departure_datetime_utc"`
	ArrivalDateTimeUTC      string     `json:"arrival_datetime_utc"`
	Notes                   string     `json:"notes"`
	Published               bool       `json:"published"`
}

type Passenger struct {
	ID   PassengerID `json:"id"`
	Name string      `json:"name"`
}

type Itinerary struct {
	ID      ItineraryID `json:"id"`
	RecordID string      `json:"record_id"`
}

type SeatAssignment struct {
	ID          int64       `json:"id"`
	ItineraryID ItineraryID `json:"itinerary_id"`
	PassengerID PassengerID `json:"passenger_id"`
	FlightID    FlightID    `json:"flight_id"`
	Seat        string      `json:"seat"`
}