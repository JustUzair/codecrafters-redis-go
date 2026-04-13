#!/bin/bash

echo "🔌 Opening tunnels to Redis-Go..."

# 1. Port-forward the Redis service in the background
# We use 'svc/redis-go' to match your service.yaml metadata name
echo "⚙️  Redis (Go): localhost:6379"
kubectl port-forward svc/redis-go 6379:6379 > /dev/null 2>&1 &
PID_REDIS=$!

# 2. If you have a frontend/other services, add them here
# echo "🌐 Frontend: localhost:3000"
# kubectl port-forward svc/frontend 3000:3000 > /dev/null 2>&1 &

# 3. Use AppleScript to open a dashboard or logs in a new window
echo "📊 Opening Minikube Dashboard and Logs..."
osascript -e 'tell app "Terminal" to do script "minikube dashboard"'
osascript -e 'tell app "Terminal" to do script "kubectl logs -l app=redis-go -f --all-containers --prefix"'

echo "✅ Tunnels active (Redis PID: $PID_REDIS). Press Ctrl+C to stop this script and cleanup."

# Optional: Wait for user to exit to kill the background port-forward
trap "kill $PID_REDIS" EXIT
wait