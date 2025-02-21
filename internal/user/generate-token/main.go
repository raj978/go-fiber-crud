package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"

	"github.com/raj978/go-fiber-crud/internal/user/generatetoken"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Define your secret key
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}

	// Generate the JWT token
	token, err := generatetoken.GenerateJWT(secretKey)
	if err != nil {
		log.Fatalf("Error generating JWT token: %v", err)
	}

	fmt.Printf("Generated JWT token: %s\n", token)
}
