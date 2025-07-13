#!/bin/bash

# Microservices Deployment Script
# Usage: ./scripts/deploy.sh [environment] [action]
# Environment: staging | production
# Action: deploy | update | rollback | status

set -e

# Configuration
NAMESPACE=""
ENVIRONMENT=""
ACTION=""
REGISTRY="ghcr.io"
IMAGE_PREFIX="your-org/microservices"
VERSION="latest"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Show usage
usage() {
    echo "Usage: $0 [environment] [action]"
    echo "Environment: staging | production"
    echo "Action: deploy | update | rollback | status | destroy"
    echo ""
    echo "Examples:"
    echo "  $0 staging deploy    # Deploy to staging"
    echo "  $0 production update # Update production"
    echo "  $0 staging status    # Check staging status"
    exit 1
}

# Check prerequisites
check_prerequisites() {
    print_info "Checking prerequisites..."
    
    # Check kubectl
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl not found. Please install kubectl."
        exit 1
    fi
    
    # Check helm
    if ! command -v helm &> /dev/null; then
        print_error "helm not found. Please install helm."
        exit 1
    fi
    
    # Check cluster connection
    if ! kubectl cluster-info &> /dev/null; then
        print_error "Cannot connect to Kubernetes cluster."
        exit 1
    fi
    
    print_info "Prerequisites check passed."
}

# Set environment variables
set_environment() {
    case $ENVIRONMENT in
        staging)
            NAMESPACE="microservices-staging"
            ;;
        production)
            NAMESPACE="microservices"
            ;;
        *)
            print_error "Invalid environment: $ENVIRONMENT"
            usage
            ;;
    esac
    
    print_info "Environment: $ENVIRONMENT"
    print_info "Namespace: $NAMESPACE"
}

# Deploy infrastructure components
deploy_infrastructure() {
    print_info "Deploying infrastructure components..."
    
    # Create namespace
    kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -
    
    # Deploy ConfigMap
    kubectl apply -f k8s/configmap.yaml -n $NAMESPACE
    
    # Deploy cert-manager ClusterIssuer (only for production)
    if [ "$ENVIRONMENT" = "production" ]; then
        kubectl apply -f k8s/cert-manager.yaml
    fi
    
    print_info "Infrastructure components deployed."
}

# Deploy monitoring stack
deploy_monitoring() {
    print_info "Deploying monitoring stack..."
    
    # Deploy Prometheus
    kubectl apply -f k8s/monitoring/prometheus.yaml
    
    # Deploy Grafana
    kubectl apply -f k8s/monitoring/grafana.yaml
    
    # Deploy Fluentd
    kubectl apply -f k8s/monitoring/fluentd.yaml
    
    print_info "Monitoring stack deployed."
}

# Deploy microservices
deploy_microservices() {
    print_info "Deploying microservices..."
    
    # Services to deploy
    services=("user-service" "product-service" "order-service" "api-gateway")
    
    for service in "${services[@]}"; do
        print_info "Deploying $service..."
        
        # Update image tags in manifests
        sed -i "s|image: microservices/${service}:latest|image: ${REGISTRY}/${IMAGE_PREFIX}/${service}:${VERSION}|g" k8s/${service}.yaml
        
        # Apply manifest
        kubectl apply -f k8s/${service}.yaml -n $NAMESPACE
        
        # Wait for deployment
        kubectl rollout status deployment/${service} -n $NAMESPACE --timeout=300s
        
        print_info "$service deployed successfully."
    done
    
    print_info "All microservices deployed."
}

# Deploy ingress and networking
deploy_networking() {
    print_info "Deploying networking components..."
    
    # Deploy network policies
    kubectl apply -f k8s/network-policy.yaml -n $NAMESPACE
    
    # Deploy ingress
    kubectl apply -f k8s/ingress.yaml -n $NAMESPACE
    
    # Deploy HPA
    kubectl apply -f k8s/hpa.yaml -n $NAMESPACE
    
    print_info "Networking components deployed."
}

