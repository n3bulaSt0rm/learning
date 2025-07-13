package main

import (
	"log"

	"learning/internal/common"
	"learning/internal/user"
	pb "learning/pkg/user/pb"
)

func main() {
	// Setup logger
	common.SetupLogger()
	defer common.Close()

	// Load configuration
	config := common.LoadUserServiceConfig()
	log.Printf("Starting User Service on %s", config.GetGRPCAddress())

	// Initialize repository
	repo := user.NewInMemoryRepository()

	// Initialize service
	service := user.NewService(repo)

	// Initialize gRPC handler
	handler := user.NewHandler(service)

	// Create gRPC server
	server := common.NewGRPCServer(config.GetGRPCAddress())

	// Register service
	pb.RegisterUserServiceServer(server.GetServer(), handler)

	// Set service as healthy
	server.SetHealthy("user")

	// Start server (this will block until shutdown signal)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
