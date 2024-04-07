package main

import (
	"log"
	"os"

	"github.com/akashsharma99/passbook-app/internal/initializers"
	"github.com/akashsharma99/passbook-app/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("PASSBOOK_ENV")
	// if the code is running in local environment, load the environment variables from the .env file
	if env == "DEV" {
		log.Println("Running in DEV environment so loading env variables from dev.env file")
		err := godotenv.Load("dev.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} else {
		// the code is running in prod/containerized environment, the environment variables are already set
		log.Println("Running in PROD environment")
	}

	// intialize the database connection pool
	initializers.InitializeDBConnection()

	// initialize the router
	router := routes.NewRouter()
	router.Run(":8080")
	log.Println("Server running on port 8080")
	defer initializers.DB.Close()
}
