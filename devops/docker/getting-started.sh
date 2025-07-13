#!/bin/bash

# Docker Level 1 - Getting Started Script
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🐳 Welcome to Level 1: Docker Basics!${NC}"
echo "=================================================================="
echo

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to print step
print_step() {
    echo -e "${GREEN}➤ $1${NC}"
}

# Function to print warning
print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Function to print error
print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check prerequisites
print_step "Checking prerequisites..."

if command_exists docker; then
    echo "✅ Docker is installed: $(docker --version)"
else
    print_error "Docker is not installed!"
    echo "Please install Docker first:"
    echo "- Windows/Mac: https://www.docker.com/products/docker-desktop"
    echo "- Linux: Follow instructions in docs/00-docker-basics.md"
    exit 1
fi

if command_exists docker-compose; then
    echo "✅ Docker Compose is installed: $(docker-compose --version)"
else
    print_warning "Docker Compose not found. Checking for 'docker compose' command..."
    if docker compose version >/dev/null 2>&1; then
        echo "✅ Docker Compose is available as 'docker compose'"
        COMPOSE_CMD="docker compose"
    else
        print_error "Docker Compose is not available!"
        echo "Please install Docker Compose or use Docker Desktop"
        exit 1
    fi
fi

# Set compose command
COMPOSE_CMD=${COMPOSE_CMD:-"docker-compose"}

# Check if Docker daemon is running
print_step "Checking Docker daemon..."
if docker info >/dev/null 2>&1; then
    echo "✅ Docker daemon is running"
else
    print_error "Docker daemon is not running!"
    echo "Please start Docker Desktop or Docker daemon"
    exit 1
fi

echo
print_step "Let's start with a simple test..."

# Run hello-world
echo "Running 'docker run hello-world'..."
if docker run hello-world >/dev/null 2>&1; then
    echo "✅ Docker hello-world test passed!"
else
    print_error "Docker hello-world test failed!"
    exit 1
fi

echo
print_step "Checking project structure..."

# Check if we're in the right directory
if [[ ! -f "go.mod" ]]; then
    print_error "Please run this script from the project root directory (where go.mod is located)"
    exit 1
fi

echo "✅ Found go.mod - looks like you're in the right directory"

# Check for Dockerfiles
echo "Checking for Dockerfiles..."
dockerfile_count=0

for service in cmd/user-service cmd/product-service cmd/order-service cmd/api-gateway; do
    if [[ -f "$service/Dockerfile" ]]; then
        echo "✅ Found $service/Dockerfile"
        ((dockerfile_count++))
    else
        print_warning "Missing $service/Dockerfile"
        echo "   You'll need to create this in the next steps"
    fi
done

if [[ $dockerfile_count -eq 0 ]]; then
    print_warning "No Dockerfiles found yet - that's okay! We'll create them."
fi

# Check for docker-compose.yml
if [[ -f "docker-compose.yml" ]]; then
    echo "✅ Found docker-compose.yml"
else
    print_warning "Missing docker-compose.yml - you'll create this following the guide"
fi

echo
print_step "What's next?"
echo "1. 📖 Read the theory guide: devops/docker/docs/00-docker-basics.md"
echo "2. 🔨 Create Dockerfiles for your microservices"
echo "3. 📝 Create docker-compose.yml"
echo "4. 🚀 Run 'docker-compose up -d'"
echo "5. 🧪 Test your services"

echo
echo -e "${BLUE}📚 Quick Commands Reference:${NC}"
echo "• View this guide: cat devops/docker/docs/00-docker-basics.md"
echo "• Build image: docker build -t service-name ."
echo "• Run container: docker run -d -p 8080:8080 service-name"
echo "• View logs: docker logs container-name"
echo "• Start all services: $COMPOSE_CMD up -d"
echo "• Stop all services: $COMPOSE_CMD down"
echo "• View running containers: docker ps"

echo
echo -e "${GREEN}Ready to start Level 1! 🚀${NC}"
echo "Follow the guide in devops/docker/docs/00-docker-basics.md"

# Optional: Create sample Dockerfile template
read -p "Would you like to create a sample Dockerfile template for user-service? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    mkdir -p cmd/user-service
    if [[ ! -f "cmd/user-service/Dockerfile" ]]; then
        cat > cmd/user-service/Dockerfile << 'EOF'
# Multi-stage build for Go application
FROM golang:1.19-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY ../../go.mod ../../go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .
COPY ../../internal ./internal

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service .

# Final stage - smaller image
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/user-service .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./user-service"]
EOF
        echo "✅ Created sample Dockerfile at cmd/user-service/Dockerfile"
        echo "You may need to adjust the paths based on your actual project structure"
    else
        echo "Dockerfile already exists, skipping..."
    fi
fi

echo
echo "=================================================================="
echo -e "${BLUE}Happy Dockerizing! 🐳${NC}" 