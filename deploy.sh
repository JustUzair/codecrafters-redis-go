#!/bin/bash

# --- CONFIGURATION ---
DOCKER_USER="0xjustuzair"
BACKEND_IMAGE="redis-go"
TAG="latest"

# Exit immediately if a command fails
set -e

echo "🚀 Starting Deployment for Redis-go Stack..."

# 1. CLEANUP: Clear existing "broken" pods and deployments
echo "🧹 Clearing old deployments..."
kubectl delete deployments --all --ignore-not-found
# Give K8s a second to breathe
sleep 2

# ... (Configuration above) ...

echo "🚀 Starting Deployment for Redis-go Stack..."

# Check if K8s is reachable
if ! kubectl cluster-info > /dev/null 2>&1; then
    echo "🔄 Cluster unreachable. Attempting to repair context..."
    minikube update-context
    
    # If it's still dead, then do the heavy lift
    if ! kubectl cluster-info > /dev/null 2>&1; then
        echo "⚠️  Context repair failed. Restarting Minikube..."
        minikube start
    fi
fi

# 2. ENVIRONMENT: Create Secrets and ConfigMaps
# We delete them first to ensure we aren't using old values
echo "🔑 Updating Kubernetes Secrets and Configs..."

# Backend Env
# kubectl delete secret backend-secrets --ignore-not-found
# kubectl create secret generic backend-secrets --from-env-file=./agent/.env

# kubectl delete configmap backend-config --ignore-not-found
# kubectl create configmap backend-config --from-literal=NODE_ENV=production


# 3. BUILD: Build the images
echo "📦 Building Docker images..."
docker build -t $DOCKER_USER/$BACKEND_IMAGE:$TAG ./


# 4. PUSH: Push to Docker Hub
echo "📤 Pushing images to Docker Hub..."
docker push $DOCKER_USER/$BACKEND_IMAGE:$TAG

# 5. DEPLOY: Apply Kubernetes Manifests
echo "☸️  Applying Kubernetes manifests..."
kubectl apply -f app/k8s/deployment.yaml
kubectl apply -f app/k8s/service.yaml

echo "✅ Deployment complete!"
echo "📍 Watch your pods turn green: kubectl get pods -w"

# Wait for pods to actually be ready before forwarding
echo "⏳ Waiting for pods to be ready..."
kubectl wait --for=condition=ready pod -l app=redis-go --timeout=90s

./serve.sh