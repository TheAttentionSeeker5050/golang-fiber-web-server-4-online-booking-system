package config

import (
	"log"

	"github.com/joho/godotenv"
)

// function to set the from .env file found in the root directory or server configuration
func SetEnvironmentFromFile() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
