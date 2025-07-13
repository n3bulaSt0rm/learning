# 🐳 Level 1: Docker Basics for Beginners

## 🎯 **Mục tiêu Level này**

Sau khi hoàn thành level này, bạn sẽ:
- ✅ Hiểu container là gì và tại sao cần Docker
- ✅ Biết cách viết Dockerfile đơn giản
- ✅ Có thể containerize microservices của mình
- ✅ Chạy được multi-service với Docker Compose

**Thời gian**: 3-5 ngày (2-3 giờ/ngày)

---

## 📚 **1. Lý thuyết cơ bản**

### **Container là gì?**

Hãy tưởng tượng container như một **"hộp đóng gói hoàn hảo"** cho ứng dụng:

```
🏠 Traditional Deployment          📦 Container Deployment
┌─────────────────────┐           ┌─────────────────────┐
│   Your App          │           │  📦 Container        │
│   ┌─────────────┐   │           │  ┌─────────────┐    │
│   │    Code     │   │           │  │    Code     │    │
│   └─────────────┘   │           │  │    Runtime  │    │
│   Dependencies      │     VS    │  │    Libs     │    │
│   Runtime           │           │  │    OS       │    │
│   Libraries         │           │  └─────────────┘    │
│   OS                │           └─────────────────────┘
└─────────────────────┘           Runs anywhere! 🚀
```

### **Tại sao cần Docker?**

**Vấn đề thường gặp:**
- "Code chạy được trên máy tôi mà!" 😅
- Setup môi trường phức tạp
- Khác biệt giữa dev/staging/production
- Dependency conflicts

**Docker giải quyết:**
- ✅ **Consistency**: Chạy giống nhau ở mọi nơi
- ✅ **Isolation**: Mỗi service riêng biệt
- ✅ **Portability**: Chạy được trên bất kỳ platform nào
- ✅ **Scalability**: Dễ dàng scale up/down

### **Docker vs Virtual Machine**

```
Virtual Machine                    Docker Container
┌─────────────────────┐           ┌─────────────────────┐
│     App A           │           │     App A           │
│  ┌─────────────┐    │           │  ┌─────────────┐    │
│  │   Bins/Libs │    │           │  │   Bins/Libs │    │
│  └─────────────┘    │           │  └─────────────┘    │
│    Guest OS         │           │                     │
└─────────────────────┘           └─────────────────────┘
│     App B           │           │     App B           │
│  ┌─────────────┐    │           │  ┌─────────────┐    │
│  │   Bins/Libs │    │           │  │   Bins/Libs │    │
│  └─────────────┘    │           │  └─────────────┘    │
│    Guest OS         │           │                     │
└─────────────────────┘           └─────────────────────┘
    Hypervisor                      Docker Engine
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
       Host OS                         Host OS
```

**Ưu điểm của Container:**
- Nhẹ hơn (share OS kernel)
- Khởi động nhanh hơn (giây vs phút)
- Sử dụng ít tài nguyên hơn

---

## 🛠️ **2. Cài đặt Docker**

### **Windows/Mac:**
1. Download Docker Desktop: https://www.docker.com/products/docker-desktop
2. Install và restart máy
3. Kiểm tra: `docker --version`

### **Linux (Ubuntu):**
```bash
# Update package index
sudo apt update

# Install prerequisites
sudo apt install apt-transport-https ca-certificates curl software-properties-common

# Add Docker's official GPG key
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# Add Docker repository
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker
sudo apt update
sudo apt install docker-ce docker-ce-cli containerd.io

# Add user to docker group (to run without sudo)
sudo usermod -aG docker $USER
newgrp docker

# Verify installation
docker --version
```

---

## 🚀 **3. Thực hành đầu tiên**

### **3.1. Hello World với Docker**

```bash
# Chạy container đầu tiên
docker run hello-world

# Xem các container đang chạy
docker ps

# Xem tất cả container (bao gồm đã dừng)
docker ps -a

# Xem các image
docker images
```

### **3.2. Chạy web server đơn giản**

```bash
# Chạy nginx
docker run -d -p 8080:80 --name my-nginx nginx

# Kiểm tra: mở browser http://localhost:8080
# Dừng container
docker stop my-nginx

# Xóa container
docker rm my-nginx
```

**Giải thích các flags:**
- `-d`: Chạy ở background (detached mode)
- `-p 8080:80`: Map port 8080 (host) → 80 (container)
- `--name`: Đặt tên cho container

---

## 📝 **4. Dockerfile - Tạo Image riêng**

### **4.1. Dockerfile cơ bản**

Tạo file `Dockerfile` (không có extension):

```dockerfile
# Dockerfile for a simple Node.js app
FROM node:16-alpine

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Expose port
EXPOSE 3000

# Command to run the app
CMD ["npm", "start"]
```

### **4.2. Dockerfile cho User Service**

Tạo `user-service/Dockerfile`:

```dockerfile
# Start from official Go image
FROM golang:1.19-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd/user-service

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
```

### **4.3. Build và run image**

