# 🐳 Level 1: Docker Basics

Welcome to your first step in the DevOps journey! This level will teach you Docker fundamentals through hands-on practice with your microservices project.

## 📋 **What You'll Learn**

- **Container Fundamentals**: What containers are and why they matter
- **Docker Basics**: Images, containers, and Docker lifecycle
- **Dockerfile Creation**: Building custom images for your services
- **Docker Compose**: Orchestrating multi-service applications
- **Development Workflow**: Using Docker for local development

## ⏱️ **Time Estimate**

**3-5 days** (2-3 hours per day)

## 🎯 **Learning Objectives**

By the end of this level, you will:
- ✅ Understand the difference between containers and VMs
- ✅ Be able to write Dockerfiles for Go microservices
- ✅ Successfully containerize all 4 microservices
- ✅ Use Docker Compose to run the entire application stack
- ✅ Debug common Docker issues

## 🚀 **Quick Start**

1. **Run the getting started script:**
   ```bash
   cd devops/docker
   chmod +x getting-started.sh
   ./getting-started.sh
   ```

2. **Follow the comprehensive guide:**
   ```bash
   # View the main guide
   cat docs/00-docker-basics.md
   
   # Or open in your favorite editor
   code docs/00-docker-basics.md
   ```

## 📚 **Resources**

| Resource | Description | Time |
|----------|-------------|------|
| [Docker Basics Guide](docs/00-docker-basics.md) | Complete beginner guide | 📖 2-3 hours |
| [Troubleshooting Guide](docs/troubleshooting.md) | Common issues & solutions | 🔧 As needed |
| [Getting Started Script](getting-started.sh) | Automated setup checker | ⚡ 5 minutes |

## 📁 **Directory Structure**

```
docker/
├── README.md                    # This file
├── getting-started.sh           # Quick start script
└── docs/
    ├── 00-docker-basics.md      # Main learning guide
    └── troubleshooting.md        # Common issues & fixes
```

## 🎯 **Success Criteria**

**You'll know you've completed this level when:**

### **✅ Knowledge Check**
- [ ] Can explain what a container is in simple terms
- [ ] Understand the difference between images and containers
- [ ] Know when to use Docker vs other solutions
- [ ] Understand basic Dockerfile instructions

### **✅ Practical Skills**
- [ ] Successfully installed Docker
- [ ] Created Dockerfiles for all microservices
- [ ] Built and ran individual containers
- [ ] Created docker-compose.yml for the full stack
- [ ] All services communicate properly

### **✅ Final Deliverable**
- [ ] **Running Application**: All 4 microservices + database working together
- [ ] **API Gateway**: Can route requests to backend services
- [ ] **Health Checks**: All services respond to health endpoints
- [ ] **Data Persistence**: Database data survives container restarts

## 🧪 **Testing Your Success**

Run these commands to verify everything works:

```bash
# Start the entire stack
docker-compose up -d

# Check all services are running
docker-compose ps

# Test API Gateway
curl http://localhost:8080/health

# Test individual services
curl http://localhost:8081/health  # user-service
curl http://localhost:8082/health  # product-service
curl http://localhost:8083/health  # order-service

# Create test data
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com"}'

# Verify data persists after restart
docker-compose restart
curl http://localhost:8080/api/users
```

## 🚨 **Common Issues**

| Issue | Quick Fix | Details |
|-------|-----------|---------|
| `docker: command not found` | Install Docker Desktop | [Installation Guide](docs/00-docker-basics.md#installation) |
| `Permission denied` | Add user to docker group | [Troubleshooting](docs/troubleshooting.md#permission-issues) |
| `Port already in use` | Change port mapping | [Port Conflicts](docs/troubleshooting.md#runtime-issues) |
| Build fails | Check file paths | [Build Issues](docs/troubleshooting.md#build-issues) |

## 🛠️ **Development Tips**

### **Faster Development Cycle**
```bash
# Rebuild specific service
docker-compose build user-service

# View logs for debugging
docker-compose logs -f user-service

# Clean up when needed
docker system prune
```

### **Best Practices**
- Use `.dockerignore` to exclude unnecessary files
- Leverage multi-stage builds for smaller images
- Use specific image tags, not `latest`
- Run containers as non-root user when possible

## 🎓 **What's Next?**

After completing Level 1, you'll be ready for:

**Level 2**: [Git & GitHub Workflow](../ci-cd/docs/01-git-basics.md)
- Version control best practices
- Collaborative development workflow
- Basic GitHub Actions for testing

**Or dive deeper into Docker:**
- Container security best practices
- Docker networking deep dive
- Production Docker patterns

## 🤝 **Getting Help**

**Stuck on something?**

1. **Check the troubleshooting guide**: [troubleshooting.md](docs/troubleshooting.md)
2. **Review common patterns**: Look at the examples in the main guide
3. **Debug step by step**: Use `docker logs` and `docker exec` commands
4. **Start fresh**: Sometimes `docker-compose down` and starting over helps

**Still need help?**
- Make sure you've followed all prerequisites
- Include error messages and your environment details
- Double-check file paths and permissions

## 🏆 **Congratulations!**

When you complete this level, you'll have:
- A solid understanding of containerization
- Working Docker setup for your microservices
- Foundation for more advanced DevOps practices
- Confidence to tackle the next level!

**Ready to containerize your microservices? Let's go! 🚀**

---

📖 **Start here**: [Docker Basics Guide](docs/00-docker-basics.md)

⚡ **Quick start**: `./getting-started.sh` 