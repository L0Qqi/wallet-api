package main

import (
	"log"
	"os"

	"github.com/L0Qqi/wallet-api/internal/config"
	"github.com/L0Qqi/wallet-api/internal/handler"
	"github.com/L0Qqi/wallet-api/internal/repository"
	"github.com/L0Qqi/wallet-api/internal/service"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db := config.InitDB()
	defer db.Close()

	repo := repository.NewWalletRepository(db)
	svc := service.NewWalletService(repo)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default port
	}

	router := handler.SetupRouter(svc)

	log.Printf("Starting server on :%s", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
