#!/bin/bash

set -e

NAMESPACE="fritzy"
MONITORING_NAMESPACE="monitoring"

echo "Deploying Fritzy Backend to Kubernetes..."

echo "Creating namespaces..."
kubectl apply -f base/namespace.yaml
kubectl create namespace $MONITORING_NAMESPACE --dry-run=client -o yaml | kubectl apply -f -

echo "Applying security configurations..."
kubectl apply -f security/pod-security.yaml
kubectl apply -f security/network-policies.yaml

echo "Creating secrets and configmaps..."
kubectl apply -f base/secrets.yaml
kubectl apply -f base/configmap.yaml
kubectl apply -f base/resource-quotas.yaml
kubectl apply -f base/pdb.yaml

echo "Deploying databases..."
kubectl apply -f databases/postgres-account.yaml
kubectl apply -f databases/postgres-order.yaml
kubectl apply -f databases/elasticsearch.yaml
kubectl apply -f databases/kafka.yaml

echo "Waiting for databases to be ready..."
kubectl wait --for=condition=ready pod -l app=postgres-account -n $NAMESPACE --timeout=300s
kubectl wait --for=condition=ready pod -l app=postgres-order -n $NAMESPACE --timeout=300s
kubectl wait --for=condition=ready pod -l app=elasticsearch -n $NAMESPACE --timeout=300s
kubectl wait --for=condition=ready pod -l app=zookeeper -n $NAMESPACE --timeout=300s
kubectl wait --for=condition=ready pod -l app=kafka -n $NAMESPACE --timeout=300s

echo "Initializing Kafka topics..."
kubectl apply -f databases/kafka-topics.yaml

echo "Deploying microservices..."
kubectl apply -f services/account.yaml
kubectl apply -f services/catalog.yaml
kubectl apply -f services/order.yaml
kubectl apply -f services/graphql.yaml

echo "Deploying ingress..."
kubectl apply -f ingress/ingress.yaml

echo "Deploying monitoring stack..."
kubectl apply -f monitoring/prometheus.yaml
kubectl apply -f monitoring/exporters.yaml
kubectl apply -f monitoring/grafana.yaml
kubectl apply -f monitoring/loki.yaml
kubectl apply -f monitoring/alertmanager.yaml
kubectl apply -f monitoring/alerts.yaml

echo "Deploying backup configurations..."
kubectl apply -f backup/postgres-backup.yaml
kubectl apply -f backup/kafka-backup.yaml
kubectl apply -f backup/velero-backup.yaml

echo "Waiting for services to be ready..."
kubectl wait --for=condition=ready pod -l app=account -n $NAMESPACE --timeout=300s
kubectl wait --for=condition=ready pod -l app=catalog -n $NAMESPACE --timeout=300s
kubectl wait --for=condition=ready pod -l app=order -n $NAMESPACE --timeout=300s
kubectl wait --for=condition=ready pod -l app=graphql -n $NAMESPACE --timeout=300s

echo "Deployment complete!"
echo ""
echo "Services:"
echo "  GraphQL: kubectl get svc graphql-service -n $NAMESPACE"
echo "  Grafana: kubectl get svc grafana -n $MONITORING_NAMESPACE"
echo "  Prometheus: kubectl get svc prometheus -n $MONITORING_NAMESPACE"
echo ""
echo "Port forwarding examples:"
echo "  GraphQL: kubectl port-forward svc/graphql-service 8080:8080 -n $NAMESPACE"
echo "  Grafana: kubectl port-forward svc/grafana 3000:3000 -n $MONITORING_NAMESPACE"
echo "  Prometheus: kubectl port-forward svc/prometheus 9090:9090 -n $MONITORING_NAMESPACE"
echo "  Kafka: kubectl port-forward svc/kafka-0 9092:9092 -n $NAMESPACE"
