# Terraform Infrastructure

This directory contains Terraform configurations for provisioning infrastructure.

## Structure

```
terraform/
├── modules/           # Reusable Terraform modules
│   ├── eks/          # EKS cluster module
│   ├── vpc/          # VPC module
│   └── rds/          # Database module
├── environments/     # Environment-specific configurations
│   ├── staging/      # Staging environment
│   └── production/   # Production environment
└── shared/          # Shared resources
```

## Usage

1. **Initialize Terraform:**
   ```bash
   cd environments/staging
   terraform init
   ```

2. **Plan deployment:**
   ```bash
   terraform plan -var-file="terraform.tfvars"
   ```

3. **Apply changes:**
   ```bash
   terraform apply -var-file="terraform.tfvars"
   ```

## Best Practices

- Use modules for reusable infrastructure components
- Separate environments with different state files
- Use terraform.tfvars for environment-specific variables
- Always run `terraform plan` before `terraform apply`
- Use remote state storage (S3 + DynamoDB) 