# Cluster Information
output "cluster_id" {
  description = "The ID of the EKS cluster"
  value       = aws_eks_cluster.main.id
}

output "cluster_arn" {
  description = "The Amazon Resource Name (ARN) of the cluster"
  value       = aws_eks_cluster.main.arn
}

output "cluster_name" {
  description = "The name of the EKS cluster"
  value       = aws_eks_cluster.main.name
}

output "cluster_endpoint" {
  description = "Endpoint for your Kubernetes API server"
  value       = aws_eks_cluster.main.endpoint
}

output "cluster_version" {
  description = "The Kubernetes version for the cluster"
  value       = aws_eks_cluster.main.version
}

output "cluster_platform_version" {
  description = "Platform version for the cluster"
  value       = aws_eks_cluster.main.platform_version
}

output "cluster_status" {
  description = "Status of the EKS cluster. One of `CREATING`, `ACTIVE`, `DELETING`, `FAILED`"
  value       = aws_eks_cluster.main.status
}

# Cluster Authentication
output "cluster_certificate_authority_data" {
  description = "Base64 encoded certificate data required to communicate with the cluster"
  value       = aws_eks_cluster.main.certificate_authority[0].data
}

output "cluster_oidc_issuer_url" {
  description = "The URL on the EKS cluster OIDC Issuer"
  value       = aws_eks_cluster.main.identity[0].oidc[0].issuer
}

output "oidc_provider_arn" {
  description = "The ARN of the OIDC Provider"
  value       = aws_iam_openid_connect_provider.eks.arn
}

# IAM Roles
output "cluster_iam_role_name" {
  description = "IAM role name associated with EKS cluster"
  value       = aws_iam_role.eks_cluster_role.name
}

output "cluster_iam_role_arn" {
  description = "IAM role ARN associated with EKS cluster"
  value       = aws_iam_role.eks_cluster_role.arn
}

output "node_group_iam_role_name" {
  description = "IAM role name associated with EKS node group"
  value       = aws_iam_role.eks_node_group_role.name
}

output "node_group_iam_role_arn" {
  description = "IAM role ARN associated with EKS node group"
  value       = aws_iam_role.eks_node_group_role.arn
}

# Node Group Information
output "node_group_id" {
  description = "EKS node group ID"
  value       = aws_eks_node_group.main.id
}

output "node_group_arn" {
  description = "Amazon Resource Name (ARN) of the EKS Node Group"
  value       = aws_eks_node_group.main.arn
}

output "node_group_status" {
  description = "Status of the EKS Node Group"
  value       = aws_eks_node_group.main.status
}

output "node_group_capacity_type" {
  description = "Type of capacity associated with the EKS Node Group"
  value       = aws_eks_node_group.main.capacity_type
}

output "node_group_instance_types" {
  description = "List of instance types associated with EKS Node Group"
  value       = aws_eks_node_group.main.instance_types
}

output "node_group_ami_type" {
  description = "Type of Amazon Machine Image (AMI) associated with the EKS Node Group"
  value       = aws_eks_node_group.main.ami_type
}

output "node_group_disk_size" {
  description = "Disk size in GiB for worker nodes"
  value       = aws_eks_node_group.main.disk_size
}

output "node_group_resources" {
  description = "List of objects containing information about underlying resources"
  value       = aws_eks_node_group.main.resources
}

# Fargate Profile Information
output "fargate_profile_id" {
  description = "EKS Fargate Profile ID"
  value       = var.enable_fargate ? aws_eks_fargate_profile.main[0].id : null
}

output "fargate_profile_arn" {
  description = "Amazon Resource Name (ARN) of the EKS Fargate Profile"
  value       = var.enable_fargate ? aws_eks_fargate_profile.main[0].arn : null
}

output "fargate_profile_status" {
  description = "Status of the EKS Fargate Profile"
  value       = var.enable_fargate ? aws_eks_fargate_profile.main[0].status : null
}

# Security Groups
output "cluster_security_group_id" {
  description = "ID of the cluster security group"
  value       = aws_security_group.eks_cluster.id
}

output "cluster_security_group_arn" {
  description = "ARN of the cluster security group"
  value       = aws_security_group.eks_cluster.arn
}

output "cluster_primary_security_group_id" {
  description = "Cluster security group that was created by Amazon EKS for the cluster"
  value       = aws_eks_cluster.main.vpc_config[0].cluster_security_group_id
}

# CloudWatch Log Group
output "cloudwatch_log_group_name" {
  description = "Name of cloudwatch log group for EKS cluster"
  value       = aws_cloudwatch_log_group.eks.name
}

output "cloudwatch_log_group_arn" {
  description = "ARN of cloudwatch log group for EKS cluster"
  value       = aws_cloudwatch_log_group.eks.arn
}

# KMS Key
output "kms_key_arn" {
  description = "The Amazon Resource Name (ARN) of the KMS key"
  value       = aws_kms_key.eks.arn
}

