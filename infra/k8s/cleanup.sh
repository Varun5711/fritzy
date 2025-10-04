#!/bin/bash

set -e

NAMESPACE="fritzy"
MONITORING_NAMESPACE="monitoring"

echo "Removing Fritzy Backend from Kubernetes..."

echo "Deleting services..."
kubectl delete -f services/ --ignore-not-found=true

echo "Deleting databases..."
kubectl delete -f databases/ --ignore-not-found=true

echo "Deleting ingress..."
kubectl delete -f ingress/ --ignore-not-found=true

echo "Deleting monitoring stack..."
kubectl delete -f monitoring/ --ignore-not-found=true

echo "Deleting secrets and configmaps..."
kubectl delete -f base/ --ignore-not-found=true

echo "Deleting namespaces..."
kubectl delete namespace $NAMESPACE --ignore-not-found=true
kubectl delete namespace $MONITORING_NAMESPACE --ignore-not-found=true

echo "Cleanup complete!"
