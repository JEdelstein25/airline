module github.com/JEdelstein25/airline/services/fleet-service

go 1.21

require (
	github.com/gorilla/mux v1.8.1
	github.com/JEdelstein25/airline/shared v0.0.0-00010101000000-000000000000
	modernc.org/sqlite v1.28.0
)

replace github.com/JEdelstein25/airline/shared => ../../shared