package main

import (
	"log"

	"learning/internal/common"
	"learning/internal/order"
	pb "learning/pkg/order/pb"
)

func main() {
	// Setup logger
	common.SetupLogger()
	defer common.Close()

	// Load configuration
	config := common.LoadOrderServiceConfig()
	log.Printf("Starting Order Service on %s", config.GetGRPCAddress())

	// Initialize external service clients
	userClient, err := order.NewUserServiceClient(config.UserServiceAddress)
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer userClient.Close()

	productClient, err := order.NewProductServiceClient(config.ProductServiceAddress)
	if err != nil {
		log.Fatalf("Failed to connect to product service: %v", err)
	}
	defer productClient.Close()

	// Initialize repository
	repo := order.NewInMemoryRepository()

	// Initialize service with external clients
	service := order.NewService(repo, userClient, productClient)
	defer service.Close()

	// Initialize gRPC handler
	handler := order.NewHandler(service)

	// Create gRPC server
	server := common.NewGRPCServer(config.GetGRPCAddress())

	// Register service
	pb.RegisterOrderServiceServer(server.GetServer(), handler)

	// Set service as healthy
	server.SetHealthy("order")

	// Start server (this will block until shutdown signal)
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
