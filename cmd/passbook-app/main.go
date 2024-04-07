package main

import (
	"log"

	"github.com/akashsharma99/passbook-app/internal/initializers"
	"github.com/akashsharma99/passbook-app/internal/routes"
)

func main() {

	// intialize the database connection pool
	initializers.InitializeDBConnection()

	// initialize the router
	router := routes.NewRouter()
	router.Run(":8080")
	log.Println("Server running on port 8080")
	defer initializers.DB.Close()
}
