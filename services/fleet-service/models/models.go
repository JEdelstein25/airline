package models

import (
	"fmt"
	"strconv"
)

// Aircraft represents an aircraft
type Aircraft struct {
	Id           int      `json:"id"`
	Registration string   `json:"registration"`
	AircraftType string   `json:"aircraftType"`
	Airline      Airline  `json:"airline"`
}

// AircraftView represents an aircraft view from the database
type AircraftView struct {
	ID             int64
	Registration   string
	AircraftType   string
	AirlineID      int64
	AirlineIataCode string
	AirlineName    string
}

// Airline represents an airline
type Airline struct {
	Id       int    `json:"id"`
	IataCode string `json:"iataCode"`
	Name     string `json:"name"`
}

// FromDBAircraft converts a database aircraft view to an API aircraft
func FromDBAircraft(a AircraftView) Aircraft {
	return Aircraft{
		Id:           int(a.ID),
		Registration: a.Registration,
		AircraftType: a.AircraftType,
		Airline: Airline{
			Id:       int(a.AirlineID),
			IataCode: a.AirlineIataCode,
			Name:     a.AirlineName,
		},
	}
}

// Fleet represents a fleet
type Fleet struct {
	Id          int     `json:"id"`
	Airline     Airline `json:"airline"`
	Code        string  `json:"code"`
	Description string  `json:"description"`
}

// FleetsView represents a fleet view from the database
type FleetsView struct {
	ID              int64
	AirlineID       int64
	AirlineIataCode string
	AirlineName     string
	Code            string
	Description     string
}

// FromDBFleet converts a database fleet view to an API fleet
func FromDBFleet(f FleetsView) Fleet {
	return Fleet{
		Id: int(f.ID),
		Airline: Airline{
			Id:       int(f.AirlineID),
			IataCode: f.AirlineIataCode,
			Name:     f.AirlineName,
		},
		Code:        f.Code,
		Description: f.Description,
	}
}

// AircraftSpec represents a specification for finding an aircraft
type AircraftSpec struct {
	Value string `json:"value"` 
}

// AsAircraftID tries to convert the spec to an aircraft ID
func (s AircraftSpec) AsAircraftID() (int64, error) {
	return strconv.ParseInt(s.Value, 10, 64)
}

// AsAircraftRegistration tries to convert the spec to an aircraft registration
func (s AircraftSpec) AsAircraftRegistration() (string, error) {
	if s.Value == "" {
		return "", fmt.Errorf("empty registration")
	}
	return s.Value, nil
}

// AirlineSpec represents a specification for finding an airline
type AirlineSpec struct {
	Value string `json:"value"`
}

// GetAirlineID gets the airline ID from the spec
func (s AirlineSpec) GetAirlineID() (int64, error) {
	return strconv.ParseInt(s.Value, 10, 64)
}

// FleetSpec represents a specification for finding a fleet
type FleetSpec struct {
	Value string `json:"value"`
}

// AsFleetID tries to convert the spec to a fleet ID
func (s FleetSpec) AsFleetID() (int64, error) {
	return strconv.ParseInt(s.Value, 10, 64)
}

// AsFleetCode tries to convert the spec to a fleet code
func (s FleetSpec) AsFleetCode() (string, error) {
	if s.Value == "" {
		return "", fmt.Errorf("empty fleet code")
	}
	return s.Value, nil
}

// AircraftType represents an aircraft type
type AircraftType struct {
	IcaoCode string `json:"icaoCode"`
	Name     string `json:"name"`
}

// AircraftTypes is a list of aircraft types
var AircraftTypes = []AircraftType{
	{IcaoCode: "A318", Name: "Airbus A318"},
	{IcaoCode: "A319", Name: "Airbus A319"},
	{IcaoCode: "A320", Name: "Airbus A320"},
	{IcaoCode: "A321", Name: "Airbus A321"},
	{IcaoCode: "A19N", Name: "Airbus A319neo"},
	{IcaoCode: "A20N", Name: "Airbus A320neo"},
	{IcaoCode: "A21N", Name: "Airbus A321neo"},
	{IcaoCode: "A332", Name: "Airbus A330-200"},
	{IcaoCode: "A338", Name: "Airbus A330-800"},
	{IcaoCode: "A339", Name: "Airbus A330-900"},
	{IcaoCode: "A342", Name: "Airbus A340-200"},
	{IcaoCode: "A343", Name: "Airbus A340-300"},
	{IcaoCode: "A345", Name: "Airbus A340-500"},
	{IcaoCode: "A346", Name: "Airbus A340-600"},
	{IcaoCode: "A359", Name: "Airbus A350-900"},
	{IcaoCode: "A35K", Name: "Airbus A350-1000"},
	{IcaoCode: "A388", Name: "Airbus A380-800"},
	{IcaoCode: "B37M", Name: "Boeing 737 MAX 7"},
	{IcaoCode: "B38M", Name: "Boeing 737 MAX 8"},
	{IcaoCode: "B39M", Name: "Boeing 737 MAX 9"},
	{IcaoCode: "B3XM", Name: "Boeing 737 MAX 10"},
	{IcaoCode: "B712", Name: "Boeing 717"},
	{IcaoCode: "B737", Name: "Boeing 737-700"},
	{IcaoCode: "B738", Name: "Boeing 737-800"},
	{IcaoCode: "B739", Name: "Boeing 737-900"},
	{IcaoCode: "B744", Name: "Boeing 747-400"},
	{IcaoCode: "B748", Name: "Boeing 747-8I"},
	{IcaoCode: "B752", Name: "Boeing 757-200"},
	{IcaoCode: "B753", Name: "Boeing 757-300"},
	{IcaoCode: "B762", Name: "Boeing 767-200"},
	{IcaoCode: "B763", Name: "Boeing 767-300"},
	{IcaoCode: "B764", Name: "Boeing 767-400"},
	{IcaoCode: "B772", Name: "Boeing 777-200"},
	{IcaoCode: "B77L", Name: "Boeing 777-200LR"},
	{IcaoCode: "B773", Name: "Boeing 777-300"},
	{IcaoCode: "B77W", Name: "Boeing 777-300ER"},
	{IcaoCode: "B788", Name: "Boeing 787-8"},
	{IcaoCode: "B789", Name: "Boeing 787-9"},
	{IcaoCode: "B78X", Name: "Boeing 787-10"},
	{IcaoCode: "CRJ2", Name: "Bombardier CRJ-200"},
	{IcaoCode: "CRJ7", Name: "Bombardier CRJ-700"},
	{IcaoCode: "CRJ9", Name: "Bombardier CRJ-900"},
	{IcaoCode: "CRJX", Name: "Bombardier CRJ-1000"},
	{IcaoCode: "E170", Name: "Embraer ERJ-170"},
	{IcaoCode: "E75L", Name: "Embraer ERJ-175"},
	{IcaoCode: "E190", Name: "Embraer ERJ-190"},
	{IcaoCode: "E195", Name: "Embraer ERJ-195"},
	{IcaoCode: "BCS1", Name: "Airbus A220-100"},
	{IcaoCode: "BCS3", Name: "Airbus A220-300"},
	{IcaoCode: "DH8D", Name: "Bombardier Q400"},
	{IcaoCode: "AT72", Name: "ATR 72"},
	{IcaoCode: "AT76", Name: "ATR 72-600"},
	{IcaoCode: "AT75", Name: "ATR 72-500"},
	{IcaoCode: "AT45", Name: "ATR 42-500"},
}