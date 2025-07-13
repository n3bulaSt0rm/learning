# 🚀 Kubernetes Learning Hub

## 📁 **Folder Structure**

```
kubernetes/
├── docs/           # Theory and learning materials
├── manifests/      # Kubernetes YAML manifests
└── helm/          # Helm charts and templates
```

## 📚 **Learning Path**

### **Phase 2, Level 3: Kubernetes Basics**
- **Start here**: `docs/KUBERNETES_THEORY.md`
- **Time**: 5-7 days
- **Goal**: Deploy microservices on local K8s cluster

### **Learning Sequence:**
1. **Theory** (`docs/`): Understand K8s concepts
2. **Practice** (`manifests/`): Deploy using raw YAML
3. **Advanced** (`helm/`): Use Helm for templating

---

## 🎯 **What's Inside**

### **📖 docs/**
- `KUBERNETES_THEORY.md` - Complete K8s architecture guide
- Pods, Services, Deployments concepts
- Networking and Storage basics

### **⚙️ manifests/**
- `namespace.yaml` - Namespace setup
- `configmap.yaml` - Configuration management
- `*-service.yaml` - Microservice deployments
- `ingress.yaml` - Traffic routing
- `monitoring/` - Prometheus, Grafana configs

### **📦 helm/**
- `microservices/` - Helm chart for all services
- Templates for easy deployment
- Values files for different environments

---

## 🛠️ **Quick Start**

### **1. Local Setup**
```bash
# Install kind or minikube
kind create cluster --name learning

# Deploy namespace
kubectl apply -f manifests/namespace.yaml
```

### **2. Deploy Services**
```bash
# Deploy all services
kubectl apply -f manifests/

# Check deployment
kubectl get pods -n microservices
```

### **3. Using Helm**
```bash
# Install with Helm
helm install microservices helm/microservices/ -n microservices

# Upgrade
helm upgrade microservices helm/microservices/ -n microservices
```

---

## 🔄 **From Docker to K8s**

If you completed **Phase 1 (Docker)**, you already have:
- ✅ Containerized applications
- ✅ Docker Compose experience
- ✅ Container networking basics

**Now you'll learn:**
- 🎯 Container orchestration
- 🎯 Service discovery
- 🎯 Auto-scaling and healing
- 🎯 Production-ready deployments

---

## 📈 **Success Criteria**

By the end of Phase 2, Level 3, you should be able to:
- [ ] Deploy microservices to local K8s cluster
- [ ] Configure service-to-service communication
- [ ] Setup ingress for external access
- [ ] Monitor services with Prometheus/Grafana
- [ ] Perform rolling updates and rollbacks

---

## 🔗 **Next Steps**

After mastering local K8s:
- **Phase 3**: AWS & Terraform foundations
- **Phase 4**: Production EKS deployment
- **Phase 5**: Advanced topics (Security, Service Mesh)

---

## 🆘 **Need Help?**

- Check `docs/KUBERNETES_THEORY.md` for concepts
- Use `kubectl describe` to debug issues
- Join K8s community forums
- Practice with simple deployments first 