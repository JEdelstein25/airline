package main

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"time"

	"github.com/stellora/airline/api-server/api"
	"github.com/stellora/airline/shared/db"
)

// fromDBSchedule converts a database schedule to an API schedule
func fromDBSchedule(a db.SchedulesView) api.Schedule {
	daysOfWeek, _ := parseDaysOfWeek(a.DaysOfWeek)
	airline := fromDBAirline(db.Airline{
		ID:       a.AirlineID,
		IataCode: a.AirlineIataCode,
		Name:     a.AirlineName,
	})
	fleet := fromDBFleet(db.FleetsView{
		ID:          a.FleetID,
		AirlineID:   a.FleetAirlineID,
		Code:        a.FleetCode,
		Description: a.FleetDescription,
	})
	fleet.Airline = airline
	b := api.Schedule{
		Id:      int(a.ID),
		Airline: airline,
		Number:  a.Number,
		OriginAirport: fromDBAirport(db.Airport{
			ID:       a.OriginAirportID,
			IataCode: a.OriginAirportIataCode,
			OadbID:   a.OriginAirportOadbID,
		}),
		DestinationAirport: fromDBAirport(db.Airport{
			ID:       a.DestinationAirportID,
			IataCode: a.DestinationAirportIataCode,
			OadbID:   a.DestinationAirportOadbID,
		}),
		Fleet:         fleet,
		StartDate:     a.StartLocaldate.String(),
		EndDate:       a.EndLocaldate.String(),
		DaysOfWeek:    daysOfWeek,
		DepartureTime: a.DepartureLocaltime.String(),
		DurationSec:   int(a.DurationSec),
		Published:     a.Published,
	}
	b.DistanceMiles = distanceMilesBetweenAirports(b.OriginAirport, b.DestinationAirport)
	return b
}

// fromDBFlight converts a database flight to an API flight
func fromDBFlight(f db.FlightsView) api.Flight {
	airline := fromDBAirline(db.Airline{
		ID:       f.AirlineID,
		IataCode: f.AirlineIataCode,
		Name:     f.AirlineName,
	})
	fleet := fromDBFleet(db.FleetsView{
		ID:          f.FleetID,
		AirlineID:   f.FleetAirlineID,
		Code:        f.FleetCode,
		Description: f.FleetDescription,
	})
	fleet.Airline = airline

	var scheduleID *int
	if f.SourceScheduleID.Valid {
		id := int(f.SourceScheduleID.Int64)
		scheduleID = &id
	}

	var scheduleInstanceDate *string
	if f.SourceScheduleInstanceLocaldate != nil {
		date := f.SourceScheduleInstanceLocaldate.String()
		scheduleInstanceDate = &date
	}

	var aircraft *api.Aircraft
	if f.AircraftID.Valid {
		aircraft = &api.Aircraft{
			Id:           int(f.AircraftID.Int64),
			Registration: f.AircraftRegistration.String,
			AircraftType: f.AircraftTypeIcaoCode.String,
			Airline:      airline,
		}
	}

	flight := api.Flight{
		Id:                   int(f.ID),
		ScheduleID:           scheduleID,
		ScheduleInstanceDate: scheduleInstanceDate,
		Airline:              airline,
		Number:               f.Number,
		OriginAirport: fromDBAirport(db.Airport{
			ID:       f.OriginAirportID,
			IataCode: f.OriginAirportIataCode,
			OadbID:   f.OriginAirportOadbID,
		}),
		DestinationAirport: fromDBAirport(db.Airport{
			ID:       f.DestinationAirportID,
			IataCode: f.DestinationAirportIataCode,
			OadbID:   f.DestinationAirportOadbID,
		}),
		Fleet:             fleet,
		Aircraft:          aircraft,
		DepartureDateTime: *f.DepartureDatetime,
		ArrivalDateTime:   *f.ArrivalDatetime,
		Published:         f.Published,
	}

	flight.DistanceMiles = distanceMilesBetweenAirports(flight.OriginAirport, flight.DestinationAirport)
	return flight
}

