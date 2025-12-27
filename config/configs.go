package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// GetDatabaseURL constructs the database URL from environment variables
var DATABASEURL string

// add defauklt images url github user images Rtarun3606k
var DEFAULTIMAGEURL = "https://avatars.githubusercontent.com/u/97576326?v=4"

// all collection in modngodb v2
var Collection = []string{"projects", "positions", "services", "blogs"}

// LoadEnv loads environment variables from a .env file
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}
	log.Println("Environment variables loaded successfully")
	DATABASEURL = getEnv("DATABASEURL", "mongodb://localhost:27017")

}

// GetEnv retrieves the value of the environment variable named by the key.
// It returns the value, or fallback if the variable is not present.
func getEnv(key, fallback string) string {
	exists := os.Getenv(key)
	// log.Println("Environment variable", key, "=", exists)
	if exists == "" {
		return fallback
	}
	return exists
}
