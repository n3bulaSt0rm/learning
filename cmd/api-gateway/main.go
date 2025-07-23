package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"learning/internal/common"
	graphqlserver "learning/internal/graphql"
	"learning/internal/order"
	"learning/internal/product"
	"learning/internal/user"
	orderpb "learning/pkg/order/pb"
	productpb "learning/pkg/product/pb"
	userpb "learning/pkg/user/pb"
)

func main() {
	// Setup logger
	common.SetupLogger()
	defer common.Close()

	// Load configuration
	config := common.LoadGatewayConfig()
	log.Printf("Starting API Gateway on %s", config.GetHTTPAddress())

	// Create context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create gRPC-Gateway mux
	mux := runtime.NewServeMux()

	// Setup service connections
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register User Service
	err := userpb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, config.UserServiceAddress, opts)
	if err != nil {
		log.Fatalf("Failed to register user service handler: %v", err)
	}
	log.Printf("Registered User  Service proxy to %s", config.UserServiceAddress)

	// Register Product Service
	err = productpb.RegisterProductServiceHandlerFromEndpoint(ctx, mux, config.ProductServiceAddress, opts)
	if err != nil {
		log.Fatalf("Failed to register product service handler: %v", err)
	}
	log.Printf("Registered Product Service proxy to %s", config.ProductServiceAddress)

	// Register Order Service
	err = orderpb.RegisterOrderServiceHandlerFromEndpoint(ctx, mux, config.OrderServiceAddress, opts)
	if err != nil {
		log.Fatalf("Failed to register order service handler: %v", err)
	}
	log.Printf("Registered Order Service proxy to %s", config.OrderServiceAddress)

	// Create a main mux that handles API, GraphQL and health endpoints
	mainMux := http.NewServeMux()

	// Add health check endpoint
	mainMux.HandleFunc("/health", healthCheckHandler)

	// Add gRPC-Gateway routes under /api/
	mainMux.Handle("/api/", mux)

	// Setup GraphQL if enabled
	if config.GraphQLEnabled {
		// Create repositories (using in-memory for demo)
		userRepo := user.NewInMemoryRepository()
		productRepo := product.NewInMemoryRepository()
		orderRepo := order.NewInMemoryRepository()

		// Create GraphQL server
		gqlConfig := &graphqlserver.Config{
			PlaygroundEnabled:    config.GraphQLPlaygroundEnabled,
			IntrospectionEnabled: true,
		}
		gqlServer := graphqlserver.NewServer(gqlConfig, userRepo, productRepo, orderRepo)

		// Add GraphQL endpoints
		mainMux.Handle("/graphql", gqlServer.Handler())
		if config.GraphQLPlaygroundEnabled {
			mainMux.Handle("/playground", gqlServer.PlaygroundHandler())
		}

		log.Printf("GraphQL endpoint available at: %s/graphql", config.GetHTTPAddress())
		if config.GraphQLPlaygroundEnabled {
			log.Printf("GraphQL Playground available at: %s/playground", config.GetHTTPAddress())
		}
	}

	// Create HTTP server with middleware
	server := &http.Server{
		Addr:         config.GetHTTPAddress(),
		Handler:      corsMiddleware(loggingMiddleware(mainMux)),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("API Gateway started on %s", config.GetHTTPAddress())
	log.Printf("Available endpoints:")
	log.Printf("  REST API:")
	log.Printf("    Users: %s/api/v1/users", config.GetHTTPAddress())
	log.Printf("    Products: %s/api/v1/products", config.GetHTTPAddress())
	log.Printf("    Orders: %s/api/v1/orders", config.GetHTTPAddress())
	log.Printf("  Health: %s/health", config.GetHTTPAddress())
	if config.GraphQLEnabled {
		log.Printf("  GraphQL:")
		log.Printf("    Endpoint: %s/graphql", config.GetHTTPAddress())
		if config.GraphQLPlaygroundEnabled {
			log.Printf("    Playground: %s/playground", config.GetHTTPAddress())
		}
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

// corsMiddleware adds CORS headers
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap ResponseWriter to capture status code
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		log.Printf("HTTP %s %s %d %v %s", r.Method, r.URL.Path, ww.statusCode, duration, r.RemoteAddr)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// healthCheckHandler handles health check requests
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
}
