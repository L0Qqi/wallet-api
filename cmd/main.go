package main

import (
	"log"
	"os"

	"github.com/L0Qqi/wallet-api/internal/handler"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}

	router := handler.SetupRouter()

	log.Printf("Starting server on :%s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
