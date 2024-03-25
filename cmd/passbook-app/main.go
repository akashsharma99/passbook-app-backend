package main

import (
	"log"
	"os"

	"github.com/akashsharma99/passbook-app/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	env := os.Getenv("PASSBOOK_ENV")
	log.Printf("Loaded %s.env variables.", env)
	err := godotenv.Load(env + ".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	router := routes.NewRouter()
	router.Run(":8080")
	// log the port on which the server is running
	log.Println("Server running on port 8080")
}
