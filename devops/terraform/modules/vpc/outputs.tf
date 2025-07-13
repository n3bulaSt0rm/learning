# VPC
output "vpc_id" {
  description = "ID of the VPC"
  value       = aws_vpc.main.id
}

output "vpc_arn" {
  description = "ARN of the VPC"
  value       = aws_vpc.main.arn
}

output "vpc_cidr_block" {
  description = "CIDR block of the VPC"
  value       = aws_vpc.main.cidr_block
}

output "vpc_main_route_table_id" {
  description = "ID of the main route table associated with this VPC"
  value       = aws_vpc.main.main_route_table_id
}

output "vpc_default_network_acl_id" {
  description = "ID of the default network ACL"
  value       = aws_vpc.main.default_network_acl_id
}

output "vpc_default_security_group_id" {
  description = "ID of the security group created by default on VPC creation"
  value       = aws_vpc.main.default_security_group_id
}

# Internet Gateway
output "igw_id" {
  description = "ID of the Internet Gateway"
  value       = aws_internet_gateway.main.id
}

output "igw_arn" {
  description = "ARN of the Internet Gateway"
  value       = aws_internet_gateway.main.arn
}

# Public Subnets
output "public_subnet_ids" {
  description = "List of IDs of public subnets"
  value       = aws_subnet.public[*].id
}

output "public_subnet_arns" {
  description = "List of ARNs of public subnets"
  value       = aws_subnet.public[*].arn
}

output "public_subnet_cidrs" {
  description = "List of CIDR blocks of public subnets"
  value       = aws_subnet.public[*].cidr_block
}

output "public_subnet_availability_zones" {
  description = "List of availability zones of public subnets"
  value       = aws_subnet.public[*].availability_zone
}

# Private Subnets
output "private_subnet_ids" {
  description = "List of IDs of private subnets"
  value       = aws_subnet.private[*].id
}

output "private_subnet_arns" {
  description = "List of ARNs of private subnets"
  value       = aws_subnet.private[*].arn
}

output "private_subnet_cidrs" {
  description = "List of CIDR blocks of private subnets"
  value       = aws_subnet.private[*].cidr_block
}

output "private_subnet_availability_zones" {
  description = "List of availability zones of private subnets"
  value       = aws_subnet.private[*].availability_zone
}

# NAT Gateways
output "nat_gateway_ids" {
  description = "List of IDs of NAT Gateways"
  value       = aws_nat_gateway.main[*].id
}

output "nat_gateway_public_ips" {
  description = "List of public IPs of NAT Gateways"
  value       = aws_nat_gateway.main[*].public_ip
}

output "nat_gateway_private_ips" {
  description = "List of private IPs of NAT Gateways"
  value       = aws_nat_gateway.main[*].private_ip
}

# Elastic IPs
output "eip_ids" {
  description = "List of IDs of Elastic IPs"
  value       = aws_eip.nat[*].id
}

output "eip_public_ips" {
  description = "List of public IPs of Elastic IPs"
  value       = aws_eip.nat[*].public_ip
}

# Route Tables
output "public_route_table_id" {
  description = "ID of the public route table"
  value       = aws_route_table.public.id
}

output "private_route_table_ids" {
  description = "List of IDs of private route tables"
  value       = aws_route_table.private[*].id
}

# VPC Endpoints
output "s3_vpc_endpoint_id" {
  description = "ID of the S3 VPC endpoint"
  value       = var.enable_s3_endpoint ? aws_vpc_endpoint.s3[0].id : null
}

output "ec2_vpc_endpoint_id" {
  description = "ID of the EC2 VPC endpoint"
  value       = var.enable_ec2_endpoint ? aws_vpc_endpoint.ec2[0].id : null
}

output "ecr_api_vpc_endpoint_id" {
  description = "ID of the ECR API VPC endpoint"
  value       = var.enable_ecr_endpoint ? aws_vpc_endpoint.ecr_api[0].id : null
}

output "ecr_dkr_vpc_endpoint_id" {
  description = "ID of the ECR DKR VPC endpoint"
  value       = var.enable_ecr_endpoint ? aws_vpc_endpoint.ecr_dkr[0].id : null
}

output "vpc_endpoint_security_group_id" {
  description = "ID of the VPC endpoint security group"
  value       = var.enable_ec2_endpoint || var.enable_ecr_endpoint ? aws_security_group.vpc_endpoint[0].id : null
}

# Flow Logs
output "vpc_flow_log_id" {
  description = "ID of the VPC Flow Log"
  value       = var.enable_flow_logs ? aws_flow_log.vpc_flow_log[0].id : null
}

output "vpc_flow_log_cloudwatch_log_group_name" {
  description = "Name of the CloudWatch Log Group for VPC Flow Logs"
  value       = var.enable_flow_logs ? aws_cloudwatch_log_group.vpc_flow_log[0].name : null
}

# DHCP Options
output "dhcp_options_id" {
  description = "ID of the DHCP options set"
  value       = var.enable_dhcp_options ? aws_vpc_dhcp_options.main[0].id : null
}

# Network Configuration for EKS
output "eks_subnet_ids" {
  description = "List of all subnet IDs for EKS cluster"
  value       = concat(aws_subnet.public[*].id, aws_subnet.private[*].id)
}

output "eks_private_subnet_ids" {
  description = "List of private subnet IDs for EKS node groups"
  value       = aws_subnet.private[*].id
}

output "eks_public_subnet_ids" {
  description = "List of public subnet IDs for EKS load balancers"
  value       = aws_subnet.public[*].id
}

# Availability Zones
output "availability_zones" {
  description = "List of availability zones used"
  value       = data.aws_availability_zones.available.names
}

# CIDR Information
output "vpc_cidr" {
  description = "CIDR block of the VPC"
  value       = var.vpc_cidr
}

output "public_subnet_cidrs_output" {
  description = "List of CIDR blocks for public subnets"
  value       = var.public_subnet_cidrs
}

output "private_subnet_cidrs_output" {
  description = "List of CIDR blocks for private subnets"
  value       = var.private_subnet_cidrs
}

# Network Summary
output "network_summary" {
  description = "Summary of network configuration"
  value = {
    vpc_id                = aws_vpc.main.id
    vpc_cidr              = aws_vpc.main.cidr_block
    availability_zones    = data.aws_availability_zones.available.names
    public_subnet_ids     = aws_subnet.public[*].id
    private_subnet_ids    = aws_subnet.private[*].id
    nat_gateway_count     = length(aws_nat_gateway.main)
    vpc_endpoints_enabled = {
      s3  = var.enable_s3_endpoint
      ec2 = var.enable_ec2_endpoint
      ecr = var.enable_ecr_endpoint
    }
    flow_logs_enabled = var.enable_flow_logs
  }
}

# Tags
output "tags" {
  description = "Tags applied to resources"
  value       = var.tags
} 