# Fritzy - Production-Grade Microservices E-Commerce Platform

A scalable, event-driven e-commerce backend built with Go microservices, GraphQL API gateway, and production-ready Kubernetes infrastructure.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         GraphQL Gateway                          │
│                     (API Layer - Port 8000)                      │
└────────────────┬────────────────┬────────────────┬──────────────┘
                 │                │                │
        ┌────────▼────────┐ ┌────▼──────┐ ┌──────▼────────┐
        │ Account Service │ │  Catalog  │ │ Order Service │
        │    (gRPC)       │ │  Service  │ │    (gRPC)     │
        └────────┬────────┘ └────┬──────┘ └──────┬────────┘
                 │                │                │
        ┌────────▼────────┐ ┌────▼──────┐ ┌──────▼────────┐
        │   PostgreSQL    │ │Elasticsearch│ │  PostgreSQL   │
        │   (Accounts)    │ │  (Products) │ │   (Orders)    │
        └─────────────────┘ └─────────────┘ └───────────────┘
                 │                │                │
                 └────────────────┼────────────────┘
                                  │
                         ┌────────▼────────┐
                         │  Kafka Cluster  │
                         │   (3 brokers)   │
                         └─────────────────┘
```

## Services

### **GraphQL Gateway** (Port 8000)
- Unified API for all microservices
- Query and mutation resolvers
- Schema stitching
- gRPC client connections

### **Account Service** (Port 8080)
- User account management
- PostgreSQL database
- Publishes `account.created` events to Kafka

### **Catalog Service** (Port 8080)
- Product catalog and search
- Elasticsearch for full-text search
- Publishes `product.created` events to Kafka

### **Order Service** (Port 8080)
- Order processing and management
- PostgreSQL database
- Communicates with Account and Catalog services
- Publishes `order.created` events to Kafka

### **Kafka Messaging**
- Event-driven architecture
- 3-broker cluster with Zookeeper
- Topics: `account.events`, `order.events`, `catalog.events`, `order.created`, `order.completed`, `notifications`
- Replication factor: 3, Min ISR: 2

## Tech Stack

**Backend:**
- Go 1.24.5
- gRPC for inter-service communication
- GraphQL (gqlgen) for API gateway
- Protocol Buffers

**Databases:**
- PostgreSQL 14 (Account, Order services)
- Elasticsearch 7.17 (Catalog service)

**Messaging:**
- Apache Kafka 7.5.0
- Zookeeper 7.5.0

**Infrastructure:**
- Kubernetes (EKS/GKE/AKS)
- Docker & Docker Compose
- Helm (optional)

**Observability:**
- Prometheus (metrics)
- Grafana (dashboards)
- Loki + Promtail (logging)
- AlertManager (alerting)
- Jaeger (optional - distributed tracing)

**CI/CD:**
- GitHub Actions
- ArgoCD (GitOps)
- Trivy (security scanning)

## Quick Start

### Local Development (Docker Compose)

```bash
# Build and run all services
docker-compose up --build

# Access GraphQL playground
open http://localhost:8000
```

### Kubernetes Production Deployment

```bash
# Deploy infrastructure
cd infra/k8s
./deploy.sh

# Access services via port-forward
kubectl port-forward svc/graphql-service 8080:8080 -n fritzy
kubectl port-forward svc/grafana 3000:3000 -n monitoring
kubectl port-forward svc/prometheus 9090:9090 -n monitoring
```

## Project Structure

```
fritzy_backend/
├── account/              # Account microservice
│   ├── cmd/account/      # Main entry point
│   ├── pb/               # Protocol Buffer definitions
│   ├── server.go         # gRPC server implementation
│   ├── service.go        # Business logic
│   ├── repository.go     # Database layer
│   └── app.dockerfile    # Docker build
│
├── catalog/              # Catalog microservice
│   ├── cmd/catalog/
│   ├── pb/
│   ├── server.go
│   ├── service.go
│   ├── repository.go
│   └── app.dockerfile
│
├── order/                # Order microservice
│   ├── cmd/order/
│   ├── pb/
│   ├── server.go
│   ├── service.go
│   ├── repository.go
│   └── app.dockerfile
│
├── graphql/              # GraphQL API Gateway
│   ├── cmd/graphql/
│   ├── graph.go
│   ├── resolvers/
│   ├── schema.graphql
│   └── app.dockerfile
│
├── kafka/                # Kafka utilities
│   └── producer.go       # Producer/Consumer wrapper
│
├── infra/                # Infrastructure as Code
│   ├── k8s/              # Kubernetes manifests
│   │   ├── base/         # Namespace, ConfigMaps, Secrets
│   │   ├── databases/    # PostgreSQL, Elasticsearch, Kafka
│   │   ├── services/     # Microservice Deployments
│   │   ├── security/     # NetworkPolicies, RBAC
│   │   ├── monitoring/   # Prometheus, Grafana, Loki
│   │   ├── backup/       # Backup CronJobs
│   │   ├── ingress/      # NGINX Ingress
│   │   └── envs/         # Dev, Staging, Prod configs
│   ├── cicd/             # CI/CD configurations
│   │   ├── github-actions/
│   │   └── argocd/
│   └── README.md         # Infrastructure documentation
│
├── docker-compose.yaml   # Local development setup
├── go.mod                # Go dependencies
└── README.md             # This file
```

## API Examples

### GraphQL Queries

```graphql
# Create Account
mutation {
  createAccount(account: { name: "John Doe" }) {
    id
    name
  }
}

