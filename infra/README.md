# Fritzy Backend - Production-Grade Kubernetes Infrastructure

## Overview

Production-ready Kubernetes infrastructure with comprehensive security, monitoring, backup, and CI/CD capabilities.

## Architecture

### Services
- **GraphQL Gateway**: API gateway (3-20 replicas)
- **Account Service**: User management (2-10 replicas)
- **Catalog Service**: Product catalog (2-10 replicas)
- **Order Service**: Order processing (2-10 replicas)

### Databases
- **PostgreSQL**: Account/Order DBs with streaming replication (2 replicas)
- **Elasticsearch**: Catalog search (3-node cluster, RF=2)

### Observability
- **Prometheus**: Metrics collection (30-day retention)
- **Grafana**: Dashboards and visualization
- **Loki + Promtail**: Centralized logging
- **AlertManager**: Alert routing to Slack/PagerDuty
- **Exporters**: Postgres, Elasticsearch, Node metrics

## Production Features

### Security
- ✅ NetworkPolicies for service isolation
- ✅ Pod Security Contexts (non-root, read-only FS)
- ✅ RBAC with service accounts
- ✅ Secret management
- ✅ TLS ingress with cert-manager
- ✅ Seccomp profiles

### High Availability
- ✅ PodDisruptionBudgets (min 2 replicas)
- ✅ Pod anti-affinity rules
- ✅ Multi-replica databases
- ✅ HPA (CPU/memory-based autoscaling)
- ✅ Graceful shutdown (preStop hooks)

### Health & Resilience
- ✅ Liveness probes
- ✅ Readiness probes
- ✅ Startup probes
- ✅ Resource limits & requests
- ✅ Resource quotas per namespace

### Backup & Recovery
- ✅ Velero (daily/weekly cluster backups)
- ✅ PostgreSQL automated backups (daily)
- ✅ Elasticsearch snapshots
- ✅ 7-day retention

### Monitoring & Alerting
- ✅ Service-level metrics
- ✅ Database metrics
- ✅ Infrastructure metrics
- ✅ Custom alerting rules
- ✅ SLO tracking (99.5% availability)
- ✅ Pre-configured Grafana dashboards

### Multi-Environment
- ✅ Dev, Staging, Production configs
- ✅ Environment-specific resource limits
- ✅ Separate credentials per env

### CI/CD
- ✅ GitHub Actions (build/push/scan)
- ✅ ArgoCD GitOps deployment
- ✅ Automated image security scanning (Trivy)
- ✅ Deployment notifications

## Quick Start

### Prerequisites
```bash
kubectl
helm
docker
```

### Deploy Infrastructure
```bash
cd infra/k8s
./deploy.sh
```

### Deploy to Specific Environment
```bash
kubectl apply -f base/
kubectl apply -f security/
kubectl apply -f databases/
kubectl apply -f services/
kubectl apply -f monitoring/
kubectl apply -f backup/
kubectl apply -f envs/prod/kustomization-patch.yaml
```

## Directory Structure

```
infra/
├── k8s/
│   ├── base/              # Namespace, secrets, configmaps, quotas, PDBs
│   ├── databases/         # PostgreSQL, Elasticsearch StatefulSets
│   ├── services/          # Microservice deployments with HPA
│   ├── ingress/           # NGINX ingress + TLS
│   ├── security/          # NetworkPolicies, RBAC
│   ├── monitoring/        # Prometheus, Grafana, Loki, AlertManager
│   ├── backup/            # Velero, DB backup CronJobs
│   ├── envs/              # Dev, staging, prod configs
│   ├── deploy.sh          # Automated deployment
│   ├── cleanup.sh         # Cleanup script
│   └── values.yaml        # Configuration values
├── cicd/
│   ├── github-actions/    # CI/CD workflows
│   └── argocd/            # GitOps application configs
└── README.md
```

## Configuration

### Update Database Credentials
```bash
kubectl create secret generic db-credentials \
  --from-literal=POSTGRES_USER=prod_user \
  --from-literal=POSTGRES_PASSWORD=secure_password \
  --from-literal=POSTGRES_DB=fritzy_prod \
  -n fritzy --dry-run=client -o yaml | kubectl apply -f -
```

