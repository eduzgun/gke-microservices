package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/eduzgun/gke-microservices/internal/auth"
	"github.com/eduzgun/gke-microservices/internal/db"
	"github.com/eduzgun/gke-microservices/internal/logger"
	"github.com/eduzgun/gke-microservices/internal/philosopher"
	"github.com/eduzgun/gke-microservices/internal/routes"
	"github.com/eduzgun/gke-microservices/internal/session"
	"github.com/eduzgun/gke-microservices/internal/user"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		panic("couldn't load env file")
	}

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

	// Connect to Postgres db
	pgDB, err := db.NewPostgresDB(log)
	if err != nil {
		log.Error("database initialisation failed")
		os.Exit(1)
	}
	defer pgDB.Close()

	// SQLC queries
	queries := db.New(pgDB.Pool)

	philosopherRepo := philosopher.NewPhilosopherRepo(queries)
	philosopherService := philosopher.NewPhilosopherService(philosopherRepo)
	philosopherController := philosopher.NewPhilosopherController(philosopherService, log)

	userRepo := user.NewUserRepo(queries)
	userService := user.NewUserService(userRepo)

	// Connect to gRPC sesion service
	sessionServiceAddr := os.Getenv("SESSION_SERVICE_ADDR")
	if env == "DEV" {
		sessionServiceAddr = "localhost:9090" // default for dev
	}

	sessionClient, err := session.NewClient(sessionServiceAddr)
	if err != nil {
		log.Error("Failed to create session client", "error", err)
		os.Exit(1)
	}
	defer sessionClient.Close()

	authService := auth.NewAuthService(userService, sessionClient)
	authController := auth.NewAuthController(authService, userService, log)

	// Build router
	router := routes.NewRouter()
	router.RegisterPhilosopherRoutes(philosopherController, sessionClient)
	router.RegisterAuthRoutes(authController, sessionClient)

	// Start server
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Info("Server starting", "address", ":"+port, "environment", env)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Error("Server failed", "error", err)
		os.Exit(1)
	}

}
