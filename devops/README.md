# 🚀 DevOps Learning Hub

A complete DevOps learning environment designed for beginners to advance step-by-step from basic containerization to production-ready deployments.

## 🎯 **Start Here: Learning Roadmap**

**📖 [Complete Learning Roadmap](LEARNING_ROADMAP.md)** - Your guide from beginner to pro

### **Current Phase**: Level 1 - Docker Basics 🐳

Ready to begin your DevOps journey? Start with Docker:

```bash
cd devops/docker
./getting-started.sh
```

## 📁 **Directory Structure**

```
devops/
├── 📋 LEARNING_ROADMAP.md        # Your learning path (START HERE!)
├── 🐳 docker/                   # Level 1: Container fundamentals
│   ├── docs/00-docker-basics.md # Beginner Docker guide
│   ├── getting-started.sh       # Quick start script
│   └── docs/troubleshooting.md  # Common issues & fixes
├── ☸️ kubernetes/               # Level 3: Orchestration basics
│   └── docs/KUBERNETES_THEORY.md
├── 📊 monitoring/               # Level 4: Observability
│   └── docs/MONITORING_THEORY.md
├── 🔄 ci-cd/                   # Level 2 & 8: Automation
│   └── docs/CI_CD_THEORY.md
├── ☁️ aws/                     # Level 5-7: Cloud & EKS
│   ├── docs/EKS_THEORY.md
│   └── terraform/              # Infrastructure as Code
├── 🛠️ scripts/                # Automation scripts
├── 🔐 security/               # Security best practices
└── ⚙️ configs/                # Environment configurations
```

## 🎓 **Learning Phases**

| Phase | Duration | Focus | Outcome |
|-------|----------|-------|---------|
| **📚 Phase 1** | 1-2 weeks | **Foundations** | Docker + Git workflow |
| **🏗️ Phase 2** | 1-2 weeks | **Local Orchestration** | Kubernetes + Basic monitoring |
| **☁️ Phase 3** | 2-3 weeks | **Cloud Basics** | AWS + Terraform |
| **🚀 Phase 4** | 3-4 weeks | **Production Ready** | EKS + CI/CD + Monitoring |
| **🎯 Phase 5** | 1-2 months | **Advanced** | Security + Optimization |

**Total time**: 2-3 months (part-time study)

## 🏁 **Quick Start Options**

### **🆕 Complete Beginner**
```bash
# Start with Docker basics
cd devops/docker
./getting-started.sh
```

### **🐳 Have Docker Experience**
```bash
# Jump to Kubernetes
cd devops/kubernetes
# [Setup guides coming in Level 3]
```

### **☸️ Know Kubernetes**
```bash
# Go to cloud deployment
cd devops/aws
# [EKS guides available now]
```

## 📚 **Documentation by Technology**

### **Level 1: Docker** 🐳
- [Docker Basics Guide](docker/docs/00-docker-basics.md) - Container fundamentals
- [Troubleshooting](docker/docs/troubleshooting.md) - Common Docker issues

### **Level 3: Kubernetes** ☸️
- [Kubernetes Theory](kubernetes/docs/KUBERNETES_THEORY.md) - Core concepts

### **Level 4: Monitoring** 📊
- [Monitoring Theory](monitoring/docs/MONITORING_THEORY.md) - Observability fundamentals

### **Level 8: CI/CD** 🔄
- [CI/CD Theory](ci-cd/docs/CI_CD_THEORY.md) - Pipeline automation

### **Level 7: AWS & EKS** ☁️
- [EKS Theory](aws/docs/EKS_THEORY.md) - Managed Kubernetes on AWS

## 🎯 **Success Milestones**

### **✅ Phase 1 Complete**
- [ ] All microservices running with Docker Compose
- [ ] Git workflow established
- [ ] Basic automation with GitHub Actions

### **✅ Phase 2 Complete**
- [ ] Local Kubernetes cluster running
- [ ] Services deployed with kubectl
- [ ] Basic monitoring dashboard

### **✅ Phase 3 Complete**
- [ ] AWS infrastructure with Terraform
- [ ] EC2 deployment working
- [ ] Cloud networking configured

### **✅ Phase 4 Complete**
- [ ] Production EKS cluster
- [ ] Full CI/CD pipeline
- [ ] Production monitoring

## 🛠️ **Tools You'll Learn**

| Level | Tools | Purpose |
|-------|-------|---------|
| **1-2** | Docker, Git, GitHub Actions | Local development |
| **3-4** | Kubernetes, Prometheus, Grafana | Orchestration & monitoring |
| **5-6** | AWS, Terraform, VPC, EC2 | Cloud infrastructure |
| **7-9** | EKS, ECR, ALB, CI/CD | Production deployment |
| **10+** | Security tools, Service mesh | Enterprise features |

## 🤝 **Getting Help**

### **📋 Check Prerequisites**
Each level has clear prerequisites and setup instructions.

### **🔧 Troubleshooting**
Each technology has its own troubleshooting guide with common issues.

### **📚 Theory First**
Every practical level starts with theory to build understanding.

### **🧪 Hands-on Practice**
All concepts are immediately applied to your microservices project.

## 🌟 **Why This Approach Works**

- **📈 Progressive**: Each level builds on the previous
- **🎯 Practical**: Use real microservices project throughout
- **🔧 Troubleshooting-focused**: Common issues documented
- **⏱️ Time-boxed**: Clear expectations for each level
- **🎓 Beginner-friendly**: No assumed knowledge

## 🚀 **Ready to Start?**

1. **📖 Read the roadmap**: [LEARNING_ROADMAP.md](LEARNING_ROADMAP.md)
2. **🐳 Begin with Docker**: [docker/README.md](docker/README.md)
3. **⚡ Quick start**: `cd devops/docker && ./getting-started.sh`

**Remember**: Take your time, practice each level thoroughly, and don't rush to the next level until you're comfortable with the current one.

**Happy learning! 🎉** 