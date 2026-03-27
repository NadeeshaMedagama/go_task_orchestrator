#!/bin/bash

# Go Task Orchestrator - Minikube Health Check Script
# This script verifies that all services are running correctly

set -e

echo "🔍 Go Task Orchestrator - Health Check Script"
echo "=============================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

NAMESPACE="go-task-orchestrator"

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
    fi
}

echo "📦 Checking Pods Status..."
echo "=========================="
kubectl get pods -n $NAMESPACE

echo ""
echo "🔗 Checking Services..."
echo "======================"
kubectl get svc -n $NAMESPACE

echo ""
echo "🌐 Checking Ingress..."
echo "====================="
kubectl get ingress -n $NAMESPACE

echo ""
echo "📋 Checking Service Details..."
echo "=============================="
MINIKUBE_IP=$(minikube ip)
echo "Minikube IP: $MINIKUBE_IP"
echo ""

# Check if services are accessible
echo "🧪 Testing Service Connectivity..."
echo "==================================="

# Get pod names
API_GATEWAY_POD=$(kubectl get pods -n $NAMESPACE -l app=api-gateway -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
AUTH_SERVICE_POD=$(kubectl get pods -n $NAMESPACE -l app=auth-service -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")
TASK_SERVICE_POD=$(kubectl get pods -n $NAMESPACE -l app=task-service -o jsonpath='{.items[0].metadata.name}' 2>/dev/null || echo "")

if [ -z "$API_GATEWAY_POD" ]; then
    echo -e "${YELLOW}⚠ API Gateway pod not found. Checking alternate labels...${NC}"
    API_GATEWAY_POD=$(kubectl get pods -n $NAMESPACE | grep api-gateway | head -1 | awk '{print $1}')
fi

if [ -z "$AUTH_SERVICE_POD" ]; then
    AUTH_SERVICE_POD=$(kubectl get pods -n $NAMESPACE | grep auth-service | head -1 | awk '{print $1}')
fi

if [ -z "$TASK_SERVICE_POD" ]; then
    TASK_SERVICE_POD=$(kubectl get pods -n $NAMESPACE | grep task-service | head -1 | awk '{print $1}')
fi

echo ""
echo "Testing API Gateway Health..."
if kubectl exec -n $NAMESPACE $API_GATEWAY_POD -- curl -s http://localhost:8080/health > /dev/null 2>&1; then
    print_status 0 "API Gateway is responding"
else
    print_status 1 "API Gateway is not responding"
fi

echo ""
echo "📋 How to Access Your Services:"
echo "=============================="
echo ""
echo "Option 1: Using Ingress (Recommended)"
echo "-------------------------------------"
echo "1. Add to /etc/hosts:"
echo "   echo '${MINIKUBE_IP} go-task-orchestrator.local' | sudo tee -a /etc/hosts"
echo ""
echo "2. Base URL for Postman:"
echo "   http://go-task-orchestrator.local/api"
echo ""
echo "Option 2: Using Port Forwarding"
echo "-------------------------------"
echo "1. Run in terminal:"
echo "   kubectl port-forward -n $NAMESPACE svc/api-gateway 8080:8080"
echo ""
echo "2. Base URL for Postman:"
echo "   http://localhost:8080/api"
echo ""

echo "✅ Health check complete!"
echo ""
echo "📚 For detailed instructions, see: MINIKUBE_TESTING_GUIDE.md"

