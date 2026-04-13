#!/bin/bash

echo "🛑 Stopping Axiom Stack (redis-go)..."

# 1. Kill background port-forwarding processes
echo "🔌 Closing network tunnels..."
# This finds any 'kubectl port-forward' process and kills it
pkill -f "kubectl port-forward" || echo "No active tunnels found."

# 2. Delete Kubernetes resources
echo "☸️  Removing Kubernetes manifests..."
# Using the files ensures we delete exactly what we created
kubectl delete -f src/k8s/service.yaml --ignore-not-found
kubectl delete -f src/k8s/deployment.yaml --ignore-not-found

# 3. Handle Persistent Data (Optional/Mindful)
# StatefulSets do NOT delete PVCs automatically to prevent data loss.
# If you want a TRULY clean state, uncomment the lines below:
# echo "💾 Wiping persistent volumes..."
# kubectl delete pvc -l app=redis-go

# 4. Cleanup Local Docker (Compose)
if [ -f "docker-compose.yml" ]; then
    echo "🐳 Stopping Docker Compose..."
    docker-compose down
fi

echo "✅ All services stopped and cleaned up."