package connections

import (
	"avito_test_fix/utilities"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

var (
	PostgresConn *pgxpool.Pool
	PostgresErr  error
)

const (
	ErrNoRowsInResult = "no rows in result set"
)

func InitPostgresConnection() {
	host, err := os.LookupEnv("POSTGRES_HOST")
	if !err {
		utilities.LogStrErr("POSTGRES-PANIC", "POSTGRES_HOST not found in .env file")
	}

	port, err := os.LookupEnv("POSTGRES_PORT")
	if !err {
		utilities.LogStrErr("POSTGRES-PANIC", "POSTGRES_PORT not found in .env file")
	}

	username, err := os.LookupEnv("POSTGRES_USERNAME")
	if !err {
		utilities.LogStrErr("POSTGRES-PANIC", "POSTGRES_USERNAME not found in .env file")
	}

	password, err := os.LookupEnv("POSTGRES_PASSWORD")
	if !err {
		utilities.LogStrErr("POSTGRES-PANIC", "POSTGRES_PASSWORD not found in .env file")
	}

	database, err := os.LookupEnv("POSTGRES_DB")
	if !err {
		utilities.LogStrErr("POSTGRES-PANIC", "POSTGRES_DB not found in .env file")
	}

	PostgresConn, PostgresErr = pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, database))
	utilities.PanicIfErr(PostgresErr)

	fmt.Printf("\nConnect to Postgres was successfully")
}
