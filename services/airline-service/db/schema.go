package db

// Schema specific to the Airline service
// During the initial migration phase, this uses the same schema as the monolith

// SQL statements to initialize the database schema
const SchemaSQL = `
CREATE TABLE IF NOT EXISTS airlines (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    iata_code TEXT NOT NULL UNIQUE,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index on iata_code for faster lookups
CREATE INDEX IF NOT EXISTS idx_airlines_iata_code ON airlines(iata_code);
`