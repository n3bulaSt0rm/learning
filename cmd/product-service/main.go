package main

import (
	"log"

	"learning/internal/common"
	"learning/internal/product"
	pb "learning/pkg/product/pb"
)

func main() {
	// Setup logger
	common.SetupLogger()
	defer common.Close()

	// Load configuration
	config := common.LoadProductServiceConfig()
	log.Printf("Starting Product Service on %s", config.GetGRPCAddress())

	// Initialize repository
	repo := product.NewInMemoryRepository()

	// Initialize service
	service := product.NewService(repo)

	// Initialize gRPC handler
	handler := product.NewHandler(service)

	// Create gRPC server
	server := common.NewGRPCServer(config.GetGRPCAddress())

	// Register service
	pb.RegisterProductServiceServer(server.GetServer(), handler)

	// Set service as healthy
	server.SetHealthy("product")

	// Start server (this will block until shutdown signal)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
