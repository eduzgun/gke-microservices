// cmd/session-service/main.go
package main

import (
	"log/slog"
	"net"
	"os"

	"google.golang.org/grpc"

	"github.com/eduzgun/gke-microservices/internal/logger"
	pb "github.com/eduzgun/gke-microservices/internal/proto/gen/go/session"
	"github.com/eduzgun/gke-microservices/internal/session"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		panic("couldn't load env file")
	}

	var log *slog.Logger
	if os.Getenv("ENV") == "prod" {
		log = logger.NewLogger()
	} else {
		log = logger.NewDevLogger()
	}

	// You could use Redis later; for now, in-memory store
	store := session.NewInMemorySessionStore()
	server := session.NewGRPCSessionServer(store, log)

	grpcServer := grpc.NewServer()
	pb.RegisterSessionServiceServer(grpcServer, server)

	addr := os.Getenv("SESSION_SERVICE_ADDR")
	if addr == "" {
		addr = ":9090" // default to localhost:9090
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("Failed to listen", "error", err)
		os.Exit(1)
	}

	log.Info("Session gRPC server starting", "address", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
