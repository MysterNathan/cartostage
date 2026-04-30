module enterprise

go 1.25.0

replace shared => ../../shared

require (
	github.com/gorilla/mux v1.8.1
	github.com/jmoiron/sqlx v1.4.0
	shared v0.0.0-00010101000000-000000000000
)

require (
	github.com/gabriel-vasile/mimetype v1.4.12 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.30.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.46.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
)
