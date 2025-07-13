# CI/CD Theory & Best Practices Guide

## 1. CI/CD Fundamentals

### 1.1 Definitions
```
┌─────────────────────────────────────────────────────────────┐
│                    DevOps Pipeline                         │
├─────────────────────────────────────────────────────────────┤
│  CI (Continuous Integration)    │  CD (Continuous Delivery) │
│  ┌─────────────────────────┐    │  ┌─────────────────────┐  │
│  │ • Code commits          │    │  │ • Automated deploy │  │
│  │ • Automated builds      │    │  │ • Testing in prod  │  │
│  │ • Automated tests       │    │  │ • Release ready    │  │
│  │ • Code quality checks   │    │  │ • Manual approval  │  │
│  └─────────────────────────┘    │  └─────────────────────┘  │
│                                 │                           │
│                                 │  CD (Continuous Deploy)  │
│                                 │  ┌─────────────────────┐  │
│                                 │  │ • Fully automated   │  │
│                                 │  │ • No manual steps   │  │
│                                 │  │ • Direct to prod    │  │
│                                 │  └─────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 1.2 Benefits
- **Faster delivery**: Reduced time to market
- **Better quality**: Automated testing catches issues early
- **Reduced risk**: Smaller, frequent deployments
- **Improved collaboration**: Shared responsibility
- **Faster feedback**: Quick detection of issues

## 2. CI/CD Pipeline Architecture

### 2.1 Pipeline Stages
```
┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐
│ Source  │→ │  Build  │→ │  Test   │→ │ Deploy  │→ │ Monitor │
│ Control │  │         │  │         │  │         │  │         │
└─────────┘  └─────────┘  └─────────┘  └─────────┘  └─────────┘
     │            │            │            │            │
     ▼            ▼            ▼            ▼            ▼
┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐
│ • Git   │  │ • Compile│ │ • Unit   │ │ • Staging│ │ • Metrics│
│ • Branch│  │ • Package│ │ • Integr.│ │ • Prod   │ │ • Alerts │
│ • PR    │  │ • Docker │ │ • E2E    │ │ • Rollback│ │ • Logs  │
└─────────┘  └─────────┘  └─────────┘  └─────────┘  └─────────┘
```

### 2.2 Pipeline Triggers
- **Push triggers**: Code commits, PR merges
- **Schedule triggers**: Nightly builds, maintenance
- **Manual triggers**: Release deployments
- **Event triggers**: Dependency updates, security patches

### 2.3 Environment Strategy
```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│ Development │→ │   Staging   │→ │    UAT      │→ │ Production  │
├─────────────┤  ├─────────────┤  ├─────────────┤  ├─────────────┤
│ • Local dev │  │ • Integration│ │ • User      │  │ • Live      │
│ • Feature   │  │   testing   │ │   acceptance │ │   traffic   │
│   branches  │  │ • Auto      │ │ • Manual    │  │ • Blue/Green│
│ • Fast      │  │   deploy    │ │   testing   │ │ • Monitoring│
│   feedback  │  │ • Smoke     │ │ • Business  │  │ • Rollback  │
│             │  │   tests     │ │   validation │ │   ready     │
└─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘
```

## 3. Source Control Best Practices

### 3.1 Branching Strategies

#### **GitFlow**
```
master  ────●─────●─────●────  (Production releases)
           /     /     /
release   ●─────●─────●        (Release candidates)
         /     /     /
develop ●─────●─────●─────●    (Integration branch)
       /     /     /     /
feature ●───●     ●─────●      (Feature development)
```

#### **GitHub Flow**
```
main    ────●─────●─────●────  (Always deployable)
           /     /     /
feature   ●─────●     ●─────●  (Feature branches)
```

#### **GitLab Flow**
```
main     ────●─────●─────●────  (Source of truth)
            /     /     /
pre-prod   ●─────●─────●        (Staging environment)
          /     /     /
prod     ●─────●─────●          (Production environment)
```

### 3.2 Commit Conventions
```bash
# Conventional Commits
<type>[optional scope]: <description>

