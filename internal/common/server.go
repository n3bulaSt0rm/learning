package common

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

// GRPCServer wraps grpc.Server with additional functionality
type GRPCServer struct {
	server       *grpc.Server
	healthServer *health.Server
	address      string
}

// NewGRPCServer creates a new gRPC server with health checks and reflection
func NewGRPCServer(address string) *GRPCServer {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)

	// Enable reflection for easier debugging
	reflection.Register(server)

	return &GRPCServer{
		server:       server,
		healthServer: healthServer,
		address:      address,
	}
}

// GetServer returns the underlying grpc.Server
func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}

// SetHealthy marks the service as healthy
func (s *GRPCServer) SetHealthy(service string) {
	s.healthServer.SetServingStatus(service, grpc_health_v1.HealthCheckResponse_SERVING)
}

// SetUnhealthy marks the service as unhealthy
func (s *GRPCServer) SetUnhealthy(service string) {
	s.healthServer.SetServingStatus(service, grpc_health_v1.HealthCheckResponse_NOT_SERVING)
}

// Start starts the gRPC server with graceful shutdown
func (s *GRPCServer) Start() error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	log.Printf("gRPC server starting on %s", s.address)

	// Start server in a goroutine
	go func() {
		if err := s.server.Serve(listener); err != nil {
			log.Printf("gRPC server error: %v", err)
		}
	}()

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for signal
	<-sigChan
	log.Println("Shutting down gRPC server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Println("gRPC server stopped gracefully")
	case <-ctx.Done():
		log.Println("gRPC server shutdown timeout, forcing stop")
		s.server.Stop()
	}

	return nil
}

// Stop stops the server immediately
func (s *GRPCServer) Stop() {
	s.server.Stop()
}

// loggingInterceptor logs gRPC requests
func loggingInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()

	// Call the handler
	resp, err := handler(ctx, req)

	// Log the request
	duration := time.Since(start)
	status := "OK"
	if err != nil {
		status = "ERROR"
	}

	log.Printf("gRPC %s %s %v %s", info.FullMethod, status, duration, getClientIP(ctx))

	return resp, err
}

// getClientIP extracts client IP from context (simplified version)
func getClientIP(ctx context.Context) string {
	// In a real implementation, you would extract this from metadata
	return "unknown"
}
