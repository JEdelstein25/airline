package api

// API definitions and documentation for the Booking Service

// These annotations can be used for documentation generation
// or for API contract validation

// BookingService API definitions
// @title Booking Service API
// @version 1.0
// @description API for managing passengers, seat assignments, and itineraries
// @basePath /

// PassengerEndpoints defines the passenger management endpoints
// @Summary Passenger management endpoints
type PassengerEndpoints struct {
	// ListPassengers returns all passengers
	// @GET /passengers
	// @Produce json
	// @Success 200 {array} models.Passenger
	// @Failure 500 {object} ErrorResponse
	ListPassengers string

	// GetPassenger returns a specific passenger
	// @GET /passengers/{id}
	// @Param id path int true "Passenger ID"
	// @Produce json
	// @Success 200 {object} models.Passenger
	// @Failure 400 {object} ErrorResponse
	// @Failure 404 {object} ErrorResponse
	// @Failure 500 {object} ErrorResponse
	GetPassenger string

	// CreatePassenger creates a new passenger
	// @POST /passengers
	// @Accept json
	// @Produce json
	// @Param passenger body models.Passenger true "Passenger object"
	// @Success 201 {object} models.Passenger
	// @Failure 400 {object} ErrorResponse
	// @Failure 500 {object} ErrorResponse
	CreatePassenger string

	// UpdatePassenger updates an existing passenger
	// @PUT /passengers/{id}
	// @Param id path int true "Passenger ID"
	// @Accept json
	// @Produce json
	// @Param passenger body models.Passenger true "Passenger object"
	// @Success 200 {object} models.Passenger
	// @Failure 400 {object} ErrorResponse
	// @Failure 500 {object} ErrorResponse
	UpdatePassenger string

	// DeletePassenger removes a passenger
	// @DELETE /passengers/{id}
	// @Param id path int true "Passenger ID"
	// @Produce json
	// @Success 200 {object} SuccessResponse
	// @Failure 400 {object} ErrorResponse
	// @Failure 500 {object} ErrorResponse
	DeletePassenger string
}

// SeatAssignmentEndpoints defines the seat assignment management endpoints
// @Summary Seat assignment management endpoints
type SeatAssignmentEndpoints struct {
	// GetSeatAssignment returns a specific seat assignment
	// @GET /seat-assignments/{id}
	// @Param id path int true "Seat Assignment ID"
	// @Produce json
	// @Success 200 {object} models.SeatAssignment
	// @Failure 400 {object} ErrorResponse
	// @Failure 404 {object} ErrorResponse
	// @Failure 500 {object} ErrorResponse
	GetSeatAssignment string

	// GetSeatAssignmentsForFlight returns all seat assignments for a flight
	// @GET /flights/{id}/seat-assignments
	// @Param id path int true "Flight ID"
	// @Produce json
	// @Success 200 {array} models.SeatAssignment
	// @Failure 400 {object} ErrorResponse
	// @Failure 500 {object} ErrorResponse
	GetSeatAssignmentsForFlight string

	// CreateSeatAssignment creates a new seat assignment
	// @POST /seat-assignments
	// @Accept json
	// @Produce json
	// @Param seatAssignment body models.SeatAssignment true "Seat Assignment object"
	// @Success 201 {object} models.SeatAssignment
	// @Failure 400 {object} ErrorResponse
	// @Failure 500 {object} ErrorResponse
	CreateSeatAssignment string

	// UpdateSeatAssignment updates an existing seat assignment
	// @PUT /seat-assignments/{id}
	// @Param id path int true "Seat Assignment ID"
	// @Accept json
	// @Produce json
	// @Param seatAssignment body models.SeatAssignment true "Seat Assignment object"
	// @Success 200 {object} models.SeatAssignment
	// @Failure 400 {object} ErrorResponse
	// @Failure 500 {object} ErrorResponse
	UpdateSeatAssignment string

	// DeleteSeatAssignment removes a seat assignment
	// @DELETE /seat-assignments/{id}
	// @Param id path int true "Seat Assignment ID"
	// @Produce json
	// @Success 200 {object} SuccessResponse
	// @Failure 400 {object} ErrorResponse
	// @Failure 500 {object} ErrorResponse
	DeleteSeatAssignment string
}

// ItineraryEndpoints defines the itinerary management endpoints
// @Summary Itinerary management endpoints
type ItineraryEndpoints struct {
	// GetItinerary returns a specific itinerary
	// @GET /itineraries/{id}
	// @Param id path int true "Itinerary ID"
	// @Produce json
	// @Success 200 {object} models.Itinerary
	// @Failure 400 {object} ErrorResponse
	// @Failure 404 {object} ErrorResponse
	// @Failure 500 {object} ErrorResponse
	GetItinerary string

	// GetPassengerItineraries returns all itineraries for a passenger
	// @GET /passengers/{id}/itineraries
	// @Param id path int true "Passenger ID"
	// @Produce json
	// @Success 200 {array} models.Itinerary
	// @Failure 400 {object} ErrorResponse
	// @Failure 500 {object} ErrorResponse
	GetPassengerItineraries string

	// CreateItinerary creates a new itinerary
	// @POST /itineraries
	// @Accept json
	// @Produce json
	// @Param itinerary body models.Itinerary true "Itinerary object"
	// @Success 201 {object} models.Itinerary
	// @Failure 400 {object} ErrorResponse
	// @Failure 500 {object} ErrorResponse
	CreateItinerary string

	// DeleteItinerary removes an itinerary
	// @DELETE /itineraries/{id}
	// @Param id path int true "Itinerary ID"
	// @Produce json
	// @Success 200 {object} SuccessResponse
	// @Failure 400 {object} ErrorResponse
	// @Failure 500 {object} ErrorResponse
	DeleteItinerary string
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Result string `json:"result"`
}