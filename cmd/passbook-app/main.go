package main

import (
	"log"
	"os"

	"github.com/akashsharma99/passbook-app/internal/initializers"
	"github.com/akashsharma99/passbook-app/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	// load the environment variables
	env := os.Getenv("PASSBOOK_ENV")
	log.Printf("Loaded %s.env variables.", env)
	err := godotenv.Load(env + ".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	// intialize the database connection
	initializers.InitializeDBConnection()

	// initialize the router
	router := routes.NewRouter()
	router.Run(":8080")
	log.Println("Server running on port 8080")
	defer initializers.DB.Close()
}
