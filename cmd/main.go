package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

}
