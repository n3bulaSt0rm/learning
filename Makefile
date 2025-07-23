.PHONY: proto build clean run-user run-product run-order run-gateway docker-build docker-up docker-down

# Variables
PROTO_DIR = api/proto
PKG_DIR = pkg
BIN_DIR = bin

# Install dependencies
deps:
	go mod tidy
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Generate protobuf files
proto: clean-proto
	mkdir -p $(PKG_DIR)/user/pb $(PKG_DIR)/product/pb $(PKG_DIR)/order/pb
	
	# Generate User service
	protoc --proto_path=$(PROTO_DIR) --proto_path=. \
		--go_out=$(PKG_DIR)/user/pb --go_opt=paths=source_relative \
		--go-grpc_out=$(PKG_DIR)/user/pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(PKG_DIR)/user/pb --grpc-gateway_opt=paths=source_relative \
		$(PROTO_DIR)/user.proto
	
	# Generate Product service
	protoc --proto_path=$(PROTO_DIR) --proto_path=. \
		--go_out=$(PKG_DIR)/product/pb --go_opt=paths=source_relative \
		--go-grpc_out=$(PKG_DIR)/product/pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(PKG_DIR)/product/pb --grpc-gateway_opt=paths=source_relative \
		$(PROTO_DIR)/product.proto
	
	# Generate Order service
	protoc --proto_path=$(PROTO_DIR) --proto_path=. \
		--go_out=$(PKG_DIR)/order/pb --go_opt=paths=source_relative \
		--go-grpc_out=$(PKG_DIR)/order/pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(PKG_DIR)/order/pb --grpc-gateway_opt=paths=source_relative \
		$(PROTO_DIR)/order.proto

# Build all services
build: proto
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/user-service ./cmd/user-service
	go build -o $(BIN_DIR)/product-service ./cmd/product-service
	go build -o $(BIN_DIR)/order-service ./cmd/order-service
	go build -o $(BIN_DIR)/api-gateway ./cmd/api-gateway

# Clean generated files
clean:
	rm -rf $(BIN_DIR)
	rm -rf $(PKG_DIR)

clean-proto:
	rm -rf $(PKG_DIR)

# Run services individually
run-user:
	go run ./cmd/user-service

run-product:
	go run ./cmd/product-service

run-order:
	go run ./cmd/order-service

run-gateway:
	go run ./cmd/api-gateway

# Docker commands
docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

# Development
dev: proto
	@echo "Starting all services..."
	@make run-user &
	@make run-product &
	@make run-order &
	@make run-gateway &
	@wait 

# GraphQL commands
.PHONY: graphql-generate
graphql-generate:
	@echo "Generating GraphQL code..."
	go run github.com/99designs/gqlgen generate

.PHONY: graphql-init
graphql-init:
	@echo "Initializing GraphQL..."
	mkdir -p internal/graphql/{schema,resolvers,models,generated}
	$(MAKE) graphql-generate

.PHONY: dev-graphql
dev-graphql:
	@echo "Starting GraphQL playground..."
	@echo "Visit http://localhost:8080/graphql for GraphQL playground"
	$(MAKE) run-gateway

# DataLoader generation
.PHONY: generate-dataloaders
generate-dataloaders:
	@echo "Generating DataLoaders..."
	go generate ./internal/graphql/dataloaders/... 