# Full deployment
deploy() {
    print_info "Starting deployment to $ENVIRONMENT..."
    
    deploy_infrastructure
    
    if [ "$ENVIRONMENT" = "production" ]; then
        deploy_monitoring
    fi
    
    deploy_microservices
    deploy_networking
    
    print_info "Deployment completed successfully!"
    
    # Show status
    show_status
}

# Update deployment
update() {
    print_info "Updating deployment in $ENVIRONMENT..."
    
    # Update microservices
    deploy_microservices
    
    print_info "Update completed successfully!"
    
    # Show status
    show_status
}

# Rollback deployment
rollback() {
    print_info "Rolling back deployment in $ENVIRONMENT..."
    
    services=("user-service" "product-service" "order-service" "api-gateway")
    
    for service in "${services[@]}"; do
        print_info "Rolling back $service..."
        kubectl rollout undo deployment/${service} -n $NAMESPACE
        kubectl rollout status deployment/${service} -n $NAMESPACE --timeout=300s
        print_info "$service rolled back successfully."
    done
    
    print_info "Rollback completed successfully!"
    
    # Show status
    show_status
}

# Show deployment status
show_status() {
    print_info "Checking deployment status..."
    
    echo ""
    echo "=== Namespace: $NAMESPACE ==="
    kubectl get namespace $NAMESPACE 2>/dev/null || echo "Namespace not found"
    
    echo ""
    echo "=== Deployments ==="
    kubectl get deployments -n $NAMESPACE
    
    echo ""
    echo "=== Services ==="
    kubectl get services -n $NAMESPACE
    
    echo ""
    echo "=== Pods ==="
    kubectl get pods -n $NAMESPACE
    
    echo ""
    echo "=== Ingress ==="
    kubectl get ingress -n $NAMESPACE
    
    echo ""
    echo "=== HPA ==="
    kubectl get hpa -n $NAMESPACE
    
    echo ""
    echo "=== Recent Events ==="
    kubectl get events -n $NAMESPACE --sort-by=.metadata.creationTimestamp | tail -10
}

# Destroy deployment
destroy() {
    print_warning "This will destroy all resources in $NAMESPACE."
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_info "Destroying deployment..."
        
        # Delete all resources in namespace
        kubectl delete namespace $NAMESPACE --ignore-not-found=true
        
        # Delete cluster resources (only for production)
        if [ "$ENVIRONMENT" = "production" ]; then
            kubectl delete clusterissuer letsencrypt-prod --ignore-not-found=true
            kubectl delete clusterissuer letsencrypt-staging --ignore-not-found=true
        fi
        
        print_info "Deployment destroyed."
    else
        print_info "Destroy cancelled."
    fi
}

# Health check
health_check() {
    print_info "Performing health check..."
    
    # Check if all pods are running
    kubectl get pods -n $NAMESPACE --field-selector=status.phase!=Running --no-headers | wc -l | xargs -I {} [ {} -eq 0 ] && {
        print_info "All pods are running."
    } || {
        print_warning "Some pods are not running:"
        kubectl get pods -n $NAMESPACE --field-selector=status.phase!=Running
    }
    
    # Check ingress
    INGRESS_IP=$(kubectl get ingress -n $NAMESPACE -o jsonpath='{.items[0].status.loadBalancer.ingress[0].ip}' 2>/dev/null)
    if [ -n "$INGRESS_IP" ]; then
        print_info "Ingress IP: $INGRESS_IP"
    else
        print_warning "Ingress IP not ready yet."
    fi
}

# Main function
main() {
    if [ $# -lt 2 ]; then
        usage
    fi
    
    ENVIRONMENT=$1
    ACTION=$2
    
    if [ $# -ge 3 ]; then
        VERSION=$3
    fi
    
    check_prerequisites
    set_environment
    
    case $ACTION in
        deploy)
            deploy
            ;;
        update)
            update
            ;;
        rollback)
            rollback
            ;;
        status)
            show_status
            ;;
        destroy)
            destroy
            ;;
        health)
            health_check
            ;;
        *)
            print_error "Invalid action: $ACTION"
            usage
            ;;
    esac
}

# Run main function
main "$@" 