# Examples
feat(auth): add JWT token validation
fix(api): resolve memory leak in user service
docs(readme): update deployment instructions
test(user): add unit tests for registration
refactor(db): optimize query performance
```

### 3.3 Pull Request Process
```
1. Create feature branch
2. Implement changes
3. Add tests
4. Update documentation
5. Create pull request
6. Code review
7. Automated checks
8. Merge to main branch
```

## 4. Build Automation

### 4.1 Build Pipeline Steps
```yaml
# Example GitHub Actions
name: CI Pipeline
on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.23'
    - name: Run tests
      run: go test -v ./...
    - name: Code coverage
      run: go test -coverprofile=coverage.out ./...
```

### 4.2 Build Optimization
- **Parallel execution**: Run tests in parallel
- **Caching**: Cache dependencies and build artifacts
- **Incremental builds**: Only build changed components
- **Build matrix**: Test across multiple environments

### 4.3 Artifact Management
```
Build Artifacts:
├── Binaries
├── Docker Images
├── Configuration Files
├── Database Migrations
├── Documentation
└── Test Reports
```

## 5. Testing Strategy

### 5.1 Testing Pyramid
```
                ┌─────────────┐
                │     E2E     │  ← Few, Slow, Expensive
                │   Tests     │
              ┌─┴─────────────┴─┐
              │  Integration    │  ← Some, Medium
              │     Tests       │
            ┌─┴─────────────────┴─┐
            │     Unit Tests      │  ← Many, Fast, Cheap
            └─────────────────────┘
```

### 5.2 Test Types

#### **Unit Tests**
```go
func TestUserRegistration(t *testing.T) {
    user := &User{
        Email:    "test@example.com",
        Password: "password123",
    }
    
    err := user.Validate()
    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", user.Email)
}
```

#### **Integration Tests**
```go
func TestUserServiceIntegration(t *testing.T) {
    // Setup test database
    db := setupTestDB()
    defer cleanupTestDB(db)
    
    service := NewUserService(db)
    user, err := service.CreateUser("test@example.com", "password")
    
    assert.NoError(t, err)
    assert.NotEmpty(t, user.ID)
}
```

#### **End-to-End Tests**
```javascript
describe('User Registration Flow', () => {
  it('should register new user successfully', async () => {
    await page.goto('/register');
    await page.fill('[name="email"]', 'test@example.com');
    await page.fill('[name="password"]', 'password123');
    await page.click('[type="submit"]');
    
    await expect(page).toHaveURL('/dashboard');
  });
});
```

### 5.3 Test Data Management
- **Test fixtures**: Predefined test data
- **Data factories**: Generate test data programmatically
- **Database seeding**: Initialize test databases
- **Mock services**: Simulate external dependencies

## 6. Quality Gates

### 6.1 Static Code Analysis
```yaml
# SonarQube Quality Gate
quality_gate:
  conditions:
    - coverage > 80%
    - duplicated_lines_density < 3%
    - security_rating == A
    - maintainability_rating == A
    - reliability_rating == A
```

### 6.2 Security Scanning
```yaml
# Security Pipeline
security_checks:
  - name: Dependency Check
    tool: OWASP Dependency Check
    action: fail_on_high_severity
    
  - name: Static Analysis
    tool: Gosec
    action: report_and_continue
    
  - name: Container Scan
    tool: Trivy
    action: fail_on_critical
```

### 6.3 Performance Testing
```yaml
# Performance Gates
performance_criteria:
  response_time_p95: < 200ms
  throughput: > 1000 RPS
  error_rate: < 0.1%
  cpu_usage: < 70%
  memory_usage: < 80%
```

## 7. Deployment Strategies

### 7.1 Blue-Green Deployment
```
┌─────────────────┐    ┌─────────────────┐
│   Blue Env      │    │   Green Env     │
│  (Current)      │    │   (New)         │
├─────────────────┤    ├─────────────────┤
│ • Version 1.0   │    │ • Version 2.0   │
│ • Live traffic  │    │ • Warming up    │
│ • Monitoring    │    │ • Testing       │
└─────────────────┘    └─────────────────┘
         │                       │
         └──── Switch Traffic ───┘
