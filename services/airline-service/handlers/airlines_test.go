package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	// Use in-memory SQLite for testing
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create the schema
	_, err = db.Exec(`
		CREATE TABLE airlines (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			iata_code TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL
		);
	`)
	if err != nil {
		t.Fatalf("Failed to create test schema: %v", err)
	}

	return db
}

func insertTestAirlines(t *testing.T, db *sql.DB) {
	_, err := db.Exec("INSERT INTO airlines (iata_code, name) VALUES ('XX', 'Test Airline 1')")
	if err != nil {
		t.Fatalf("Failed to insert test airline: %v", err)
	}

	_, err = db.Exec("INSERT INTO airlines (iata_code, name) VALUES ('YY', 'Test Airline 2')")
	if err != nil {
		t.Fatalf("Failed to insert test airline: %v", err)
	}
}

func TestGetAirline(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	insertTestAirlines(t)

	handler := NewHandler(db)

	t.Run("Get by ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/airlines/1", nil)
		rec := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("airlineSpec", "1")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.GetAirline(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var got Airline
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if got.Id != 1 || got.IATACode != "XX" {
			t.Errorf("Expected airline ID 1 with IATA code XX, got %+v", got)
		}
	})

	t.Run("Get by IATA code", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/airlines/XX", nil)
		rec := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("airlineSpec", "XX")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.GetAirline(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
		}

		var got Airline
		if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
			t.Fatalf("Error decoding response: %v", err)
		}

		if got.Id != 1 || got.IATACode != "XX" {
			t.Errorf("Expected airline ID 1 with IATA code XX, got %+v", got)
		}
	})

	t.Run("Not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/airlines/ZZ", nil)
		rec := httptest.NewRecorder()

		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("airlineSpec", "ZZ")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

		handler.GetAirline(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status %d, got %d", http.StatusNotFound, rec.Code)
		}
	})
}

func TestListAirlines(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	insertTestAirlines(t)

	handler := NewHandler(db)

	req := httptest.NewRequest("GET", "/airlines", nil)
	rec := httptest.NewRecorder()

	handler.ListAirlines(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var got []Airline
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if len(got) != 2 {
		t.Errorf("Expected 2 airlines, got %d", len(got))
	}

	if got[0].IATACode != "XX" || got[1].IATACode != "YY" {
		t.Errorf("Expected airlines with IATA codes XX and YY, got %+v", got)
	}
}

func TestCreateAirline(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	handler := NewHandler(db)

	reqBody := bytes.NewBufferString(`{"iata_code":"ZZ","name":"New Airline"}`) 
	req := httptest.NewRequest("POST", "/airlines", reqBody)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.CreateAirline(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var got Airline
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("Error decoding response: %v", err)
	}

	if got.Id != 1 || got.IATACode != "ZZ" || got.Name != "New Airline" {
		t.Errorf("Expected airline ID 1 with IATA code ZZ and name 'New Airline', got %+v", got)
	}

	// Verify in database
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM airlines WHERE iata_code = 'ZZ'").Scan(&count)
	if err != nil {
		t.Fatalf("Database query failed: %v", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 airline in database, got %d", count)
	}
}