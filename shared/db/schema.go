package db

// These types represent the core database schema shared across services

// Airline represents an airline entity
type Airline struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	Country string `json:"country"`
}

// Aircraft represents an aircraft entity
type Aircraft struct {
	ID          int    `json:"id"`
	TypeID      int    `json:"type_id"`
	TailNumber  string `json:"tail_number"`
	AirlineID   int    `json:"airline_id"`
	Capacity    int    `json:"capacity"`
	Manufactured string `json:"manufactured"`
}

// Airport represents an airport entity
type Airport struct {
	ID        int     `json:"id"`
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Timezone  string  `json:"timezone"`
}

// Flight represents a flight entity
type Flight struct {
	ID             int    `json:"id"`
	FlightNumber   string `json:"flight_number"`
	AirlineID      int    `json:"airline_id"`
	DepartureTime  string `json:"departure_time"`
	ArrivalTime    string `json:"arrival_time"`
	OriginID       int    `json:"origin_id"`
	DestinationID  int    `json:"destination_id"`
	AircraftID     int    `json:"aircraft_id"`
	Status         string `json:"status"`
}

// Schedule represents a flight schedule entity
type Schedule struct {
	ID           int    `json:"id"`
	FlightNumber string `json:"flight_number"`
	AirlineID    int    `json:"airline_id"`
	OriginID     int    `json:"origin_id"`
	DestinationID int   `json:"destination_id"`
	Departure    string `json:"departure"`
	Arrival      string `json:"arrival"`
	DaysOfWeek   string `json:"days_of_week"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
}

// Passenger represents a passenger entity
type Passenger struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

// SeatAssignment represents a seat assignment entity
type SeatAssignment struct {
	ID          int    `json:"id"`
	FlightID    int    `json:"flight_id"`
	PassengerID int    `json:"passenger_id"`
	SeatNumber  string `json:"seat_number"`
	Class       string `json:"class"`
}