```bash
# Build image
cd user-service
docker build -t user-service:v1.0 .

# Run container
docker run -d -p 8081:8080 --name user-service user-service:v1.0

# Check logs
docker logs user-service

# Test the service
curl http://localhost:8081/health
```

---

## 🎭 **5. Docker Compose - Multi-service**

### **5.1. Tại sao cần Docker Compose?**

Thay vì chạy từng container riêng lẻ:
```bash
docker run -d -p 8081:8080 user-service
docker run -d -p 8082:8080 product-service
docker run -d -p 8083:8080 order-service
docker run -d -p 8080:8080 api-gateway
docker run -d -p 5432:5432 postgres
```

Ta có thể dùng 1 file YAML để quản lý tất cả!

### **5.2. docker-compose.yml cơ bản**

Tạo file `docker-compose.yml` ở root project:

```yaml
version: '3.8'

services:
  # Database
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: microservices
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password123
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 30s
      timeout: 10s
      retries: 3

  # User Service
  user-service:
    build:
      context: ./cmd/user-service
      dockerfile: Dockerfile
    ports:
      - "8081:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=microservices
      - DB_USER=postgres
      - DB_PASSWORD=password123
    depends_on:
      postgres:
        condition: service_healthy

  # Product Service
  product-service:
    build:
      context: ./cmd/product-service
      dockerfile: Dockerfile
    ports:
      - "8082:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=microservices
      - DB_USER=postgres
      - DB_PASSWORD=password123
    depends_on:
      postgres:
        condition: service_healthy

  # Order Service
  order-service:
    build:
      context: ./cmd/order-service
      dockerfile: Dockerfile
    ports:
      - "8083:8080"
    environment:
      - DB_HOST=postgres
      - DB_PORT=5432
      - DB_NAME=microservices
      - DB_USER=postgres
      - DB_PASSWORD=password123
      - USER_SERVICE_URL=http://user-service:8080
      - PRODUCT_SERVICE_URL=http://product-service:8080
    depends_on:
      postgres:
        condition: service_healthy

  # API Gateway
  api-gateway:
    build:
      context: ./cmd/api-gateway
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - USER_SERVICE_URL=http://user-service:8080
      - PRODUCT_SERVICE_URL=http://product-service:8080
      - ORDER_SERVICE_URL=http://order-service:8080
    depends_on:
      - user-service
      - product-service
      - order-service

volumes:
  postgres_data:
```

### **5.3. Commands cơ bản**

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f user-service

# Stop all services
docker-compose down

# Rebuild and restart
docker-compose up -d --build

# View running services
docker-compose ps
```

---

## 🧪 **6. Testing Setup**

### **6.1. Health checks**

Test từng service:
```bash
# User Service
curl http://localhost:8081/health

# Product Service  
curl http://localhost:8082/health

# Order Service
curl http://localhost:8083/health

# API Gateway
curl http://localhost:8080/health
```

### **6.2. Create sample data**

```bash
# Create a user
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

# Create a product
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{"name":"iPhone 15","price":999.99,"description":"Latest iPhone"}'

# Get users
curl http://localhost:8080/api/users
```

---

## 📋 **7. Checklist hoàn thành Level 1**

### **Lý thuyết:**
- [ ] Hiểu container vs VM
- [ ] Biết Docker architecture cơ bản
- [ ] Hiểu Dockerfile instructions
- [ ] Biết khi nào dùng Docker Compose

### **Thực hành:**
- [ ] Cài đặt Docker thành công
- [ ] Chạy được container đầu tiên
- [ ] Viết được Dockerfile cho ít nhất 1 service
- [ ] Setup được docker-compose.yml
- [ ] All services chạy được với `docker-compose up`
- [ ] Test API qua API Gateway thành công

### **Troubleshooting:**
- [ ] Biết cách xem logs: `docker logs <container>`
- [ ] Biết cách debug: `docker exec -it <container> sh`
- [ ] Biết cách clean up: `docker system prune`

---

## 🎉 **Deliverable**

**Mục tiêu cuối Level 1:**
- ✅ Tất cả 4 microservices chạy được với Docker Compose
- ✅ API Gateway có thể gọi được các microservices
- ✅ Database connection hoạt động
- ✅ Health checks pass cho tất cả services

**Files cần có:**
```
project/
├── docker-compose.yml
├── cmd/
│   ├── user-service/Dockerfile
│   ├── product-service/Dockerfile
│   ├── order-service/Dockerfile
│   └── api-gateway/Dockerfile
└── README.md (updated with Docker instructions)
```

---

## 🚀 **Next Steps**

Sau khi hoàn thành Level 1:
1. **Level 2**: Git & GitHub Workflow
2. **Level 3**: Kubernetes Basics
3. **Level 4**: Basic Monitoring

**Questions or issues?** Check [troubleshooting.md](./troubleshooting.md)

**Ready for Level 2?** 👉 [Level 2: Git & GitHub Workflow](../ci-cd/docs/01-git-basics.md) 