package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/eduzgun/gke-microservices/internal/db"
	"github.com/eduzgun/gke-microservices/internal/logger"
)

func main() {
	// Setup base logger
	var log *slog.Logger
	if os.Getenv("ENV") == "prod" {
		log = logger.NewLogger()
	} else {
		log = logger.NewDevLogger()
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		// Default to dev if not set
		env = "DEV"
	}

	// Load postgres environment
	db, err := db.NewPostgresDB(log)
	if err != nil {
		log.Error("database initialisation failed")
		os.Exit(1)
	}
	defer db.Close()

	fmt.Println("The app has begun")
}
