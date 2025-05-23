package initializers

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PgxPoolIface defines an interface for *pgxpool.Pool to allow for mocking in tests.
// It includes methods used across the application.
type PgxPoolIface interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	Ping(ctx context.Context) error
	Close()
}

var DB PgxPoolIface

// InitializeDBConnection initializes the database connection
func InitializeDBConnection() {
	// get the database connection pool
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("PGSQL_DB_URL"))
	if err != nil {
		log.Fatal("Unable to create a connection pool", err)
	}
	// assign the connection to the global variable
	DB = dbpool
	// log the connection status
	log.Println("Connected to db pool")
}
