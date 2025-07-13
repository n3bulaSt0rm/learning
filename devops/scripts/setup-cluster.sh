#!/bin/bash

# Kubernetes Cluster Setup Script
# This script installs and configures required components for the microservices

set -e

# Configuration
INGRESS_NAMESPACE="ingress-nginx"
CERT_MANAGER_NAMESPACE="cert-manager"
MONITORING_NAMESPACE="monitoring"

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

# Check if running as root
check_root() {
    if [ "$EUID" -eq 0 ]; then
        print_error "Please don't run this script as root"
        exit 1
    fi
}

# Check prerequisites
check_prerequisites() {
    print_info "Checking prerequisites..."
    
    # Check kubectl
    if ! command -v kubectl &> /dev/null; then
        print_error "kubectl not found. Please install kubectl first."
        exit 1
    fi
    
    # Check helm
    if ! command -v helm &> /dev/null; then
        print_error "helm not found. Installing helm..."
        install_helm
    fi
    
    # Check cluster connection
    if ! kubectl cluster-info &> /dev/null; then
        print_error "Cannot connect to Kubernetes cluster."
        exit 1
    fi
    
    print_info "Prerequisites check passed."
}

# Install Helm
install_helm() {
    print_info "Installing Helm..."
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash
    print_info "Helm installed successfully."
}

# Add Helm repositories
add_helm_repos() {
    print_info "Adding Helm repositories..."
    
    helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
    helm repo add jetstack https://charts.jetstack.io
    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    helm repo add grafana https://grafana.github.io/helm-charts
    helm repo add elastic https://helm.elastic.co
    
    helm repo update
    
    print_info "Helm repositories added successfully."
}

# Install NGINX Ingress Controller
install_ingress() {
    print_info "Installing NGINX Ingress Controller..."
    
    # Create namespace
    kubectl create namespace $INGRESS_NAMESPACE --dry-run=client -o yaml | kubectl apply -f -
    
    # Install NGINX Ingress Controller
    helm upgrade --install ingress-nginx ingress-nginx/ingress-nginx \
        --namespace $INGRESS_NAMESPACE \
        --set controller.replicaCount=2 \
        --set controller.metrics.enabled=true \
        --set controller.metrics.serviceMonitor.enabled=true \
        --set controller.podAnnotations."prometheus\.io/scrape"="true" \
        --set controller.podAnnotations."prometheus\.io/port"="10254" \
        --set controller.service.type=LoadBalancer \
        --wait --timeout=300s
    
    print_info "NGINX Ingress Controller installed successfully."
}

# Install cert-manager
install_cert_manager() {
    print_info "Installing cert-manager..."
    
    # Create namespace
    kubectl create namespace $CERT_MANAGER_NAMESPACE --dry-run=client -o yaml | kubectl apply -f -
    
    # Install cert-manager
    helm upgrade --install cert-manager jetstack/cert-manager \
        --namespace $CERT_MANAGER_NAMESPACE \
        --version v1.13.0 \
        --set installCRDs=true \
        --set prometheus.enabled=true \
        --set prometheus.servicemonitor.enabled=true \
        --wait --timeout=300s
    
    print_info "cert-manager installed successfully."
}

# Install Metrics Server
install_metrics_server() {
    print_info "Installing Metrics Server..."
    
    # Check if metrics server is already installed
    if kubectl get deployment metrics-server -n kube-system &> /dev/null; then
        print_info "Metrics Server already installed."
        return
    fi
    
    # Install metrics server
    kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
    
    # Wait for metrics server to be ready
    kubectl rollout status deployment/metrics-server -n kube-system --timeout=300s
    
    print_info "Metrics Server installed successfully."
}

# Install Prometheus stack
install_prometheus() {
    print_info "Installing Prometheus monitoring stack..."
    
    # Create namespace
    kubectl create namespace $MONITORING_NAMESPACE --dry-run=client -o yaml | kubectl apply -f -
    
    # Install Prometheus stack
    helm upgrade --install prometheus prometheus-community/kube-prometheus-stack \
        --namespace $MONITORING_NAMESPACE \
        --set prometheus.prometheusSpec.retention=30d \
        --set prometheus.prometheusSpec.storageSpec.volumeClaimTemplate.spec.resources.requests.storage=50Gi \
        --set grafana.adminPassword=admin123 \
        --set grafana.service.type=LoadBalancer \
        --set grafana.persistence.enabled=true \
        --set grafana.persistence.size=10Gi \
        --set alertmanager.alertmanagerSpec.storage.volumeClaimTemplate.spec.resources.requests.storage=10Gi \
        --wait --timeout=600s
    
    print_info "Prometheus monitoring stack installed successfully."
}

