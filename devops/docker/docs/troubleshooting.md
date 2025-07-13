# 🔧 Docker Troubleshooting Guide

## 🚨 **Common Issues & Solutions**

### **1. Docker Installation Issues**

#### **Problem**: `docker: command not found`
**Solution**:
```bash
# Check if Docker is installed
which docker

# If not installed:
# - Windows/Mac: Install Docker Desktop
# - Linux: Follow installation guide in 00-docker-basics.md
```

#### **Problem**: `Cannot connect to the Docker daemon`
**Solutions**:
```bash
# Check if Docker daemon is running
docker info

# On Linux - start Docker service
sudo systemctl start docker
sudo systemctl enable docker

# On Windows/Mac - start Docker Desktop

# Add user to docker group (Linux)
sudo usermod -aG docker $USER
newgrp docker
```

---

### **2. Permission Issues**

#### **Problem**: `Permission denied while trying to connect to Docker daemon socket`
**Solution**:
```bash
# On Linux - add user to docker group
sudo usermod -aG docker $USER
newgrp docker

# Or run with sudo (not recommended for regular use)
sudo docker ps
```

#### **Problem**: `dial unix /var/run/docker.sock: connect: permission denied`
**Solution**:
```bash
# Fix socket permissions (Linux)
sudo chmod 666 /var/run/docker.sock

# Better solution - add user to docker group
sudo usermod -aG docker $USER
```

---

### **3. Build Issues**

#### **Problem**: `COPY failed: no such file or directory`
**Solution**:
```dockerfile
# Make sure paths are correct relative to build context
# Incorrect:
COPY /absolute/path/file .

# Correct:
COPY ./relative/path/file .
COPY . .
```

#### **Problem**: `go.mod not found` in Dockerfile
**Solution**:
```dockerfile
# For Go projects - copy go.mod from correct location
# If Dockerfile is in cmd/service-name/:
COPY ../../go.mod ../../go.sum ./

# Or build from project root:
# docker build -f cmd/service-name/Dockerfile .
```

#### **Problem**: Build is very slow
**Solutions**:
```dockerfile
# Use .dockerignore file
# Create .dockerignore in project root:
echo "node_modules
.git
*.log
tmp/
.env" > .dockerignore

# Use multi-stage builds
FROM golang:1.19-alpine AS builder
# ... build steps ...
FROM alpine:latest
COPY --from=builder /app/binary .
```

---

### **4. Runtime Issues**

#### **Problem**: Container exits immediately
**Solutions**:
```bash
# Check logs
docker logs container-name

# Common causes:
# 1. Main process exits
# 2. Missing executable permissions
# 3. Wrong CMD/ENTRYPOINT

# Debug by running interactively
docker run -it image-name sh
```

#### **Problem**: `Port already in use`
**Solutions**:
```bash
# Check what's using the port
lsof -i :8080        # Mac/Linux
netstat -ano | findstr :8080  # Windows

# Use different port
docker run -p 8081:8080 service-name

# Stop conflicting container
docker ps
docker stop container-name
```

#### **Problem**: `Cannot reach service on localhost`
**Solutions**:
```bash
# Make sure port is exposed and mapped
docker run -p 8080:8080 service-name  # Map host:container

# Check if service is listening on 0.0.0.0, not 127.0.0.1
# In your app, bind to 0.0.0.0:8080, not localhost:8080

# Check container is running
docker ps
```

---

### **5. Docker Compose Issues**

#### **Problem**: `docker-compose: command not found`
**Solutions**:
```bash
# Try new syntax
docker compose up

# Install docker-compose
# Linux:
sudo apt install docker-compose
# Or with pip:
pip install docker-compose

# Mac: Already included in Docker Desktop
```

#### **Problem**: Services can't communicate
**Solutions**:
```yaml
# Use service names as hostnames
services:
  user-service:
    # ...
  api-gateway:
    environment:
      - USER_SERVICE_URL=http://user-service:8080  # Use service name
```

