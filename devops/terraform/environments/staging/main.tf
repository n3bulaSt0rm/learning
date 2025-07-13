# Provider Configuration
terraform {
  required_version = ">= 1.0"
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
  }

  # Backend configuration for remote state
  backend "s3" {
    bucket         = "your-terraform-state-bucket-name"
    key            = "staging/eks/terraform.tfstate"
    region         = "us-west-2"
    dynamodb_table = "terraform-state-lock"
    encrypt        = true
  }
}

# AWS Provider
provider "aws" {
  region = var.aws_region
  
  default_tags {
    tags = {
      Environment = "staging"
      Project     = "learning-microservices"
      ManagedBy   = "terraform"
    }
  }
}

# Local Values
locals {
  name         = "learning-staging"
  cluster_name = "learning-eks-staging"
  region       = var.aws_region
  
  common_tags = {
    Environment = "staging"
    Project     = "learning-microservices"
    ManagedBy   = "terraform"
    Owner       = "devops-team"
  }
}

# Data Sources
data "aws_caller_identity" "current" {}
data "aws_availability_zones" "available" {}

# VPC Module
module "vpc" {
  source = "../../modules/vpc"

  name         = local.name
  region       = local.region
  cluster_name = local.cluster_name

  vpc_cidr             = var.vpc_cidr
  public_subnet_cidrs  = var.public_subnet_cidrs
  private_subnet_cidrs = var.private_subnet_cidrs

  enable_nat_gateway   = true
  single_nat_gateway   = true  # For cost optimization in staging
  enable_dns_hostnames = true
  enable_dns_support   = true

  # VPC Endpoints for cost optimization
  enable_s3_endpoint  = true
  enable_ec2_endpoint = true
  enable_ecr_endpoint = true

  # Flow logs
  enable_flow_logs          = true
  flow_log_retention_days   = 7  # Shorter retention for staging

  tags = local.common_tags
}

# EKS Module
module "eks" {
  source = "../../modules/eks"

  cluster_name       = local.cluster_name
  kubernetes_version = var.kubernetes_version
  region            = local.region

  # VPC Configuration
  vpc_id              = module.vpc.vpc_id
  subnet_ids          = module.vpc.eks_subnet_ids
  private_subnet_ids  = module.vpc.eks_private_subnet_ids

  # Cluster Access
  endpoint_private_access = true
  endpoint_public_access  = true
  public_access_cidrs     = var.public_access_cidrs

  # Node Group Configuration
  instance_types   = var.instance_types
  capacity_type    = "SPOT"  # Use spot instances for cost savings
  desired_size     = 2
  min_size         = 1
  max_size         = 5
  disk_size        = 30

  # Fargate Configuration
  enable_fargate     = var.enable_fargate
  fargate_namespace  = "default"

  # Add-ons
  enable_irsa                         = true
  enable_cluster_autoscaler           = true
  enable_aws_load_balancer_controller = true

  # Logging
  enabled_cluster_log_types = ["api", "audit", "authenticator", "controllerManager", "scheduler"]
  log_retention_days        = 7

  # Monitoring
  enable_cloudwatch_logging    = true
  enable_prometheus_monitoring = true

  environment = "staging"
  tags        = local.common_tags
}

# ECR Repository for microservices
resource "aws_ecr_repository" "microservices" {
  for_each = toset(var.microservices)

  name                 = "${local.name}-${each.value}"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    scan_on_push = true
  }

  lifecycle_policy {
    policy = jsonencode({
      rules = [
        {
          rulePriority = 1
          description  = "Keep last 10 images"
          selection = {
            tagStatus     = "tagged"
            tagPrefixList = ["v"]
            countType     = "imageCountMoreThan"
            countNumber   = 10
          }
          action = {
            type = "expire"
          }
        }
      ]
    })
  }

  tags = local.common_tags
}

# RDS Subnet Group for databases
resource "aws_db_subnet_group" "main" {
  name       = "${local.name}-db-subnet-group"
  subnet_ids = module.vpc.private_subnet_ids

  tags = merge(local.common_tags, {
    Name = "${local.name}-db-subnet-group"
  })
}

