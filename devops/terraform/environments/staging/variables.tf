# AWS Configuration
variable "aws_region" {
  description = "AWS region for resources"
  type        = string
  default     = "us-west-2"
}

# VPC Configuration
variable "vpc_cidr" {
  description = "CIDR block for VPC"
  type        = string
  default     = "10.0.0.0/16"
}

variable "public_subnet_cidrs" {
  description = "CIDR blocks for public subnets"
  type        = list(string)
  default     = ["10.0.1.0/24", "10.0.2.0/24"]
}

variable "private_subnet_cidrs" {
  description = "CIDR blocks for private subnets"
  type        = list(string)
  default     = ["10.0.11.0/24", "10.0.12.0/24"]
}

# EKS Configuration
variable "kubernetes_version" {
  description = "Kubernetes version for EKS cluster"
  type        = string
  default     = "1.28"
}

variable "instance_types" {
  description = "EC2 instance types for EKS worker nodes"
  type        = list(string)
  default     = ["t3.medium", "t3.large"]
}

variable "public_access_cidrs" {
  description = "CIDR blocks that can access the EKS public API server endpoint"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "enable_fargate" {
  description = "Enable Fargate profile for EKS"
  type        = bool
  default     = false
}

# Microservices Configuration
variable "microservices" {
  description = "List of microservices to create ECR repositories for"
  type        = list(string)
  default     = ["user-service", "product-service", "order-service", "api-gateway"]
}

# Database Configuration
variable "db_username" {
  description = "Database administrator username"
  type        = string
  default     = "postgres"
}

variable "db_password" {
  description = "Database administrator password"
  type        = string
  sensitive   = true
}

# Environment Configuration
variable "environment" {
  description = "Environment name"
  type        = string
  default     = "staging"
} 