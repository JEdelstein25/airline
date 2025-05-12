module github.com/JEdelstein25/airline/services/airline-service

go 1.21

require (
	github.com/go-chi/chi/v5 v5.0.10
	github.com/JEdelstein25/airline/shared v0.0.0-00010101000000-000000000000
	modernc.org/sqlite v1.27.0
)

replace github.com/JEdelstein25/airline/shared => ../../shared