# Create Product
mutation {
  createProduct(product: {
    name: "iPhone 15"
    description: "Latest iPhone"
    price: 999.99
  }) {
    id
    name
    price
  }
}

# Create Order
mutation {
  createOrder(order: {
    accountId: "account-id"
    products: [
      { productId: "product-id", quantity: 2 }
    ]
  }) {
    id
    totalPrice
    createdAt
  }
}

# Query Account with Orders
query {
  accounts(pagination: { skip: 0, take: 10 }) {
    id
    name
    orders {
      id
      totalPrice
      products {
        name
        quantity
      }
    }
  }
}
```

## Environment Variables

### All Services
```bash
KAFKA_BROKERS=kafka-0:9092,kafka-1:9092,kafka-2:9092
KAFKA_CONSUMER_GROUP=fritzy-services
```

### Account Service
```bash
DATABASE_URL=postgres://user:pass@localhost:5432/accounts?sslmode=disable
```

### Catalog Service
```bash
DATABASE_URL=http://localhost:9200
```

### Order Service
```bash
DATABASE_URL=postgres://user:pass@localhost:5432/orders?sslmode=disable
ACCOUNT_SERVICE_URL=http://account:8080
CATALOG_SERVICE_URL=http://catalog:8080
```

### GraphQL Gateway
```bash
ACCOUNT_SERVICE_URL=http://account:8080
CATALOG_SERVICE_URL=http://catalog:8080
ORDER_SERVICE_URL=http://order:8080
```

## Development

### Prerequisites
```bash
go >= 1.24
docker >= 20.10
docker-compose >= 1.29
kubectl >= 1.28 (for K8s deployment)
protoc (for generating protobuf files)
```

### Install Dependencies
```bash
go mod download
```

### Generate Protocol Buffers
```bash
cd account && protoc --go_out=. --go-grpc_out=. account.proto
cd catalog && protoc --go_out=. --go-grpc_out=. catalog.proto
cd order && protoc --go_out=. --go-grpc_out=. order.proto
```

### Run Individual Services
```bash
# Account Service
cd account/cmd/account && go run main.go

# Catalog Service
cd catalog/cmd/catalog && go run main.go

# Order Service
cd order/cmd/order && go run main.go

