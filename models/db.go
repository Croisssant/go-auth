package models

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func DbInit() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	URL := os.Getenv("TURSO_DATABASE_URL")
	secretKey := os.Getenv("TURSO_AUTH_TOKEN")
	fmt.Printf("URL: %v\nAuth: %v", URL, secretKey)
}
