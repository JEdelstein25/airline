package models

// Passenger represents a passenger in the system
type Passenger struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

// SeatAssignment represents a seat assignment for a passenger on a flight
type SeatAssignment struct {
	ID          int    `json:"id"`
	FlightID    int    `json:"flight_id"`
	PassengerID int    `json:"passenger_id"`
	SeatNumber  string `json:"seat_number"`
}

// ItinerarySegment represents one flight segment in an itinerary
type ItinerarySegment struct {
	ID          int `json:"id"`
	ItineraryID int `json:"itinerary_id"`
	FlightID    int `json:"flight_id"`
}

// Itinerary represents a complete booking itinerary for a passenger
type Itinerary struct {
	ID              int                `json:"id"`
	PassengerID     int                `json:"passenger_id"`
	BookingReference string             `json:"booking_reference"`
	Segments        []ItinerarySegment `json:"segments"`
}