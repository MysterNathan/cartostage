module student

go 1.25.0

replace shared => ../../shared

require (
	github.com/gorilla/mux v1.8.1
	github.com/jmoiron/sqlx v1.4.0
	shared v0.0.0-00010101000000-000000000000
)

require (
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.42.0 // indirect
)
