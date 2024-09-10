package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"juicedboxx.com/reverse-proxy/pkgs/server"
)

func main() {
	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}
	log.Printf("Current working directory: %s", wd)

	// Print the absolute path of the .env file to verify its existence
	envPath, err := filepath.Abs(".env")
	if err != nil {
		log.Fatalf("Error getting absolute path of .env file: %v", err)
	}
	log.Printf("Absolute path to .env file: %s", envPath)

	// Check if the .env file exists
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		log.Printf(".env file does not exist at path: %s", envPath)
	} else {
		log.Printf(".env file exists at path: %s", envPath)
	}

	err = godotenv.Load(envPath)
	if err != nil {
		log.Default().Printf("Error loading .env file: %v", err)
	} else {
		log.Println(".env file loaded successfully")
	}

	// Use standard HTTP challenge handling
	useStandardHTTPChallengeHandling := true

	// Start the server
	server.StartServer(useStandardHTTPChallengeHandling)
}
