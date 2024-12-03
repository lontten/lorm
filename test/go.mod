module test

go 1.21

replace db => ../../lorm

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/jackc/pgx/v5 v5.7.1
	github.com/lontten/lorm v0.2.23
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/gofrs/uuid v4.0.0+incompatible // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.9.0 // indirect
	github.com/pkg/errors v0.8.1 // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/text v0.18.0 // indirect
)
