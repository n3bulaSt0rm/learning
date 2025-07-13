# Kubernetes Theory Guide - Kiến thức Cơ bản

## 1. Containerization và Docker

### 1.1 Container là gì?
- **Container** là một cách đóng gói ứng dụng cùng với tất cả dependencies (thư viện, runtime, system tools)
- Khác với Virtual Machine, container chia sẻ OS kernel với host
- **Lợi ích**:
  - Lightweight (nhẹ hơn VM)
  - Portable (chạy được ở mọi nơi)
  - Consistent (môi trường nhất quán)
  - Scalable (dễ scale)

### 1.2 Docker Fundamentals
```bash
# Docker Image: Template để tạo container
# Docker Container: Instance đang chạy của image
# Dockerfile: Recipe để build image

# Example Dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o app

FROM alpine:latest
COPY --from=builder /app/app .
CMD ["./app"]
```

### 1.3 Multi-stage Build
- **Stage 1 (Builder)**: Compile code, install dependencies
- **Stage 2 (Runtime)**: Chỉ chứa những gì cần thiết để chạy
- **Lợi ích**: Image size nhỏ, security tốt hơn

## 2. Kubernetes Architecture

### 2.1 Control Plane Components
```
┌─────────────────┐
│   Control Plane │
├─────────────────┤
│ • API Server    │ ← Entry point for all operations
│ • etcd          │ ← Distributed key-value store
│ • Scheduler     │ ← Decides where to place pods
│ • Controller    │ ← Maintains desired state
└─────────────────┘
```

### 2.2 Worker Node Components
```
┌─────────────────┐
│   Worker Node   │
├─────────────────┤
│ • kubelet       │ ← Node agent, manages pods
│ • kube-proxy    │ ← Network proxy, load balancing
│ • Container RT  │ ← Docker/containerd
└─────────────────┘
```

### 2.3 Key Concepts

#### **Pod**
- Smallest deployable unit
- Chứa 1 hoặc nhiều containers
- Shared network và storage
- Ephemeral (có thể bị destroy bất cứ lúc nào)

#### **Deployment**
- Manages ReplicaSets
- Handles rolling updates
- Ensures desired number of pods

#### **Service**
- Stable network endpoint
- Load balances traffic to pods
- Types: ClusterIP, NodePort, LoadBalancer

#### **Ingress**
- HTTP/HTTPS routing to services
- SSL termination
- Path-based routing

## 3. Kubernetes Objects Deep Dive

### 3.1 Pod Lifecycle
```
Pending → Running → Succeeded/Failed
     ↓
   Terminated
```

### 3.2 Deployment Strategy
```yaml
# Rolling Update (Default)
strategy:
  type: RollingUpdate
  rollingUpdate:
    maxUnavailable: 1
    maxSurge: 1

# Blue-Green Deployment
# Recreate Strategy
```

### 3.3 Service Types
```yaml
# ClusterIP - Internal only
apiVersion: v1
kind: Service
spec:
  type: ClusterIP  # Default

# NodePort - External access via node IP
spec:
  type: NodePort
  ports:
  - port: 80
    nodePort: 30000

# LoadBalancer - Cloud provider LB
spec:
  type: LoadBalancer
```

## 4. Networking trong Kubernetes

### 4.1 Cluster Networking
```
┌─────────────────────────────────────┐
│              Cluster                │
├─────────────────────────────────────┤
│  Node 1           Node 2            │
│  ┌─────┐         ┌─────┐           │
│  │Pod A│ ◄──────► │Pod B│           │
│  └─────┘         └─────┘           │
│     │               │              │
│  ┌─────────────────────┐           │
│  │    Service          │           │
│  └─────────────────────┘           │
│           │                        │
│  ┌─────────────────────┐           │
│  │     Ingress         │           │
│  └─────────────────────┘           │
└─────────────────────────────────────┘
```

### 4.2 DNS trong Kubernetes
```bash
# Service DNS
<service-name>.<namespace>.svc.cluster.local

# Example
user-service.microservices.svc.cluster.local
```

### 4.3 Network Policies
- Default: All traffic allowed
- Network Policies: Implement micro-segmentation
- Rules: Ingress/Egress, namespaceSelector, podSelector

## 5. Storage và ConfigMaps

### 5.1 ConfigMaps vs Secrets
```yaml
# ConfigMap - Non-sensitive data
apiVersion: v1
kind: ConfigMap
data:
  database_url: "postgresql://..."
  
# Secret - Sensitive data (base64 encoded)
apiVersion: v1
kind: Secret
data:
  password: cGFzc3dvcmQ=  # password in base64
```

### 5.2 Volume Types
- **emptyDir**: Temporary storage
- **hostPath**: Node filesystem
- **persistentVolume**: Persistent storage
- **configMap/secret**: Configuration injection

## 6. Resource Management

### 6.1 Resource Requests vs Limits
```yaml
resources:
  requests:    # Minimum guaranteed
    cpu: 100m  # 0.1 CPU
    memory: 128Mi
  limits:      # Maximum allowed
    cpu: 500m  # 0.5 CPU
    memory: 512Mi
```

