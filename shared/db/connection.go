package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"
	
	_ "github.com/lib/pq"
)

var (
	dbConn *sql.DB
	once    sync.Once
)

// GetConnection returns a database connection singleton
func GetConnection() (*sql.DB, error) {
	var err error
	
	once.Do(func() {
		host := getEnv("DB_HOST", "localhost")
		port := getEnv("DB_PORT", "5432")
		user := getEnv("DB_USER", "airline_user")
		password := getEnv("DB_PASSWORD", "airline_password")
		dbname := getEnv("DB_NAME", "airline_db")

		connStr := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname,
		)
		
		dbConn, err = sql.Open("postgres", connStr)
		if err != nil {
			return
		}
		
		// Set connection pool parameters
		dbConn.SetMaxOpenConns(25)
		dbConn.SetMaxIdleConns(5)
	})
	
	return dbConn, err
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}