#### **Problem**: `depends_on` not working
**Solutions**:
```yaml
# Use healthchecks for true dependency
services:
  postgres:
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 30s
      timeout: 10s
      retries: 3
  
  app:
    depends_on:
      postgres:
        condition: service_healthy  # Wait for healthy status
```

#### **Problem**: Changes not reflected after rebuild
**Solutions**:
```bash
# Force rebuild without cache
docker-compose build --no-cache

# Or rebuild specific service
docker-compose build --no-cache user-service

# Remove old images
docker image prune -f
```

---

### **6. Resource Issues**

#### **Problem**: Container killed (OOMKilled)
**Solutions**:
```bash
# Check memory usage
docker stats

# Increase memory limits
docker run -m 512m service-name

# In docker-compose.yml:
services:
  app:
    deploy:
      resources:
        limits:
          memory: 512M
```

#### **Problem**: Disk space issues
**Solutions**:
```bash
# Clean up Docker resources
docker system prune -a

# Remove unused images
docker image prune -a

# Remove unused volumes
docker volume prune

# Remove unused networks
docker network prune

# Check disk usage
docker system df
```

---

### **7. Networking Issues**

#### **Problem**: Services can't connect to external APIs
**Solutions**:
```bash
# Check if container has internet access
docker run alpine ping google.com

# Check DNS resolution
docker run alpine nslookup google.com

# Use host networking (for debugging only)
docker run --network host service-name
```

#### **Problem**: Custom network issues
**Solutions**:
```bash
# Create custom network
docker network create mynetwork

# List networks
docker network ls

# Inspect network
docker network inspect mynetwork

# Connect container to network
docker run --network mynetwork service-name
```

---

### **8. Development Workflow Issues**

#### **Problem**: Code changes require full rebuild
**Solutions**:
```yaml
# Use volume mounts for development
services:
  app:
    volumes:
      - ./src:/app/src:ro  # Mount source code
      - node_modules:/app/node_modules  # Anonymous volume for deps
```

#### **Problem**: Database data lost on restart
**Solutions**:
```yaml
# Use named volumes for persistence
services:
  postgres:
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

---

## 🔍 **Debugging Commands**

### **Container Debugging**
```bash
# View container details
docker inspect container-name

# Execute command in running container
docker exec -it container-name sh

# View real-time logs
docker logs -f container-name

# View container processes
docker top container-name

# View resource usage
docker stats container-name
```

### **Image Debugging**
```bash
# View image layers
docker history image-name

# Run container with shell (debug mode)
docker run -it --entrypoint sh image-name

# Check image details
docker inspect image-name
```

### **Network Debugging**
```bash
# List networks
docker network ls

# Inspect network
docker network inspect bridge

# Test connectivity between containers
docker exec container1 ping container2
```

---

## 📋 **Quick Fixes Checklist**

**Before asking for help, try these:**

- [ ] Check container logs: `docker logs container-name`
- [ ] Verify ports are mapped: `docker ps`
- [ ] Check if service is actually running in container: `docker exec -it container sh`
- [ ] Verify environment variables: `docker exec container env`
- [ ] Check disk space: `docker system df`
- [ ] Try rebuilding without cache: `docker build --no-cache`
- [ ] Clean up resources: `docker system prune`
- [ ] Restart Docker daemon/Desktop

---

## 🆘 **Still Need Help?**

1. **Check Docker official docs**: https://docs.docker.com/
2. **Search Docker forums**: https://forums.docker.com/
3. **Stack Overflow**: Tag your question with `docker`
4. **Include in your question**:
   - Docker version: `docker --version`
   - OS and version
   - Complete error message
   - Dockerfile content
   - docker-compose.yml content
   - Steps to reproduce

**Pro tip**: Most Docker issues are due to:
1. Incorrect file paths in Dockerfile
2. Port mapping problems
3. Permission issues
4. Resource constraints

Happy debugging! 🐳 