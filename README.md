# gRPC Microservices Architecture

Đây là một implementation hoàn chỉnh của kiến trúc microservices sử dụng gRPC với Go, tuân thủ các best practices trong industry.

## 🏗️ Kiến trúc hệ thống

```
┌─────────────────┐
│   API Gateway   │ ← REST API (Port 8080)
│  (gRPC-Gateway) │
└─────────┬───────┘
          │
          ├─────────────────────────────────────┐
          │                                     │
          ▼                                     ▼
┌─────────────────┐                   ┌─────────────────┐
│  User Service   │                   │ Product Service │
│   (Port 50051)  │                   │   (Port 50052)  │
└─────────────────┘                   └─────────────────┘
          ▲                                     ▲
          │                                     │
          └─────────────────┬───────────────────┘
                            │
                            ▼
                  ┌─────────────────┐
                  │  Order Service  │
                  │   (Port 50053)  │
                  └─────────────────┘
```

## 🎯 Tính năng chính

### Services
- **User Service**: Quản lý thông tin người dùng
- **Product Service**: Quản lý sản phẩm và inventory
- **Order Service**: Xử lý đơn hàng với inter-service communication
- **API Gateway**: REST API gateway sử dụng gRPC-Gateway

### Best Practices Implemented
- ✅ Clean Architecture (Domain, Service, Repository layers)
- ✅ Protocol Buffers cho API definitions
- ✅ gRPC inter-service communication
- ✅ RESTful API exposure qua gRPC-Gateway
- ✅ Structured logging với Zap
- ✅ Configuration management
- ✅ Health checks
- ✅ Graceful shutdown
- ✅ Error handling và validation
- ✅ Docker containerization
- ✅ Service discovery via configuration

## 🚀 Quick Start

### Prerequisites
- Go 1.23+
- Docker & Docker Compose
- Make
- Protocol Buffers compiler (protoc)

### 1. Clone và Setup
```bash
git clone <repository>
cd learning
```

### 2. Install Dependencies
```bash
make deps
```

### 3. Generate Protocol Buffer Code
```bash
make proto
```

### 4. Run với Docker Compose (Recommended)
```bash
make docker-up
```

### 5. Run Local Development
```bash
# Terminal 1: User Service
make run-user

# Terminal 2: Product Service  
make run-product

# Terminal 3: Order Service
make run-order

# Terminal 4: API Gateway
make run-gateway
```

## 📋 API Documentation

### Base URL
- REST API: `http://localhost:8080`
- Health Check: `http://localhost:8080/health`

### User Service

#### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com", 
    "phone": "+1234567890"
  }'
```

#### Get User
```bash
curl http://localhost:8080/api/v1/users/{user_id}
```

#### List Users
```bash
curl "http://localhost:8080/api/v1/users?page=1&page_size=10"
```

### Product Service

#### Create Product
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "iPhone 15",
    "description": "Latest iPhone model",
    "price": 999.99,
    "stock": 100,
    "category": "Electronics"
  }'
```

#### Get Product
```bash
curl http://localhost:8080/api/v1/products/{product_id}
```

#### List Products
```bash
curl "http://localhost:8080/api/v1/products?page=1&page_size=10&category=Electronics"
```

### Order Service

#### Create Order
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user-uuid-here",
    "items": [
      {
        "product_id": "product-uuid-here",
        "quantity": 2
      }
    ]
  }'
```

#### Get Order
```bash
curl http://localhost:8080/api/v1/orders/{order_id}
```

#### List User Orders
```bash
curl "http://localhost:8080/api/v1/users/{user_id}/orders?page=1&page_size=10"
```

## 🛠️ Development

### Project Structure
```
learning/
├── api/
│   └── proto/              # Protocol Buffer definitions
├── cmd/                    # Main applications
│   ├── user-service/
│   ├── product-service/
│   ├── order-service/
│   └── api-gateway/
├── internal/               # Private packages
│   ├── user/              # User domain
│   ├── product/           # Product domain  
│   ├── order/             # Order domain
│   └── common/            # Shared utilities
├── pkg/                   # Generated protobuf code
├── deployments/           # Docker files
└── scripts/               # Build scripts
```

### Building Services
```bash
# Build all services
make build

# Build individual services
go build ./cmd/user-service
go build ./cmd/product-service  
go build ./cmd/order-service
go build ./cmd/api-gateway
```

### Testing
```bash
# Run tests
go test ./...

# Test with coverage
go test -cover ./...
```

## 🐳 Docker Commands

```bash
# Build all containers
make docker-build

# Start all services
make docker-up

# Stop all services  
make docker-down

# View logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f user-service
```

## ⚙️ Configuration

Services được configure qua environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `HOST` | localhost | Server host |
| `USER_SERVICE_PORT` | 50051 | User service port |
| `PRODUCT_SERVICE_PORT` | 50052 | Product service port |
| `ORDER_SERVICE_PORT` | 50053 | Order service port |
| `GATEWAY_PORT` | 8080 | Gateway port |
| `USER_SERVICE_ADDRESS` | localhost:50051 | User service address |
| `PRODUCT_SERVICE_ADDRESS` | localhost:50052 | Product service address |
| `ORDER_SERVICE_ADDRESS` | localhost:50053 | Order service address |
| `LOG_LEVEL` | info | Logging level |

## 🔍 Monitoring & Health Checks

### Health Check Endpoints
- API Gateway: `http://localhost:8080/health`
- gRPC Health Check: Available on all gRPC services

### Logging
- Structured JSON logging với Zap
- Request/response logging
- Error tracking
- Performance metrics

## 📊 Architecture Decisions

### Technology Choices
- **gRPC**: High-performance, type-safe inter-service communication
- **Protocol Buffers**: Efficient serialization và schema evolution
- **gRPC-Gateway**: REST API compatibility
- **Zap**: High-performance structured logging
- **Clean Architecture**: Separation of concerns và testability

### Design Patterns
- **Repository Pattern**: Data access abstraction
- **Service Layer**: Business logic encapsulation
- **Dependency Injection**: Loose coupling
- **Circuit Breaker**: Resilience (planned)
- **Saga Pattern**: Distributed transactions (planned)

## 🔮 Roadmap

- [ ] Database integration (PostgreSQL/MongoDB)
- [ ] Message queues (RabbitMQ/Apache Kafka)
- [ ] Service discovery (Consul/etcd)
- [ ] Circuit breaker pattern
- [ ] Distributed tracing (Jaeger)
- [ ] Metrics collection (Prometheus)
- [ ] API rate limiting
- [ ] Authentication & authorization
- [ ] Kubernetes deployment

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 