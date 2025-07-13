# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the applications
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/user-service ./cmd/user-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/product-service ./cmd/product-service  
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/order-service ./cmd/order-service
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api-gateway ./cmd/api-gateway

# User Service
FROM alpine:latest AS user-service
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/user-service .
EXPOSE 50051
CMD ["./user-service"]

# Product Service  
FROM alpine:latest AS product-service
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/product-service .
EXPOSE 50052
CMD ["./product-service"]

# Order Service
FROM alpine:latest AS order-service
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/order-service .
EXPOSE 50053
CMD ["./order-service"]

# API Gateway
FROM alpine:latest AS api-gateway
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/api-gateway .
EXPOSE 8080
CMD ["./api-gateway"] 