```

**Pros**: Instant rollback, zero downtime
**Cons**: Resource intensive, requires duplicate infrastructure

### 7.2 Rolling Deployment
```
Initial:  [V1] [V1] [V1] [V1]
Step 1:   [V2] [V1] [V1] [V1]  ← Update one instance
Step 2:   [V2] [V2] [V1] [V1]  ← Update next instance
Step 3:   [V2] [V2] [V2] [V1]  ← Continue rolling
Final:    [V2] [V2] [V2] [V2]  ← All updated
```

**Pros**: Resource efficient, gradual migration
**Cons**: Mixed versions during deployment

### 7.3 Canary Deployment
```
Production Traffic: 100%
                   ↓
            ┌─────────────┐
            │Load Balancer│
            └─────────────┘
                    │
           ┌────────┴────────┐
           ▼                 ▼
    ┌─────────────┐   ┌─────────────┐
    │Stable (95%) │   │Canary (5%)  │
    │Version 1.0  │   │Version 2.0  │
    └─────────────┘   └─────────────┘
```

**Pros**: Low risk, real user feedback
**Cons**: Complex routing, monitoring required

### 7.4 A/B Testing
```
User Traffic
     │
┌────┴────┐
│ Router  │
└────┬────┘
     │
┌────┴─────────────┐
▼                  ▼
Version A (50%)    Version B (50%)
Feature Flag OFF   Feature Flag ON
```

## 8. Infrastructure as Code

### 8.1 IaC Principles
- **Declarative**: Describe desired state
- **Idempotent**: Same result every time
- **Version controlled**: Track changes
- **Immutable**: Replace rather than modify

### 8.2 Kubernetes Manifests
```yaml
# Declarative deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: user-service
spec:
  replicas: 3
  selector:
    matchLabels:
      app: user-service
  template:
    metadata:
      labels:
        app: user-service
    spec:
      containers:
      - name: user-service
        image: myapp/user-service:v1.2.3
        resources:
          requests:
            cpu: 100m
            memory: 128Mi
          limits:
            cpu: 500m
            memory: 512Mi
```

### 8.3 GitOps Workflow
```
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   Source    │→ │    Build    │→ │   Config    │→ │   Deploy    │
│    Repo     │  │   Pipeline  │  │    Repo     │  │   Agent     │
├─────────────┤  ├─────────────┤  ├─────────────┤  ├─────────────┤
│ • App code  │  │ • Build     │  │ • K8s       │  │ • ArgoCD    │
│ • Tests     │  │ • Test      │  │   manifests │  │ • Flux      │
│ • Dockerfile│  │ • Scan      │  │ • Helm      │  │ • Tekton    │
└─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘
```

## 9. Security in CI/CD

### 9.1 Pipeline Security
```yaml
# Secure pipeline practices
security_measures:
  secrets_management:
    - use_vault_or_secrets_manager
    - no_secrets_in_code
    - rotate_regularly
    
  access_control:
    - principle_of_least_privilege
    - multi_factor_authentication
    - audit_logs
    
  supply_chain:
    - dependency_scanning
    - base_image_scanning
    - signed_artifacts
```

### 9.2 Security Gates
```yaml
# Security checkpoints
gates:
  - stage: source
    checks:
      - secrets_detection
      - dependency_vulnerabilities
      
  - stage: build
    checks:
      - static_analysis
      - license_compliance
      
  - stage: deploy
    checks:
      - container_scanning
      - runtime_security
