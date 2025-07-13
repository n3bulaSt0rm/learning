# Cluster Configuration
variable "cluster_name" {
  description = "Name of the EKS cluster"
  type        = string
}

variable "kubernetes_version" {
  description = "Kubernetes version for the EKS cluster"
  type        = string
  default     = "1.28"
}

# VPC Configuration
variable "vpc_id" {
  description = "ID of the VPC where the cluster will be created"
  type        = string
}

variable "subnet_ids" {
  description = "List of subnet IDs for the EKS cluster"
  type        = list(string)
}

variable "private_subnet_ids" {
  description = "List of private subnet IDs for the EKS node groups"
  type        = list(string)
}

# Cluster Access Configuration
variable "endpoint_private_access" {
  description = "Whether the Amazon EKS private API server endpoint is enabled"
  type        = bool
  default     = true
}

variable "endpoint_public_access" {
  description = "Whether the Amazon EKS public API server endpoint is enabled"
  type        = bool
  default     = true
}

variable "public_access_cidrs" {
  description = "List of CIDR blocks that can access the Amazon EKS public API server endpoint"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "allowed_cidr_blocks" {
  description = "List of CIDR blocks allowed to access the cluster"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

# Node Group Configuration
variable "capacity_type" {
  description = "Type of capacity associated with the EKS Node Group. Valid values: ON_DEMAND, SPOT"
  type        = string
  default     = "ON_DEMAND"
}

variable "instance_types" {
  description = "List of instance types for the EKS Node Group"
  type        = list(string)
  default     = ["t3.medium"]
}

variable "ami_type" {
  description = "Type of Amazon Machine Image (AMI) associated with the EKS Node Group"
  type        = string
  default     = "AL2_x86_64"
}

variable "disk_size" {
  description = "Disk size in GiB for worker nodes"
  type        = number
  default     = 20
}

variable "desired_size" {
  description = "Desired number of worker nodes"
  type        = number
  default     = 2
}

variable "max_size" {
  description = "Maximum number of worker nodes"
  type        = number
  default     = 10
}

variable "min_size" {
  description = "Minimum number of worker nodes"
  type        = number
  default     = 1
}

variable "max_unavailable" {
  description = "Maximum number of nodes unavailable at once during a version update"
  type        = number
  default     = 1
}

# Remote Access Configuration
variable "ec2_ssh_key" {
  description = "EC2 Key Pair name for SSH access to worker nodes"
  type        = string
  default     = null
}

variable "source_security_group_ids" {
  description = "Set of EC2 Security Group IDs to allow SSH access from"
  type        = list(string)
  default     = []
}

# Fargate Configuration
variable "enable_fargate" {
  description = "Whether to enable Fargate profile"
  type        = bool
  default     = false
}

variable "fargate_namespace" {
  description = "Kubernetes namespace for Fargate profile"
  type        = string
  default     = "default"
}

variable "fargate_labels" {
  description = "Key-value map of Kubernetes labels for Fargate profile"
  type        = map(string)
  default     = {}
}

# Logging Configuration
variable "enabled_cluster_log_types" {
  description = "List of control plane logging types to enable"
  type        = list(string)
  default     = ["api", "audit", "authenticator", "controllerManager", "scheduler"]
}

variable "log_retention_days" {
  description = "Number of days to retain CloudWatch logs"
  type        = number
  default     = 14
}

# Add-on Versions
variable "vpc_cni_version" {
  description = "Version of the VPC CNI add-on"
  type        = string
  default     = null
}

variable "coredns_version" {
  description = "Version of the CoreDNS add-on"
  type        = string
  default     = null
}

variable "kube_proxy_version" {
  description = "Version of the kube-proxy add-on"
  type        = string
  default     = null
}

variable "ebs_csi_driver_version" {
  description = "Version of the EBS CSI driver add-on"
  type        = string
  default     = null
}

# Security Configuration
variable "enable_irsa" {
  description = "Whether to enable IAM Roles for Service Accounts"
  type        = bool
  default     = true
}

variable "enable_pod_security_policy" {
  description = "Whether to enable Pod Security Policy"
  type        = bool
  default     = false
}

# Monitoring Configuration
variable "enable_cloudwatch_logging" {
  description = "Whether to enable CloudWatch Container Insights"
  type        = bool
  default     = true
}

variable "enable_prometheus_monitoring" {
  description = "Whether to enable Prometheus monitoring"
  type        = bool
  default     = true
}

# Backup Configuration
variable "enable_backup" {
  description = "Whether to enable AWS Backup for the cluster"
  type        = bool
  default     = false
}

variable "backup_retention_days" {
  description = "Number of days to retain backups"
  type        = number
  default     = 7
}

# Network Configuration
variable "cluster_security_group_additional_rules" {
  description = "Additional security group rules for the cluster security group"
  type        = any
  default     = {}
}

variable "node_security_group_additional_rules" {
  description = "Additional security group rules for the node security group"
  type        = any
  default     = {}
}

# Auto Scaling Configuration
variable "enable_cluster_autoscaler" {
  description = "Whether to enable cluster autoscaler"
  type        = bool
  default     = true
}

variable "cluster_autoscaler_version" {
  description = "Version of cluster autoscaler"
  type        = string
  default     = "1.21.0"
}

# Load Balancer Configuration
variable "enable_aws_load_balancer_controller" {
  description = "Whether to enable AWS Load Balancer Controller"
  type        = bool
  default     = true
}

variable "aws_load_balancer_controller_version" {
  description = "Version of AWS Load Balancer Controller"
  type        = string
  default     = "2.4.4"
}

# Storage Configuration
variable "enable_ebs_csi_driver" {
  description = "Whether to enable EBS CSI driver"
  type        = bool
  default     = true
}

variable "enable_efs_csi_driver" {
  description = "Whether to enable EFS CSI driver"
  type        = bool
  default     = false
}

# Cost Optimization
variable "enable_spot_instances" {
  description = "Whether to enable spot instances for cost optimization"
  type        = bool
  default     = false
}

variable "spot_instance_types" {
  description = "List of spot instance types"
  type        = list(string)
  default     = ["t3.medium", "t3.large", "t3.xlarge"]
}

variable "spot_max_price" {
  description = "Maximum price for spot instances"
  type        = string
  default     = "0.05"
}

# Environment
variable "environment" {
  description = "Environment name (e.g., dev, staging, prod)"
  type        = string
  default     = "dev"
}

variable "region" {
  description = "AWS region"
  type        = string
}

# Tags
variable "tags" {
  description = "A map of tags to add to all resources"
  type        = map(string)
  default     = {}
}

variable "cluster_tags" {
  description = "A map of tags to add to the cluster"
  type        = map(string)
  default     = {}
}

variable "node_group_tags" {
  description = "A map of tags to add to the node group"
  type        = map(string)
  default     = {}
}

# Multi-AZ Configuration
variable "availability_zones" {
  description = "List of availability zones"
  type        = list(string)
  default     = []
}

# Additional Node Groups
variable "additional_node_groups" {
  description = "Map of additional node group configurations"
  type        = any
  default     = {}
}

# Networking
variable "enable_nat_gateway" {
  description = "Whether to enable NAT Gateway"
  type        = bool
  default     = true
}

variable "single_nat_gateway" {
  description = "Whether to use single NAT Gateway"
  type        = bool
  default     = false
}

# DNS Configuration
variable "enable_dns_hostnames" {
  description = "Whether to enable DNS hostnames in the VPC"
  type        = bool
  default     = true
}

variable "enable_dns_support" {
  description = "Whether to enable DNS support in the VPC"
  type        = bool
  default     = true
}

# Encryption
variable "enable_encryption_at_rest" {
  description = "Whether to enable encryption at rest"
  type        = bool
  default     = true
}

variable "kms_key_id" {
  description = "KMS key ID for encryption"
  type        = string
  default     = null
}

# Compliance
variable "enable_compliance_monitoring" {
  description = "Whether to enable compliance monitoring"
  type        = bool
  default     = false
}

variable "compliance_frameworks" {
  description = "List of compliance frameworks to monitor"
  type        = list(string)
  default     = ["SOC2", "PCI-DSS"]
} 