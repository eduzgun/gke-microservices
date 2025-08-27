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
)

func main() {
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

	port := os.Getenv("SESSION_SERVICE_PORT")
	if port == "" {
		port = "9090"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Error("Failed to listen", "error", err)
		os.Exit(1)
	}

	log.Info("Session gRPC server starting", "address", ":"+port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Error("Server failed", "error", err)
		os.Exit(1)
	}
}
