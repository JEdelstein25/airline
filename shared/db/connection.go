package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite" // SQLite driver
)

// ConnectToDB establishes a database connection with retry logic
func ConnectToDB(dbName string) (*sql.DB, error) {
	// Get database path from environment or use default
	dbPath := os.Getenv(fmt.Sprintf("%s_DB_PATH", dbName))
	if dbPath == "" {
		dbPath = fmt.Sprintf("./%s.db", dbName)
	}

	// Retry connection several times before giving up
	var db *sql.DB
	var err error
	maxRetries := 5

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("sqlite", dbPath)
		if err == nil {
			// Test the connection
			err = db.Ping()
			if err == nil {
				break
			}
		}

		log.Printf("Database connection attempt %d failed: %v, retrying...", i+1, err)
		time.Sleep(time.Second * time.Duration(i+1))
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database after %d attempts: %w", maxRetries, err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * 5)

	log.Printf("Successfully connected to %s database", dbName)
	return db, nil
}

// InitializeSchema ensures the database has the required schema
func InitializeSchema(db *sql.DB, schemaSQL string) error {
	_, err := db.Exec(schemaSQL)
	if err != nil {
		return fmt.Errorf("failed to initialize database schema: %w", err)
	}
	return nil
}