### 6.2 Horizontal Pod Autoscaler (HPA)
```yaml
# Scales pods based on metrics
spec:
  minReplicas: 1
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        averageUtilization: 70
```

### 6.3 Quality of Service (QoS)
- **Guaranteed**: requests = limits
- **Burstable**: requests < limits
- **BestEffort**: no requests/limits

## 7. Security Best Practices

### 7.1 Pod Security Standards
```yaml
securityContext:
  runAsNonRoot: true      # Don't run as root
  runAsUser: 1000         # Specific user ID
  readOnlyRootFilesystem: true
  allowPrivilegeEscalation: false
  capabilities:
    drop:
    - ALL                 # Drop all capabilities
```

### 7.2 RBAC (Role-Based Access Control)
```yaml
# Role định nghĩa permissions
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
rules:
- apiGroups: [""]
  resources: ["pods"]
  verbs: ["get", "list"]

# RoleBinding gán role cho user/serviceaccount
kind: RoleBinding
subjects:
- kind: ServiceAccount
  name: my-service-account
roleRef:
  kind: Role
  name: my-role
```

### 7.3 Network Policies
```yaml
# Default deny all
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
spec:
  podSelector: {}
  policyTypes:
  - Ingress
  - Egress
```

## 8. Observability

### 8.1 Three Pillars of Observability
```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   Metrics   │  │   Logging   │  │   Tracing   │
│             │  │             │  │             │
│ • CPU/Memory│  │ • App logs  │  │ • Request   │
│ • Request   │  │ • Error logs│  │   flow      │
│   rate      │  │ • Audit logs│  │ • Latency   │
│ • Error rate│  │             │  │   analysis  │
└─────────────┘  └─────────────┘  └─────────────┘
```

### 8.2 Prometheus Architecture
```
┌──────────────┐    ┌─────────────┐    ┌─────────────┐
│  Application │───►│ Prometheus  │───►│   Grafana   │
│  (metrics)   │    │  (scrape)   │    │ (visualize) │
└──────────────┘    └─────────────┘    └─────────────┘
                           │
                    ┌─────────────┐
                    │ Alertmanager│
                    │ (notify)    │
                    └─────────────┘
```

### 8.3 Key Metrics to Monitor
- **Golden Signals**:
  - Latency (response time)
  - Traffic (request rate)
  - Errors (error rate)
  - Saturation (resource utilization)

## 9. CI/CD Pipeline

### 9.1 GitOps Workflow
```
Developer → Git Push → CI Pipeline → Build Image → Update Manifest → CD Pipeline → Deploy
     ↑                                                                      ↓
     └──────────────────── Feedback Loop ←──────────────────────────────────┘
```

### 9.2 Pipeline Stages
```yaml
# 1. Code Quality
- Linting
- Unit Tests
- Security Scan

# 2. Build
- Docker Build
- Image Scan
- Push to Registry

# 3. Deploy
- Update Manifests
- Apply to Cluster
- Health Checks

# 4. Post-Deploy
- Integration Tests
- Performance Tests
- Monitoring Setup
```

### 9.3 Deployment Strategies
- **Rolling Update**: Gradual replacement (default)
- **Blue-Green**: Switch between environments
- **Canary**: Small percentage first
- **A/B Testing**: Feature flags

## 10. Helm Package Manager

### 10.1 Helm Concepts
```
Chart = Package (templates + values)
Release = Deployed instance
Repository = Chart collection
```

### 10.2 Template Structure
```
mychart/
├── Chart.yaml       # Chart metadata
├── values.yaml      # Default values
├── templates/       # Kubernetes manifests
│   ├── deployment.yaml
│   ├── service.yaml
│   └── _helpers.tpl # Template helpers
└── charts/         # Dependencies
```

### 10.3 Templating
```yaml
# Template with values
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Values.name }}
spec:
  replicas: {{ .Values.replicaCount }}
```

## 11. Production Readiness

### 11.1 High Availability
- Multiple master nodes
- Multiple worker nodes
- Load balancers
- Data replication

### 11.2 Disaster Recovery
- Backup strategies
- Recovery procedures
- Testing DR plans
- RTO/RPO targets

### 11.3 Monitoring và Alerting
```yaml
# Critical Alerts
- Pod crashes
- High CPU/Memory
- Disk space low
- Network issues
- Certificate expiry
```

## 12. Common Patterns

### 12.1 Microservices Patterns
- **Sidecar**: Helper container
- **Ambassador**: Proxy container
- **Adapter**: Format converter

### 12.2 Configuration Patterns
- **ConfigMap**: Environment-specific config
- **Secrets**: Sensitive data
- **Init Containers**: Setup tasks

### 12.3 Traffic Management
- **Service Mesh**: Istio, Linkerd
- **Ingress Controllers**: NGINX, Traefik
- **Load Balancing**: Round-robin, sticky sessions

## Next Steps

Sau khi hiểu các khái niệm này, bạn có thể:
1. Thực hành với cluster local (minikube/kind)
2. Deploy ứng dụng đơn giản
3. Implement monitoring
4. Setup CI/CD pipeline
5. Apply security best practices

Bạn muốn tôi giải thích sâu hơn về phần nào không? 