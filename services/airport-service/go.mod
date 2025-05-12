module github.com/JEdelstein25/airline/services/airport-service

go 1.21

require (
	github.com/gorilla/mux v1.8.1
	github.com/lib/pq v1.10.9
	github.com/JEdelstein25/airline/shared v0.0.0-00010101000000-000000000000
)

replace github.com/JEdelstein25/airline/shared => ../../shared