### Configure Alerting
Edit `monitoring/alertmanager.yaml`:
- Set Slack webhook URL
- Set PagerDuty service key

### Configure Backups
Edit `backup/velero-backup.yaml`:
- Set S3 bucket name
- Configure AWS credentials

## Monitoring

### Access Grafana
```bash
kubectl port-forward svc/grafana 3000:3000 -n monitoring
```
URL: http://localhost:3000 (admin/admin123)

### Access Prometheus
```bash
kubectl port-forward svc/prometheus 9090:9090 -n monitoring
```

### View Logs (Loki)
Accessible via Grafana → Explore → Loki datasource

## Alerts

### Critical Alerts
- Service down (2min threshold)
- High error rate (>5% for 5min)
- PostgreSQL down
- Elasticsearch cluster RED
- Pod crash looping
- SLO availability breach (<99.5%)

### Warning Alerts
- High latency (p95 >1s)
- High CPU/memory usage (>80%)
- Replication lag (>30s)
- HPA maxed out
- Low cache hit ratio

## CI/CD Workflows

### GitHub Actions

**Build & Push** (`.github/workflows/build-and-push.yml`):
- Builds Docker images
- Scans for vulnerabilities
- Pushes to registry
- Updates manifests

**Deploy** (`.github/workflows/deploy.yml`):
- Manual deployment trigger
- Supports dev/staging/prod
- Includes smoke tests
- Auto-rollback on failure

### ArgoCD

Install ArgoCD:
```bash
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
kubectl apply -f infra/cicd/argocd/
```

Access ArgoCD UI:
```bash
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

## Scaling

### Manual Scaling
```bash
kubectl scale deployment account-service --replicas=10 -n fritzy
```

### HPA Configuration
Already configured with CPU/memory targets (70%/80%)

### Database Scaling
- PostgreSQL: Add read replicas to StatefulSet
- Elasticsearch: Increase data node count

## Backup & Recovery

### Restore from Velero
```bash
velero restore create --from-backup fritzy-daily-backup-20240101
```

### Restore PostgreSQL
```bash
kubectl exec -it postgres-account-0 -n fritzy -- psql -U varun -d varun < backup.sql
```

## Resource Requirements

### Minimum Cluster Specs
- **Nodes**: 3+ (for HA)
- **CPU**: 20+ cores
- **Memory**: 40+ GB
- **Storage**: 200+ GB

### Production Recommendations
- **Nodes**: 5-10 (multi-AZ)
- **CPU**: 40+ cores
- **Memory**: 80+ GB
- **Storage**: 500+ GB SSD

## Security Checklist

- [ ] Update default passwords in secrets
- [ ] Configure TLS certificates (cert-manager)
- [ ] Set up external secrets management (Vault/AWS Secrets Manager)
- [ ] Enable Pod Security Standards
- [ ] Configure OIDC/SSO for Grafana
- [ ] Set up audit logging
- [ ] Implement rate limiting on ingress
- [ ] Configure WAF rules

## Troubleshooting

### Check pod status
```bash
kubectl get pods -n fritzy
kubectl describe pod <pod-name> -n fritzy
kubectl logs <pod-name> -n fritzy
```

### Check HPA status
```bash
kubectl get hpa -n fritzy
```

### View metrics
```bash
kubectl top pods -n fritzy
kubectl top nodes
```

### Check network policies
```bash
kubectl get networkpolicies -n fritzy
```

## Cost Optimization

- Use spot/preemptible instances for non-prod
- Configure cluster autoscaler
- Right-size resource requests/limits
- Use persistent volume claim resizing
- Implement pod priorities
- Enable vertical pod autoscaling

## SLIs/SLOs

### Availability SLO: 99.5%
- Error budget: 0.5% (3.6 hours/month)
- Measurement: 7-day rolling window
- Alert threshold: 99.5%

### Latency SLO: 95th percentile < 1s
- Measurement: 5-minute windows
- Alert threshold: >1s for 5+ minutes

### Success Rate SLO: 99.5%
- 5xx errors < 0.5%
- Measurement: 30-day rolling window

## Support

For issues or questions:
1. Check logs: `kubectl logs <pod> -n fritzy`
2. Check alerts: Grafana → Alerting
3. View metrics: Prometheus/Grafana
4. Check ArgoCD sync status

## License

MIT
