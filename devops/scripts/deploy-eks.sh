#!/bin/bash

# EKS Deployment Script
# This script deploys EKS infrastructure using Terraform

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
TERRAFORM_DIR="$PROJECT_ROOT/terraform"

# Default values
ENVIRONMENT="staging"
ACTION="plan"
AUTO_APPROVE=false
DESTROY=false

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_debug() {
    echo -e "${BLUE}[DEBUG]${NC} $1"
}

show_usage() {
    cat << EOF
Usage: $0 [OPTIONS]

Deploy EKS infrastructure using Terraform

OPTIONS:
    -e, --environment ENVIRONMENT   Environment to deploy (staging|production) [default: staging]
    -a, --action ACTION            Terraform action (plan|apply|destroy) [default: plan]
    -y, --auto-approve             Auto approve Terraform apply/destroy
    -d, --destroy                  Destroy infrastructure
    -h, --help                     Show this help message

EXAMPLES:
    $0 -e staging -a plan                    # Plan staging deployment
    $0 -e staging -a apply                   # Apply staging deployment
    $0 -e staging -a apply -y                # Apply staging with auto-approve
    $0 -e production -a apply                # Apply production deployment
    $0 -e staging -d -y                      # Destroy staging infrastructure

PREREQUISITES:
    - AWS CLI configured with appropriate credentials
    - Terraform >= 1.0 installed
    - kubectl installed
    - eksctl installed (optional, for manual operations)

EOF
}

check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check AWS CLI
    if ! command -v aws &> /dev/null; then
        log_error "AWS CLI is not installed"
        return 1
    fi
    
    # Check AWS credentials
    if ! aws sts get-caller-identity &> /dev/null; then
        log_error "AWS credentials not configured or invalid"
        return 1
    fi
    
    # Check Terraform
    if ! command -v terraform &> /dev/null; then
        log_error "Terraform is not installed"
        return 1
    fi
    
    # Check kubectl
    if ! command -v kubectl &> /dev/null; then
        log_warn "kubectl is not installed (recommended for cluster access)"
    fi
    
    # Check eksctl
    if ! command -v eksctl &> /dev/null; then
        log_warn "eksctl is not installed (optional)"
    fi
    
    log_info "Prerequisites check completed"
}

validate_environment() {
    local env="$1"
    if [[ ! "$env" =~ ^(staging|production)$ ]]; then
        log_error "Invalid environment: $env. Must be 'staging' or 'production'"
        return 1
    fi
    
    local env_dir="$TERRAFORM_DIR/environments/$env"
    if [[ ! -d "$env_dir" ]]; then
        log_error "Environment directory not found: $env_dir"
        return 1
    fi
}

setup_terraform_backend() {
    local environment="$1"
    log_info "Setting up Terraform backend for $environment..."
    
    # Note: In a real scenario, you would create S3 bucket and DynamoDB table
    # This is just a placeholder for the setup
    cat << EOF
To setup Terraform backend, you need to:

1. Create S3 bucket for state storage:
   aws s3 mb s3://your-terraform-state-bucket-name-${environment}
   
2. Create DynamoDB table for state locking:
   aws dynamodb create-table \\
     --table-name terraform-state-lock-${environment} \\
     --attribute-definitions AttributeName=LockID,AttributeType=S \\
     --key-schema AttributeName=LockID,KeyType=HASH \\
     --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5

3. Update backend configuration in main.tf

EOF
}

terraform_init() {
    local env_dir="$1"
    log_info "Initializing Terraform in $env_dir..."
    
    cd "$env_dir"
    terraform init
}

terraform_plan() {
    local env_dir="$1"
    log_info "Running Terraform plan in $env_dir..."
    
    cd "$env_dir"
    terraform plan -var-file="terraform.tfvars"
}

terraform_apply() {
    local env_dir="$1"
    local auto_approve="$2"
    
    log_info "Running Terraform apply in $env_dir..."
    
    cd "$env_dir"
    if [[ "$auto_approve" == "true" ]]; then
        terraform apply -var-file="terraform.tfvars" -auto-approve
    else
        terraform apply -var-file="terraform.tfvars"
    fi
}

terraform_destroy() {
    local env_dir="$1"
    local auto_approve="$2"
    
    log_warn "Running Terraform destroy in $env_dir..."
    
    cd "$env_dir"
    if [[ "$auto_approve" == "true" ]]; then
        terraform destroy -var-file="terraform.tfvars" -auto-approve
    else
        terraform destroy -var-file="terraform.tfvars"
    fi
}