output "kms_key_id" {
  description = "The globally unique identifier for the KMS key"
  value       = aws_kms_key.eks.key_id
}

# Add-ons
output "cluster_addons" {
  description = "Map of cluster add-on configurations"
  value = {
    vpc_cni = {
      addon_name    = aws_eks_addon.vpc_cni.addon_name
      addon_version = aws_eks_addon.vpc_cni.addon_version
      status        = aws_eks_addon.vpc_cni.status
    }
    coredns = {
      addon_name    = aws_eks_addon.coredns.addon_name
      addon_version = aws_eks_addon.coredns.addon_version
      status        = aws_eks_addon.coredns.status
    }
    kube_proxy = {
      addon_name    = aws_eks_addon.kube_proxy.addon_name
      addon_version = aws_eks_addon.kube_proxy.addon_version
      status        = aws_eks_addon.kube_proxy.status
    }
    ebs_csi_driver = {
      addon_name    = aws_eks_addon.ebs_csi_driver.addon_name
      addon_version = aws_eks_addon.ebs_csi_driver.addon_version
      status        = aws_eks_addon.ebs_csi_driver.status
    }
  }
}

# Service Account IAM Roles
output "aws_load_balancer_controller_role_arn" {
  description = "ARN of the AWS Load Balancer Controller IAM role"
  value       = aws_iam_role.aws_load_balancer_controller.arn
}

output "ebs_csi_driver_role_arn" {
  description = "ARN of the EBS CSI driver IAM role"
  value       = aws_iam_role.ebs_csi_driver.arn
}

output "cluster_autoscaler_role_arn" {
  description = "ARN of the cluster autoscaler IAM role"
  value       = aws_iam_role.cluster_autoscaler.arn
}

# Kubernetes Configuration
output "kubectl_config" {
  description = "kubectl config as generated by the module"
  value = {
    cluster_name                   = aws_eks_cluster.main.name
    endpoint                       = aws_eks_cluster.main.endpoint
    certificate_authority_data     = aws_eks_cluster.main.certificate_authority[0].data
    region                         = var.region
    aws_authenticator_command      = "aws"
    aws_authenticator_command_args = ["eks", "get-token", "--cluster-name", aws_eks_cluster.main.name]
  }
}

# Kubeconfig for local development
output "kubeconfig_filename" {
  description = "The filename of the generated kubeconfig"
  value       = "kubeconfig_${aws_eks_cluster.main.name}"
}

# AWS CLI commands
output "aws_cli_commands" {
  description = "AWS CLI commands to configure kubectl"
  value = {
    update_kubeconfig = "aws eks update-kubeconfig --region ${var.region} --name ${aws_eks_cluster.main.name}"
    describe_cluster  = "aws eks describe-cluster --region ${var.region} --name ${aws_eks_cluster.main.name}"
    list_nodegroups   = "aws eks list-nodegroups --region ${var.region} --cluster-name ${aws_eks_cluster.main.name}"
  }
}

# Cluster Tags
output "cluster_tags" {
  description = "A map of tags assigned to the cluster"
  value       = aws_eks_cluster.main.tags_all
}

# Node Group Tags
output "node_group_tags" {
  description = "A map of tags assigned to the node group"
  value       = aws_eks_node_group.main.tags_all
}

# Networking Information
output "vpc_id" {
  description = "The VPC ID where the cluster is deployed"
  value       = var.vpc_id
}

output "subnet_ids" {
  description = "List of subnet IDs used by the cluster"
  value       = var.subnet_ids
}

output "private_subnet_ids" {
  description = "List of private subnet IDs used by the node groups"
  value       = var.private_subnet_ids
}

# Cluster Configuration
output "cluster_endpoint_config" {
  description = "Cluster endpoint configuration"
  value = {
    private_access = var.endpoint_private_access
    public_access  = var.endpoint_public_access
    public_cidrs   = var.public_access_cidrs
  }
}

# Scaling Configuration
output "node_group_scaling_config" {
  description = "Node group scaling configuration"
  value = {
    desired_size = aws_eks_node_group.main.scaling_config[0].desired_size
    max_size     = aws_eks_node_group.main.scaling_config[0].max_size
    min_size     = aws_eks_node_group.main.scaling_config[0].min_size
  }
}

# Cluster Info for Monitoring
output "cluster_info" {
  description = "Comprehensive cluster information for monitoring and management"
  value = {
    cluster_name     = aws_eks_cluster.main.name
    cluster_endpoint = aws_eks_cluster.main.endpoint
    cluster_version  = aws_eks_cluster.main.version
    region           = var.region
    oidc_issuer_url  = aws_eks_cluster.main.identity[0].oidc[0].issuer
    vpc_id           = var.vpc_id
    environment      = var.environment
    created_at       = aws_eks_cluster.main.created_at
  }
} 