# Security Group for RDS
resource "aws_security_group" "rds" {
  name        = "${local.name}-rds-sg"
  description = "Security group for RDS database"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.common_tags, {
    Name = "${local.name}-rds-sg"
  })
}

# RDS PostgreSQL instance
resource "aws_db_instance" "main" {
  identifier             = "${local.name}-postgres"
  engine                 = "postgres"
  engine_version         = "15.4"
  instance_class         = "db.t3.micro"  # Small instance for staging
  allocated_storage      = 20
  max_allocated_storage  = 100
  storage_type           = "gp3"
  storage_encrypted      = true

  db_name  = "microservices"
  username = var.db_username
  password = var.db_password

  db_subnet_group_name   = aws_db_subnet_group.main.name
  vpc_security_group_ids = [aws_security_group.rds.id]

  backup_retention_period = 7
  backup_window          = "03:00-04:00"
  maintenance_window     = "sun:04:00-sun:05:00"

  skip_final_snapshot = true
  deletion_protection = false

  tags = merge(local.common_tags, {
    Name = "${local.name}-postgres"
  })
}

# ElastiCache Subnet Group
resource "aws_elasticache_subnet_group" "main" {
  name       = "${local.name}-cache-subnet-group"
  subnet_ids = module.vpc.private_subnet_ids

  tags = local.common_tags
}

# Security Group for ElastiCache
resource "aws_security_group" "elasticache" {
  name        = "${local.name}-elasticache-sg"
  description = "Security group for ElastiCache"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port   = 6379
    to_port     = 6379
    protocol    = "tcp"
    cidr_blocks = [module.vpc.vpc_cidr_block]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = merge(local.common_tags, {
    Name = "${local.name}-elasticache-sg"
  })
}

# ElastiCache Redis cluster
resource "aws_elasticache_replication_group" "main" {
  replication_group_id       = "${local.name}-redis"
  description                = "Redis cluster for microservices"
  port                       = 6379
  parameter_group_name       = "default.redis7"
  node_type                  = "cache.t3.micro"
  num_cache_clusters         = 1
  engine_version             = "7.0"
  
  subnet_group_name  = aws_elasticache_subnet_group.main.name
  security_group_ids = [aws_security_group.elasticache.id]

  at_rest_encryption_enabled = true
  transit_encryption_enabled = true

  tags = merge(local.common_tags, {
    Name = "${local.name}-redis"
  })
}

# S3 Bucket for application data
resource "aws_s3_bucket" "app_data" {
  bucket = "${local.name}-app-data-${random_id.bucket_suffix.hex}"

  tags = merge(local.common_tags, {
    Name = "${local.name}-app-data"
  })
}

resource "random_id" "bucket_suffix" {
  byte_length = 4
}

# S3 Bucket versioning
resource "aws_s3_bucket_versioning" "app_data" {
  bucket = aws_s3_bucket.app_data.id
  versioning_configuration {
    status = "Enabled"
  }
}

# S3 Bucket encryption
resource "aws_s3_bucket_server_side_encryption_configuration" "app_data" {
  bucket = aws_s3_bucket.app_data.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}

# S3 Bucket public access block
resource "aws_s3_bucket_public_access_block" "app_data" {
  bucket = aws_s3_bucket.app_data.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Output important values
output "cluster_info" {
  description = "EKS cluster information"
  value = {
    cluster_name     = module.eks.cluster_name
    cluster_endpoint = module.eks.cluster_endpoint
    cluster_version  = module.eks.cluster_version
    region           = local.region
    vpc_id           = module.vpc.vpc_id
    private_subnets  = module.vpc.private_subnet_ids
    public_subnets   = module.vpc.public_subnet_ids
  }
}

output "database_info" {
  description = "Database connection information"
  value = {
    rds_endpoint = aws_db_instance.main.endpoint
    redis_endpoint = aws_elasticache_replication_group.main.primary_endpoint_address
  }
  sensitive = true
}

output "aws_cli_commands" {
  description = "AWS CLI commands for cluster access"
  value = {
    update_kubeconfig = "aws eks update-kubeconfig --region ${local.region} --name ${local.cluster_name}"
    get_cluster_info  = "aws eks describe-cluster --region ${local.region} --name ${local.cluster_name}"
  }
} 