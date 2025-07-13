# Configuration Management

This directory contains environment-specific configuration files.

## Structure

```
configs/
├── staging/          # Staging environment configurations
│   ├── app.yaml     # Application configuration
│   ├── database.yaml # Database configuration
│   └── secrets.yaml  # Secrets (encrypted)
├── production/      # Production environment configurations
│   ├── app.yaml     # Application configuration
│   ├── database.yaml # Database configuration
│   └── secrets.yaml  # Secrets (encrypted)
└── common/          # Common configurations
    └── logging.yaml  # Logging configuration
```

## Best Practices

### 1. **Environment Separation**
- Keep staging and production configurations separate
- Use different namespaces/clusters for isolation
- Test configurations in staging before production

### 2. **Secret Management**
- Never store secrets in plain text
- Use tools like Sealed Secrets, External Secrets Operator
- Rotate secrets regularly

### 3. **Configuration Validation**
- Validate configurations before deployment
- Use schema validation where possible
- Test configuration changes in staging

### 4. **Version Control**
- Version control all configuration files
- Use meaningful commit messages
- Tag releases for traceability

## Tools

- **Kustomize**: For managing configuration variants
- **Helm**: For templating and packaging
- **Sealed Secrets**: For encrypting secrets
- **External Secrets**: For external secret management

## Usage

1. **Deploy to staging:**
   ```bash
   kubectl apply -f staging/ -n staging
   ```

2. **Deploy to production:**
   ```bash
   kubectl apply -f production/ -n production
   ```

3. **Validate configuration:**
   ```bash
   kubectl dry-run=client -f staging/
   ``` 