```

### 9.3 Compliance & Governance
- **SOC 2**: Security controls and procedures
- **PCI DSS**: Payment card industry standards
- **GDPR**: Data protection regulations
- **HIPAA**: Healthcare data security

## 10. Monitoring & Observability

### 10.1 Pipeline Monitoring
```
Build Metrics:
├── Build duration
├── Success rate
├── Deployment frequency
├── Lead time
├── Mean time to recovery
└── Change failure rate
```

### 10.2 DORA Metrics
```
┌─────────────────────────────────────────┐
│              DORA Metrics               │
├─────────────────────────────────────────┤
│ 1. Deployment Frequency                 │
│    How often: Daily, Weekly, Monthly    │
│                                         │
│ 2. Lead Time for Changes                │
│    Time: Commit → Production            │
│                                         │
│ 3. Change Failure Rate                  │
│    Percentage: Failed deployments      │
│                                         │
│ 4. Mean Time to Recovery (MTTR)         │
│    Time: Incident → Resolution          │
└─────────────────────────────────────────┘
```

### 10.3 Pipeline Alerting
```yaml
alerts:
  - name: Build Failure
    condition: build_status == "failed"
    severity: high
    notification: slack, email
    
  - name: Deployment Timeout
    condition: deployment_duration > 30m
    severity: critical
    notification: pagerduty
    
  - name: High Error Rate
    condition: error_rate > 5%
    severity: warning
    notification: slack
```

## 11. Advanced CI/CD Patterns

### 11.1 Microservices CI/CD
```
┌─────────────────────────────────────────┐
│            Monorepo                     │
├─────────────────────────────────────────┤
│ Service A │ Service B │ Service C       │
│     │     │     │     │     │           │
│     ▼     │     ▼     │     ▼           │
│ Pipeline A│ Pipeline B│ Pipeline C      │
│     │     │     │     │     │           │
│     ▼     │     ▼     │     ▼           │
│  Deploy A │  Deploy B │  Deploy C       │
└─────────────────────────────────────────┘
```

### 11.2 Feature Flags
```go
// Feature flag example
func handleUserRegistration(w http.ResponseWriter, r *http.Request) {
    if featureFlag.IsEnabled("new_registration_flow", userID) {
        // New registration logic
        handleNewRegistration(w, r)
    } else {
        // Legacy registration logic
        handleLegacyRegistration(w, r)
    }
}
```

### 11.3 Progressive Delivery
```
┌─────────────────┐
│ Feature Flag    │ ← Control exposure
├─────────────────┤
│ Canary Deploy   │ ← Gradual rollout
├─────────────────┤
│ A/B Testing     │ ← Measure impact
├─────────────────┤
│ Full Rollout    │ ← Complete deployment
└─────────────────┘
```

## 12. Troubleshooting & Optimization

### 12.1 Common Pipeline Issues
```
Pipeline Failures:
├── Environment Issues
│   ├── Network connectivity
│   ├── Resource constraints
│   └── Permission problems
├── Code Issues
│   ├── Test failures
│   ├── Build errors
│   └── Dependency conflicts
├── Infrastructure Issues
│   ├── Kubernetes cluster problems
│   ├── Registry issues
│   └── Load balancer problems
└── Configuration Issues
    ├── Wrong environment variables
    ├── Missing secrets
    └── Incorrect manifest files
```

### 12.2 Performance Optimization
```yaml
# Pipeline optimization strategies
optimization:
  parallel_execution:
    - run_tests_in_parallel
    - build_multiple_architectures
    
  caching:
    - cache_dependencies
    - cache_docker_layers
    - cache_build_artifacts
    
  resource_management:
    - right_size_runners
    - use_spot_instances
    - cleanup_resources
```

### 12.3 Cost Optimization
- **Resource right-sizing**: Match resources to workload
- **Spot instances**: Use cheaper compute when available
- **Pipeline efficiency**: Reduce build times
- **Artifact cleanup**: Remove old artifacts regularly

## Best Practices Summary

### ✅ Do's
- Implement comprehensive testing
- Use infrastructure as code
- Monitor and measure everything
- Automate security scanning
- Practice gitops principles
- Implement proper secret management
- Use feature flags for risk mitigation
- Document processes and runbooks

### ❌ Don'ts
- Skip testing to go faster
- Store secrets in code
- Deploy directly to production
- Ignore security vulnerabilities
- Deploy without monitoring
- Skip rollback planning
- Use long-lived branches
- Deploy without proper documentation

Bạn muốn tôi giải thích chi tiết hơn về phần nào trong CI/CD không? 