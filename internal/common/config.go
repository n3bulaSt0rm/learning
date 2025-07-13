package common

import (
	"fmt"
	"log"
	"os"
)

// Config holds all configuration for the services
type Config struct {
	// Server config
	Port string
	Host string

	// Service discovery
	UserServiceAddress    string
	ProductServiceAddress string
	OrderServiceAddress   string

	// Database config (for this example, we'll use in-memory storage)
	DatabaseURL string

	// Logging
	LogLevel string
}

// LoadConfig loads configuration from environment variables with defaults
func LoadConfig() *Config {
	config := &Config{
		Port:                  getEnv("PORT", "8080"),
		Host:                  getEnv("HOST", "localhost"),
		UserServiceAddress:    getEnv("USER_SERVICE_ADDRESS", "localhost:50051"),
		ProductServiceAddress: getEnv("PRODUCT_SERVICE_ADDRESS", "localhost:50052"),
		OrderServiceAddress:   getEnv("ORDER_SERVICE_ADDRESS", "localhost:50053"),
		DatabaseURL:           getEnv("DATABASE_URL", ""),
		LogLevel:              getEnv("LOG_LEVEL", "info"),
	}

	log.Printf("Configuration loaded: %+v", config)
	return config
}

// LoadUserServiceConfig loads config specifically for user service
func LoadUserServiceConfig() *Config {
	config := LoadConfig()
	config.Port = getEnv("USER_SERVICE_PORT", "50051")
	return config
}

// LoadProductServiceConfig loads config specifically for product service
func LoadProductServiceConfig() *Config {
	config := LoadConfig()
	config.Port = getEnv("PRODUCT_SERVICE_PORT", "50052")
	return config
}

// LoadOrderServiceConfig loads config specifically for order service
func LoadOrderServiceConfig() *Config {
	config := LoadConfig()
	config.Port = getEnv("ORDER_SERVICE_PORT", "50053")
	return config
}

// LoadGatewayConfig loads config specifically for API gateway
func LoadGatewayConfig() *Config {
	config := LoadConfig()
	config.Port = getEnv("GATEWAY_PORT", "8080")
	return config
}

// GetGRPCAddress returns the full gRPC address
func (c *Config) GetGRPCAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// GetHTTPAddress returns the full HTTP address
func (c *Config) GetHTTPAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}

// getEnv gets environment variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