// fromDBAirline converts a database airline to an API airline
func fromDBAirline(a db.Airline) api.Airline {
	return api.Airline{
		Id:       int(a.ID),
		IataCode: a.IataCode,
		Name:     a.Name,
	}
}

// fromDBAirport converts a database airport to an API airport
func fromDBAirport(a db.Airport) api.Airport {
	return api.Airport{
		Id:       int(a.ID),
		IataCode: a.IataCode,
	}
}

// fromDBFleet converts a database fleet to an API fleet
func fromDBFleet(f db.FleetsView) api.Fleet {
	return api.Fleet{
		Id:   int(f.ID),
		Code: f.Code,
	}
}

// mapSlice applies a function to each element of a slice
func mapSlice[T any, U any](fn func(T) U, slice []T) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = fn(v)
	}
	return result
}

// getAirlineBySpec gets an airline by its spec (ID or IATA code)
func getAirlineBySpec(ctx context.Context, queries *db.Queries, spec api.AirlineSpec) (db.Airline, error) {
	if spec.Id != 0 {
		return queries.GetAirlineByID(ctx, int64(spec.Id))
	}
	return queries.GetAirlineByIataCode(ctx, spec.IataCode)
}

// getAirportBySpec gets an airport by its spec (ID or IATA code)
func getAirportBySpec(ctx context.Context, queries *db.Queries, spec api.AirportSpec) (db.Airport, error) {
	if spec.Id != 0 {
		return queries.GetAirportByID(ctx, int64(spec.Id))
	}
	return queries.GetAirportByIataCode(ctx, spec.IataCode)
}

// getOrCreateAirportBySpec gets or creates an airport by its spec
func getOrCreateAirportBySpec(ctx context.Context, tx *sql.Tx, queries *db.Queries, spec api.AirportSpec) (db.Airport, error) {
	airport, err := getAirportBySpec(ctx, queries, spec)
	if err == nil || !errors.Is(err, sql.ErrNoRows) {
		return airport, err
	}

	// If the spec has IATA code, create the airport
	if spec.IataCode != "" {
		created, err := queries.CreateAirport(ctx, db.CreateAirportParams{
			IataCode: spec.IataCode,
		})
		if err != nil {
			return db.Airport{}, err
		}
		return queries.GetAirportByID(ctx, created)
	}

	return db.Airport{}, sql.ErrNoRows
}

// getFleetBySpec gets a fleet by its spec (ID or code)
func getFleetBySpec(ctx context.Context, queries *db.Queries, airlineID int64, spec api.FleetSpec) (db.Fleet, error) {
	if spec.Id != 0 {
		return queries.GetFleetByID(ctx, int64(spec.Id))
	}
	return queries.GetFleetByAirlineIDAndCode(ctx, db.GetFleetByAirlineIDAndCodeParams{
		AirlineID: airlineID,
		Code:      spec.Code,
	})
}

// parseDaysOfWeek parses a string like `01356` to a slice with those numbers
func parseDaysOfWeek(str string) (days []int, err error) {
	seen := make([]bool, 7)
	for _, c := range str {
		if c < '0' || c > '6' {
			return nil, errors.New("invalid day of week")
		}
		day := int(c - '0')
		if seen[day] {
			continue
		}
		seen[day] = true
		days = append(days, day)
	}
	sort.Ints(days)
	return days, nil
}

// toDBDaysOfWeek converts a slice of days to a string representation
func toDBDaysOfWeek(days []int) string {
	s := make([]byte, len(days))
	for i, day := range days {
		s[i] = byte('0' + day)
	}
	return string(s)
}

// daysOfWeekContains checks if a day of week is in the days of week slice
func daysOfWeekContains(daysOfWeek []int, day time.Weekday) bool {
	for _, d := range daysOfWeek {
		if d == int(day) {
			return true
		}
	}
	return false
}

// distanceMilesBetweenAirports calculates the distance between two airports
func distanceMilesBetweenAirports(a, b api.Airport) float64 {
	// In a real implementation, this would calculate the actual distance
	// For now, we'll return a placeholder value
	return 0
}