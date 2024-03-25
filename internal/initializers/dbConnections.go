package initializers

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

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