configure_kubectl() {
    local environment="$1"
    local env_dir="$TERRAFORM_DIR/environments/$environment"
    
    log_info "Configuring kubectl for EKS cluster..."
    
    cd "$env_dir"
    
    # Get cluster name from terraform output
    local cluster_name
    cluster_name=$(terraform output -raw cluster_info | jq -r '.cluster_name' 2>/dev/null || echo "")
    
    if [[ -z "$cluster_name" ]]; then
        log_warn "Could not get cluster name from Terraform output"
        log_info "Run: aws eks update-kubeconfig --region <region> --name <cluster-name>"
        return 1
    fi
    
    # Get region from terraform output
    local region
    region=$(terraform output -raw cluster_info | jq -r '.region' 2>/dev/null || echo "us-west-2")
    
    # Update kubeconfig
    aws eks update-kubeconfig --region "$region" --name "$cluster_name"
    
    log_info "kubectl configured for cluster: $cluster_name"
    
    # Test connection
    if kubectl get nodes &> /dev/null; then
        log_info "Successfully connected to EKS cluster"
        kubectl get nodes
    else
        log_warn "Could not connect to EKS cluster. Check your configuration."
    fi
}

install_cluster_components() {
    local environment="$1"
    
    log_info "Installing additional cluster components..."
    
    # Check if kubectl is configured
    if ! kubectl get nodes &> /dev/null; then
        log_error "kubectl is not configured. Configure it first."
        return 1
    fi
    
    # Install AWS Load Balancer Controller (if not already installed via Terraform)
    log_info "Checking AWS Load Balancer Controller..."
    if ! kubectl get deployment -n kube-system aws-load-balancer-controller &> /dev/null; then
        log_info "Installing AWS Load Balancer Controller..."
        # This would be handled by Terraform in our setup
        log_info "AWS Load Balancer Controller should be installed via Terraform"
    fi
    
    # Install Cluster Autoscaler (if not already installed via Terraform)
    log_info "Checking Cluster Autoscaler..."
    if ! kubectl get deployment -n kube-system cluster-autoscaler &> /dev/null; then
        log_info "Installing Cluster Autoscaler..."
        # This would be handled by Terraform in our setup
        log_info "Cluster Autoscaler should be installed via Terraform"
    fi
    
    # Install metrics server
    log_info "Checking Metrics Server..."
    if ! kubectl get deployment -n kube-system metrics-server &> /dev/null; then
        log_info "Installing Metrics Server..."
        kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
    fi
}

show_cluster_info() {
    local environment="$1"
    local env_dir="$TERRAFORM_DIR/environments/$environment"
    
    log_info "Cluster Information:"
    
    cd "$env_dir"
    
    # Show terraform outputs
    if terraform output cluster_info &> /dev/null; then
        echo -e "${BLUE}Terraform Outputs:${NC}"
        terraform output cluster_info
        echo
    fi
    
    # Show kubectl cluster info
    if kubectl cluster-info &> /dev/null; then
        echo -e "${BLUE}Kubectl Cluster Info:${NC}"
        kubectl cluster-info
        echo
        
        echo -e "${BLUE}Nodes:${NC}"
        kubectl get nodes -o wide
        echo
        
        echo -e "${BLUE}Pods:${NC}"
        kubectl get pods --all-namespaces
    fi
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -e|--environment)
            ENVIRONMENT="$2"
            shift 2
            ;;
        -a|--action)
            ACTION="$2"
            shift 2
            ;;
        -y|--auto-approve)
            AUTO_APPROVE=true
            shift
            ;;
        -d|--destroy)
            DESTROY=true
            ACTION="destroy"
            shift
            ;;
        -h|--help)
            show_usage
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Main execution
main() {
    log_info "Starting EKS deployment script..."
    log_info "Environment: $ENVIRONMENT"
    log_info "Action: $ACTION"
    
    # Check prerequisites
    check_prerequisites
    
    # Validate environment
    validate_environment "$ENVIRONMENT"
    
    local env_dir="$TERRAFORM_DIR/environments/$ENVIRONMENT"
    
    # Check if terraform.tfvars exists
    if [[ ! -f "$env_dir/terraform.tfvars" ]]; then
        log_error "terraform.tfvars not found in $env_dir"
        log_info "Copy terraform.tfvars.example to terraform.tfvars and customize it"
        exit 1
    fi
    
    # Initialize Terraform
    terraform_init "$env_dir"
    
    # Execute based on action
    case "$ACTION" in
        plan)
            terraform_plan "$env_dir"
            ;;
        apply)
            terraform_plan "$env_dir"
            terraform_apply "$env_dir" "$AUTO_APPROVE"
            
            # Configure kubectl after successful apply
            log_info "Waiting for cluster to be ready..."
            sleep 60  # Wait for cluster to be ready
            
            configure_kubectl "$ENVIRONMENT"
            install_cluster_components "$ENVIRONMENT"
            show_cluster_info "$ENVIRONMENT"
            ;;
        destroy)
            if [[ "$DESTROY" == "true" ]]; then
                terraform_destroy "$env_dir" "$AUTO_APPROVE"
            else
                log_error "Use -d flag to confirm destroy action"
                exit 1
            fi
            ;;
        *)
            log_error "Invalid action: $ACTION"
            show_usage
            exit 1
            ;;
    esac
    
    log_info "EKS deployment script completed successfully!"
}

# Run main function
main "$@" 