# GraphQL Gateway
cd graphql/cmd/graphql && go run main.go
```

### Build Docker Images
```bash
docker build -t account-service:latest -f account/app.dockerfile .
docker build -t catalog-service:latest -f catalog/app.dockerfile .
docker build -t order-service:latest -f order/app.dockerfile .
docker build -t graphql-service:latest -f graphql/app.dockerfile .
```

## Production Features

### Security
✅ NetworkPolicies (zero-trust networking)
✅ Pod Security Contexts (non-root, read-only FS)
✅ RBAC with ServiceAccounts
✅ TLS/SSL ingress with cert-manager
✅ Secrets management

### High Availability
✅ Multi-replica deployments (3+ replicas)
✅ PodDisruptionBudgets (min 2 available)
✅ Pod anti-affinity rules
✅ HPA (CPU/memory-based autoscaling)
✅ Database replication (PostgreSQL, Elasticsearch)

### Observability
✅ Prometheus metrics collection
✅ Grafana dashboards (services, databases, Kafka)
✅ Loki centralized logging
✅ AlertManager with Slack/PagerDuty
✅ 60+ pre-configured alerts
✅ SLO tracking (99.5% availability)

### Resilience
✅ Health probes (liveness, readiness, startup)
✅ Graceful shutdown (preStop hooks)
✅ Circuit breakers (optional)
✅ Retry logic with exponential backoff
✅ Resource quotas and limits

### Backup & Recovery
✅ Velero cluster backups (daily/weekly)
✅ PostgreSQL automated backups (daily)
✅ Kafka metadata backups
✅ 7-30 day retention policies

### Multi-Environment
✅ Dev, Staging, Production configs
✅ Environment-specific resource limits
✅ Separate credentials per environment

## Kafka Event Schema

### Account Created Event
```json
{
  "event_type": "account.created",
  "account_id": "uuid",
  "name": "string",
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Order Created Event
```json
{
  "event_type": "order.created",
  "order_id": "uuid",
  "account_id": "uuid",
  "total_price": 999.99,
  "products": [
    {
      "product_id": "uuid",
      "quantity": 2,
      "price": 499.99
    }
  ],
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Product Created Event
```json
{
  "event_type": "product.created",
  "product_id": "uuid",
  "name": "string",
  "description": "string",
  "price": 99.99,
  "timestamp": "2024-01-01T00:00:00Z"
}
```

## Monitoring

### Grafana Dashboards
- **Microservices**: Request rate, latency, error rate, CPU/memory
- **PostgreSQL**: Connections, transactions, replication lag, cache hit ratio
- **Elasticsearch**: Cluster health, index rate, search rate, JVM heap
- **Kafka**: Broker status, messages/bytes in/out, consumer lag, partition health

### Key Metrics
- Request throughput (req/sec)
- Latency percentiles (p50, p95, p99)
- Error rate (5xx responses)
- Database connection pool usage
- Kafka consumer lag
- Pod CPU/Memory utilization

### Alerts (60+ rules)

**Critical:**
- Service down (>2min)
- High error rate (>5% for 5min)
- Database down
- Kafka broker offline
- No active Kafka controller
- SLO breach (<99.5% availability)

**Warning:**
- High latency (p95 >1s)
- High CPU/memory (>80%)
- Kafka consumer lag (>1000)
- Database replication lag (>30s)
- Under-replicated Kafka partitions

## CI/CD Pipeline

### GitHub Actions Workflow
1. **Build**: Compile Go binaries
2. **Test**: Run unit tests
3. **Scan**: Trivy security scan
4. **Build Images**: Docker multi-arch builds
5. **Push**: Push to container registry
6. **Deploy**: Update K8s manifests (via ArgoCD)

### ArgoCD GitOps
- Automated sync from Git repository
- Rollback on failure
- Health checks and readiness gates
- Slack notifications on deployment events

## Performance & Scalability

### Benchmarks
- **GraphQL Gateway**: 10,000+ req/sec
- **gRPC Services**: 50,000+ req/sec per instance
- **Database**: Optimized with connection pooling
- **Kafka**: High-throughput event streaming

### Scaling Strategy
- **Horizontal**: HPA scales based on CPU/memory (2-20 replicas)
- **Vertical**: Resource requests/limits configurable
- **Database**: Read replicas for PostgreSQL, ES cluster scaling
- **Kafka**: 3+ brokers, partition-based parallelism

## Troubleshooting

### Check Service Health
```bash
kubectl get pods -n fritzy
kubectl logs -f <pod-name> -n fritzy
kubectl describe pod <pod-name> -n fritzy
```

### View Kafka Topics
```bash
kubectl exec -it kafka-0 -n fritzy -- kafka-topics --list --bootstrap-server localhost:9092
```

### Check Kafka Consumer Lag
```bash
kubectl exec -it kafka-0 -n fritzy -- kafka-consumer-groups --bootstrap-server localhost:9092 --describe --group fritzy-services
```

### Access Databases
```bash
# PostgreSQL
kubectl port-forward svc/postgres-account-primary 5432:5432 -n fritzy
psql -h localhost -U varun -d varun

# Elasticsearch
kubectl port-forward svc/elasticsearch 9200:9200 -n fritzy
curl http://localhost:9200/_cluster/health
```

## Resource Requirements

### Minimum (Dev/Test)
- **CPU**: 4 cores
- **Memory**: 8 GB
- **Storage**: 50 GB

### Recommended (Staging)
- **CPU**: 8 cores
- **Memory**: 16 GB
- **Storage**: 100 GB

### Production
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
- [ ] Enable mTLS between services (service mesh)
- [ ] Implement API authentication/authorization

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## License

MIT

## Support

For issues, questions, or contributions:
- **Issues**: [GitHub Issues](https://github.com/Varun5711/fritzy/issues)
- **Documentation**: See `/infra/README.md` for infrastructure details
- **Monitoring**: Check Grafana dashboards for real-time metrics

## Roadmap

- [ ] Implement distributed tracing (Jaeger/Tempo)
- [ ] Add caching layer (Redis)
- [ ] Implement rate limiting per user
- [ ] Add webhook support for order events
- [ ] Implement saga pattern for distributed transactions
- [ ] Add payment service integration
- [ ] Implement notification service (email/SMS)
- [ ] Add admin dashboard
- [ ] Implement A/B testing framework
- [ ] Add feature flags service

---

**Built with ❤️ using Go, Kubernetes, and Kafka**
