package db

import (
	"database/sql"
	
	"github.com/JEdelstein25/airline/shared/db"
	"github.com/JEdelstein25/airline/services/booking-service/models"
)

// PassengerRepository handles database operations for passengers
type PassengerRepository struct {
	DB *sql.DB
}

// NewPassengerRepository creates a new passenger repository instance
func NewPassengerRepository(db *sql.DB) *PassengerRepository {
	return &PassengerRepository{DB: db}
}

// GetPassenger retrieves a passenger by ID
func (r *PassengerRepository) GetPassenger(id int) (*models.Passenger, error) {
	var p models.Passenger
	row := r.DB.QueryRow("SELECT id, name, email, phone FROM passengers WHERE id = ?", id)
	err := row.Scan(&p.ID, &p.Name, &p.Email, &p.Phone)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

// ListPassengers returns all passengers
func (r *PassengerRepository) ListPassengers() ([]models.Passenger, error) {
	rows, err := r.DB.Query("SELECT id, name, email, phone FROM passengers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var passengers []models.Passenger
	for rows.Next() {
		var p models.Passenger
		if err := rows.Scan(&p.ID, &p.Name, &p.Email, &p.Phone); err != nil {
			return nil, err
		}
		passengers = append(passengers, p)
	}
	return passengers, nil
}

// CreatePassenger inserts a new passenger
func (r *PassengerRepository) CreatePassenger(p *models.Passenger) error {
	res, err := r.DB.Exec(
		"INSERT INTO passengers (name, email, phone) VALUES (?, ?, ?)",
		p.Name, p.Email, p.Phone,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	p.ID = int(id)
	return nil
}

// UpdatePassenger updates an existing passenger
func (r *PassengerRepository) UpdatePassenger(p *models.Passenger) error {
	_, err := r.DB.Exec(
		"UPDATE passengers SET name = ?, email = ?, phone = ? WHERE id = ?",
		p.Name, p.Email, p.Phone, p.ID,
	)
	return err
}

// DeletePassenger removes a passenger
func (r *PassengerRepository) DeletePassenger(id int) error {
	_, err := r.DB.Exec("DELETE FROM passengers WHERE id = ?", id)
	return err
}

// SeatAssignmentRepository handles database operations for seat assignments
type SeatAssignmentRepository struct {
	DB *sql.DB
}

// NewSeatAssignmentRepository creates a new seat assignment repository instance
func NewSeatAssignmentRepository(db *sql.DB) *SeatAssignmentRepository {
	return &SeatAssignmentRepository{DB: db}
}

// GetSeatAssignment retrieves a seat assignment by ID
func (r *SeatAssignmentRepository) GetSeatAssignment(id int) (*models.SeatAssignment, error) {
	var sa models.SeatAssignment
	row := r.DB.QueryRow(
		"SELECT id, flight_id, passenger_id, seat_number FROM seat_assignments WHERE id = ?", 
		id,
	)
	err := row.Scan(&sa.ID, &sa.FlightID, &sa.PassengerID, &sa.SeatNumber)
	if err != nil {
		return nil, err
	}
	return &sa, nil
}

// GetSeatAssignmentsForFlight retrieves all seat assignments for a flight
func (r *SeatAssignmentRepository) GetSeatAssignmentsForFlight(flightID int) ([]models.SeatAssignment, error) {
	rows, err := r.DB.Query(
		"SELECT id, flight_id, passenger_id, seat_number FROM seat_assignments WHERE flight_id = ?",
		flightID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assignments []models.SeatAssignment
	for rows.Next() {
		var sa models.SeatAssignment
		if err := rows.Scan(&sa.ID, &sa.FlightID, &sa.PassengerID, &sa.SeatNumber); err != nil {
			return nil, err
		}
		assignments = append(assignments, sa)
	}
	return assignments, nil
}

// CreateSeatAssignment creates a new seat assignment
func (r *SeatAssignmentRepository) CreateSeatAssignment(sa *models.SeatAssignment) error {
	res, err := r.DB.Exec(
		"INSERT INTO seat_assignments (flight_id, passenger_id, seat_number) VALUES (?, ?, ?)",
		sa.FlightID, sa.PassengerID, sa.SeatNumber,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	sa.ID = int(id)
	return nil
}

// UpdateSeatAssignment updates an existing seat assignment
func (r *SeatAssignmentRepository) UpdateSeatAssignment(sa *models.SeatAssignment) error {
	_, err := r.DB.Exec(
		"UPDATE seat_assignments SET flight_id = ?, passenger_id = ?, seat_number = ? WHERE id = ?",
		sa.FlightID, sa.PassengerID, sa.SeatNumber, sa.ID,
	)
	return err
}

// DeleteSeatAssignment removes a seat assignment
func (r *SeatAssignmentRepository) DeleteSeatAssignment(id int) error {
	_, err := r.DB.Exec("DELETE FROM seat_assignments WHERE id = ?", id)
	return err
}

// ItineraryRepository handles database operations for itineraries
type ItineraryRepository struct {
	DB *sql.DB
}

// NewItineraryRepository creates a new itinerary repository instance
func NewItineraryRepository(db *sql.DB) *ItineraryRepository {
	return &ItineraryRepository{DB: db}
}

// GetItinerary retrieves an itinerary by ID
func (r *ItineraryRepository) GetItinerary(id int) (*models.Itinerary, error) {
	var it models.Itinerary
	row := r.DB.QueryRow(
		"SELECT id, passenger_id, booking_reference FROM itineraries WHERE id = ?", 
		id,
	)
	err := row.Scan(&it.ID, &it.PassengerID, &it.BookingReference)
	if err != nil {
		return nil, err
	}
	
	// Get itinerary segments
	rows, err := r.DB.Query(
		"SELECT id, itinerary_id, flight_id FROM itinerary_segments WHERE itinerary_id = ?",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var seg models.ItinerarySegment
		if err := rows.Scan(&seg.ID, &seg.ItineraryID, &seg.FlightID); err != nil {
			return nil, err
		}
		it.Segments = append(it.Segments, seg)
	}
	
	return &it, nil
}

// GetPassengerItineraries retrieves all itineraries for a passenger
func (r *ItineraryRepository) GetPassengerItineraries(passengerID int) ([]models.Itinerary, error) {
	rows, err := r.DB.Query(
		"SELECT id, passenger_id, booking_reference FROM itineraries WHERE passenger_id = ?",
		passengerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var itineraries []models.Itinerary
	for rows.Next() {
		var it models.Itinerary
		if err := rows.Scan(&it.ID, &it.PassengerID, &it.BookingReference); err != nil {
			return nil, err
		}
		itineraries = append(itineraries, it)
	}
	
	// Get segments for each itinerary
	for i := range itineraries {
		segRows, err := r.DB.Query(
			"SELECT id, itinerary_id, flight_id FROM itinerary_segments WHERE itinerary_id = ?",
			itineraries[i].ID,
		)
		if err != nil {
			return nil, err
		}
		defer segRows.Close()

		for segRows.Next() {
			var seg models.ItinerarySegment
			if err := segRows.Scan(&seg.ID, &seg.ItineraryID, &seg.FlightID); err != nil {
				return nil, err
			}
			itineraries[i].Segments = append(itineraries[i].Segments, seg)
		}
	}
	
	return itineraries, nil
}

// CreateItinerary creates a new itinerary
func (r *ItineraryRepository) CreateItinerary(it *models.Itinerary) error {
	// Begin transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Insert itinerary
	res, err := tx.Exec(
		"INSERT INTO itineraries (passenger_id, booking_reference) VALUES (?, ?)",
		it.PassengerID, it.BookingReference,
	)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	it.ID = int(id)

	// Insert itinerary segments
	for i := range it.Segments {
		it.Segments[i].ItineraryID = it.ID
		res, err := tx.Exec(
			"INSERT INTO itinerary_segments (itinerary_id, flight_id) VALUES (?, ?)",
			it.Segments[i].ItineraryID, it.Segments[i].FlightID,
		)
		if err != nil {
			return err
		}
		
		segID, err := res.LastInsertId()
		if err != nil {
			return err
		}
		it.Segments[i].ID = int(segID)
	}

	return nil
}

// DeleteItinerary removes an itinerary and its segments
func (r *ItineraryRepository) DeleteItinerary(id int) error {
	// Begin transaction
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Delete segments first (foreign key constraint)
	_, err = tx.Exec("DELETE FROM itinerary_segments WHERE itinerary_id = ?", id)
	if err != nil {
		return err
	}

	// Delete itinerary
	_, err = tx.Exec("DELETE FROM itineraries WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}