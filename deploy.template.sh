#!/bin/bash
set -e

PROJECT_ID=""  # Set your actual project ID here
REGION="" # Set your default region
REPO="" # Set the artifact registry repository name

# Check if PROJECT_ID is set
if [[ -z "$PROJECT_ID" || "$PROJECT_ID" == "YOUR_PROJECT_ID_HERE" ]]; then
    echo "Error: PROJECT_ID not set. Edit this file and set your GCP project ID."
    exit 1
fi

# Make sure the tags match your own tags
FULL_IMAGE_BACKEND="$REGION-docker.pkg.dev/$PROJECT_ID/$REPO/gke-backend:latest"
FULL_IMAGE_FRONTEND="$REGION-docker.pkg.dev/$PROJECT_ID/$REPO/gke-frontend:latest"
FULL_IMAGE_SESSION="$REGION-docker.pkg.dev/$PROJECT_ID/$REPO/gke-session:latest"
FULL_IMAGE_MIGRATE="$REGION-docker.pkg.dev/$PROJECT_ID/$REPO/gke-migrate:latest"

echo "Deploying images with tag: latest"

kubectl apply -f k8s/app-config.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/postgres-pvc.yaml
kubectl apply -f k8s/postgres-deploy.yaml
kubectl apply -f k8s/postgres-service.yaml

export PUBLIC_GO_API_BASE="/api"
export CORS_ORIGINS="*"
export BACKEND_IMAGE_URL="$FULL_IMAGE_BACKEND"
export FRONTEND_IMAGE_URL="$FULL_IMAGE_FRONTEND" 
export SESSION_IMAGE_URL="$FULL_IMAGE_SESSION"
export MIGRATE_IMAGE_URL="$FULL_IMAGE_MIGRATE"

envsubst < k8s/backend-deploy.yaml | kubectl apply -f -
kubectl apply -f k8s/backend-service.yaml
envsubst < k8s/frontend-deploy.yaml | kubectl apply -f -
kubectl apply -f k8s/frontend-service.yaml
envsubst < k8s/session-deploy.yaml | kubectl apply -f -
kubectl apply -f k8s/session-service.yaml
envsubst < k8s/migrate-job.yaml | kubectl apply -f -

# Apply Gateway API resources (no env substitution needed)
kubectl apply -f k8s/gateway.yaml
kubectl apply -f k8s/frontend-route.yaml
kubectl apply -f k8s/backend-route.yaml
kubectl apply -f k8s/health-check.yaml

echo "All deployments and services applied successfully!"

# Wait for Gateway to get an IP address
echo "Waiting for Gateway to get an external IP..."
for i in {1..30}; do
    GATEWAY_IP=$(kubectl get gateway external-gateway -o jsonpath='{.status.addresses[0].value}' 2>/dev/null || echo "")
    if [[ -n "$GATEWAY_IP" ]]; then
        echo "Gateway IP: $GATEWAY_IP"
        echo "Your application will be available at: http://$GATEWAY_IP"
        echo "API endpoint: http://$GATEWAY_IP/api"
        break
    fi
    echo "Waiting for Gateway IP... (attempt $i/30)"
    sleep 10
done

if [[ -z "$GATEWAY_IP" ]]; then
    echo "Warning: Gateway IP not ready yet. Check status with:"
    echo "kubectl get gateway external-gateway"
fi