# Install Elasticsearch and Kibana
install_elasticsearch() {
    print_info "Installing Elasticsearch and Kibana..."
    
    # Install Elasticsearch
    helm upgrade --install elasticsearch elastic/elasticsearch \
        --namespace $MONITORING_NAMESPACE \
        --set replicas=3 \
        --set volumeClaimTemplate.resources.requests.storage=10Gi \
        --set esConfig."elasticsearch\.yml"="cluster.name: \"docker-cluster\"\nnetwork.host: 0.0.0.0\n" \
        --wait --timeout=600s
    
    # Install Kibana
    helm upgrade --install kibana elastic/kibana \
        --namespace $MONITORING_NAMESPACE \
        --set service.type=LoadBalancer \
        --set elasticsearchHosts="http://elasticsearch-master:9200" \
        --wait --timeout=300s
    
    print_info "Elasticsearch and Kibana installed successfully."
}

# Setup storage class (if needed)
setup_storage_class() {
    print_info "Setting up storage class..."
    
    # Check if default storage class exists
    if kubectl get storageclass -o jsonpath='{.items[?(@.metadata.annotations.storageclass\.kubernetes\.io/is-default-class=="true")].metadata.name}' | grep -q .; then
        print_info "Default storage class already exists."
    else
        print_warning "No default storage class found. You may need to configure one."
    fi
}

# Verify installation
verify_installation() {
    print_info "Verifying installation..."
    
    # Check NGINX Ingress
    kubectl get pods -n $INGRESS_NAMESPACE
    print_info "NGINX Ingress Controller status: $(kubectl get pods -n $INGRESS_NAMESPACE -l app.kubernetes.io/name=ingress-nginx -o jsonpath='{.items[0].status.phase}')"
    
    # Check cert-manager
    kubectl get pods -n $CERT_MANAGER_NAMESPACE
    print_info "cert-manager status: $(kubectl get pods -n $CERT_MANAGER_NAMESPACE -l app.kubernetes.io/name=cert-manager -o jsonpath='{.items[0].status.phase}')"
    
    # Check Metrics Server
    kubectl get pods -n kube-system -l k8s-app=metrics-server
    print_info "Metrics Server status: $(kubectl get pods -n kube-system -l k8s-app=metrics-server -o jsonpath='{.items[0].status.phase}')"
    
    # Check Prometheus
    kubectl get pods -n $MONITORING_NAMESPACE
    print_info "Prometheus stack installed in namespace: $MONITORING_NAMESPACE"
    
    # Get LoadBalancer IPs
    print_info "Getting LoadBalancer IPs..."
    INGRESS_IP=$(kubectl get svc ingress-nginx-controller -n $INGRESS_NAMESPACE -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
    GRAFANA_IP=$(kubectl get svc prometheus-grafana -n $MONITORING_NAMESPACE -o jsonpath='{.status.loadBalancer.ingress[0].ip}')
    
    if [ -n "$INGRESS_IP" ]; then
        print_info "NGINX Ingress IP: $INGRESS_IP"
    else
        print_warning "NGINX Ingress IP not ready yet. Check with: kubectl get svc ingress-nginx-controller -n $INGRESS_NAMESPACE"
    fi
    
    if [ -n "$GRAFANA_IP" ]; then
        print_info "Grafana IP: $GRAFANA_IP"
        print_info "Grafana login: admin/admin123"
    else
        print_warning "Grafana IP not ready yet. Check with: kubectl get svc prometheus-grafana -n $MONITORING_NAMESPACE"
    fi
    
    print_info "Installation verification completed."
}

# Print cluster info
print_cluster_info() {
    print_info "Cluster Information:"
    kubectl cluster-info
    
    print_info "Cluster Nodes:"
    kubectl get nodes -o wide
    
    print_info "Cluster Version:"
    kubectl version --short
}

# Main setup function
main() {
    print_info "Starting Kubernetes cluster setup..."
    
    check_root
    check_prerequisites
    add_helm_repos
    setup_storage_class
    
    # Install components
    install_ingress
    install_cert_manager
    install_metrics_server
    install_prometheus
    # install_elasticsearch  # Uncomment if you want ELK stack
    
    # Verify installation
    verify_installation
    print_cluster_info
    
    print_info "Cluster setup completed successfully!"
    print_info "You can now deploy your microservices using: ./scripts/deploy.sh"
}

# Show help
show_help() {
    echo "Kubernetes Cluster Setup Script"
    echo ""
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  --minimal      Install only essential components"
    echo "  --full         Install all components including ELK stack"
    echo ""
    echo "This script will install:"
    echo "  - NGINX Ingress Controller"
    echo "  - cert-manager"
    echo "  - Metrics Server"
    echo "  - Prometheus monitoring stack"
    echo "  - Elasticsearch and Kibana (with --full option)"
}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        --minimal)
            MINIMAL=true
            shift
            ;;
        --full)
            FULL=true
            shift
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Run main function
main 