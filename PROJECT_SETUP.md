# Project Setup Documentation

This document describes the professional setup completed for the Recontronic Server project.

## Summary

The Recontronic Server has been set up as a professional, production-ready Go application with modern DevOps practices, CI/CD pipelines, and cloud-native deployment configurations.

## What Was Created

### 1. Go Module Initialization
- **Go Version**: 1.25.3 (latest stable release)
- **Module**: `github.com/presstronic/recontronic-server`
- **Dependencies**: Chi v5, Viper, go-playground/validator, lib/pq, golang.org/x/crypto

### 2. Project Structure

```
recontronic-server/
├── cmd/                    # Application entry points
│   ├── api/               # REST/gRPC API server
│   ├── worker/            # Background job workers
│   └── cli/               # CLI client
├── internal/              # Private application code
│   ├── config/            # Configuration management (Viper)
│   ├── database/          # Database layer
│   ├── handlers/          # HTTP/gRPC handlers
│   ├── middleware/        # HTTP middleware
│   ├── models/            # Data models
│   ├── repository/        # Data access layer
│   ├── services/          # Business logic
│   └── workers/           # Background job implementations
├── pkg/                   # Public reusable libraries
│   ├── logger/            # Logging utilities
│   ├── validator/         # Input validation
│   └── utils/             # Common utilities
├── api/                   # API definitions
│   ├── proto/             # Protocol Buffer definitions (gRPC)
│   └── openapi/           # OpenAPI/Swagger specs
├── deployments/           # Deployment configurations
│   ├── docker/            # Docker & Docker Compose
│   ├── k8s/              # Kubernetes manifests
│   └── terraform/         # Infrastructure as Code
├── configs/              # Configuration files
├── migrations/           # Database migrations
├── scripts/              # Build and utility scripts
└── test/                 # Tests
    ├── integration/      # Integration tests
    └── unit/             # Unit tests
```

### 3. Configuration System

**Files Created:**
- `internal/config/config.go` - Configuration struct and loader using Viper
- `configs/config.yaml` - Default configuration file
- `configs/.env.example` - Environment variable template

**Features:**
- YAML-based configuration
- Environment variable override support
- Sensible defaults for all settings
- Support for multiple environments (dev, staging, prod)

**Configuration Sections:**
- Server (port, timeouts, graceful shutdown)
- Database (TimescaleDB connection, pool settings)
- Logging (level, format, output)
- Scanning (future: scan settings)
- Alerts (future: Discord/Slack webhooks)

### 4. Application Entry Points

**API Server** (`cmd/api/main.go`):
- REST API on port 8080 (configurable)
- Chi router with middleware stack
- Authentication system with API key support
- Health check endpoint (`/health`)
- Readiness probe endpoint (`/ready`)
- Graceful shutdown handling
- Configuration loading with Viper
- Database connection pooling

**Implemented Endpoints:**
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `GET /api/v1/auth/me` - Get current user (protected)
- `POST /api/v1/auth/keys` - Generate API key (protected)
- `GET /api/v1/auth/keys` - List API keys (protected)
- `DELETE /api/v1/auth/keys/{id}` - Revoke API key (protected)

**Worker** (`cmd/worker/main.go`):
- Background job processor placeholder
- Graceful shutdown
- Configuration loading
- Ready for reconnaissance tool integration

**CLI** (`cmd/cli/main.go`):
- This is a separate project: `github.com/presstronic/recontronic-cli`
- Server-side CLI placeholder removed

### 5. Docker Configuration

**Multi-stage Dockerfiles:**

**API Dockerfile** (`deployments/docker/Dockerfile`):
- Go 1.25.3 Alpine builder stage
- Minimal Alpine runtime
- Non-root user (security)
- Health checks
- Optimized binary size with `-ldflags="-w -s"`

**Worker Dockerfile** (`deployments/docker/Dockerfile.worker`):
- Includes reconnaissance tools (subfinder, httpx)
- Optimized for scanning workloads
- Non-root user
- Tool installation automated

**Docker Compose** (`deployments/docker/docker-compose.yml`):
- TimescaleDB service with health checks
- Redis service
- API service
- Worker service
- Volume management
- Network isolation
- Environment variable configuration

### 6. Kubernetes Manifests

**Base Manifests** (`deployments/k8s/base/`):

1. **namespace.yaml** - Recontronic namespace
2. **configmap.yaml** - Application configuration
3. **secret.yaml** - Sensitive credentials (to be updated)
4. **timescaledb.yaml** - StatefulSet with persistent storage
5. **redis.yaml** - Redis deployment
6. **api.yaml** - API server deployment with:
   - LoadBalancer service
   - 2 replicas
   - Health/readiness probes
   - HorizontalPodAutoscaler (2-10 replicas)
   - Resource limits
7. **worker.yaml** - Worker deployment with:
   - 3 replicas
   - HorizontalPodAutoscaler (3-20 replicas)
   - Resource limits

**Overlays** (structure created):
- `overlays/dev/` - Development environment
- `overlays/staging/` - Staging environment
- `overlays/prod/` - Production environment

### 7. CI/CD Pipelines

**GitHub Actions Workflows:**

**CI Pipeline** (`.github/workflows/ci.yml`):
- Triggers on push/PR to main/develop
- Lint job with:
  - `go vet` for static analysis
  - `gofmt` for formatting checks
  - `staticcheck` for advanced linting (Go 1.25 compatible)
- Test job with:
  - TimescaleDB service container
  - Redis service container
  - Race detection
  - Code coverage
  - Codecov upload
- Build job (builds API and Worker binaries)
- Docker build job (builds images)
- Caching for Go modules and Docker layers
- ✅ All checks currently passing

**Release Pipeline** (`.github/workflows/release.yml`):
- Triggers on version tags (v*)
- GoReleaser for binary releases
- Docker image build and push to GHCR
- Semantic versioning tags
- Multi-arch support ready

### 8. Makefile

**Professional Makefile** with color output and comprehensive commands:

**Build Commands:**
- `make build` - Build all binaries
- `make build-api` - Build API server
- `make build-worker` - Build worker

**Testing Commands:**
- `make test` - Run all tests
- `make test-coverage` - Generate coverage report
- `make lint` - Run linter
- `make fmt` - Format code
- `make vet` - Run go vet

**Docker Commands:**
- `make docker-build` - Build Docker images
- `make docker-up` - Start Docker Compose
- `make docker-down` - Stop Docker Compose
- `make docker-logs` - View logs

**Kubernetes Commands:**
- `make k8s-deploy` - Deploy to Kubernetes
- `make k8s-delete` - Delete from Kubernetes
- `make k8s-status` - Check deployment status

**Development Commands:**
- `make dev` - Start development server
- `make deps` - Download dependencies
- `make clean` - Clean build artifacts

**Database Commands:**
- `make migrate-up` - Run migrations
- `make migrate-down` - Rollback migrations
- `make migrate-create name=<name>` - Create migration

**Utility Commands:**
- `make proto-gen` - Generate protobuf code
- `make install-tools` - Install dev tools
- `make ci` - Run all CI checks
- `make help` - Display help

### 9. Documentation

**README.md**:
- Comprehensive project overview
- Quick start guide
- Development instructions
- Docker deployment guide
- Kubernetes deployment guide
- API documentation
- Configuration reference
- Contributing guidelines
- Roadmap

**CONTRIBUTING.md**:
- Code of conduct
- Development workflow
- Branch naming conventions
- Coding standards (Go style guide)
- Testing requirements
- Commit message format (Conventional Commits)
- Pull request guidelines
- Code review process

**LICENSE**:
- Copyright notice
- Private development status
- Planned open source release

**PROJECT_SETUP.md** (this file):
- Complete setup documentation
- What was created and why

### 10. Code Quality Tools

**Native Go Linting Tools**:
- `go vet` - Official Go static analysis tool
- `gofmt` - Official Go formatting tool
- `staticcheck` - Advanced static analysis (Go 1.25 compatible)
- Note: golangci-lint was replaced due to Go 1.25 compatibility issues

**Git Configuration**:
- `.gitignore` - Comprehensive ignore rules for:
  - Build artifacts
  - IDE files
  - Environment files
  - Secrets
  - Terraform state
  - Database files
  - Logs and temporary files

## Technology Stack

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| Language | Go | 1.25.3 | Application development |
| Router | Chi | v5.2.3 | HTTP routing and middleware |
| Config | Viper | v1.21.0 | Configuration management |
| Database | TimescaleDB | PostgreSQL 16 | Time-series data storage |
| Database Driver | lib/pq | v1.10.9 | PostgreSQL driver |
| Auth | Argon2id | golang.org/x/crypto | Password hashing |
| Validation | go-playground/validator | v10.28.0 | Input validation |
| Container | Docker | Latest | Containerization |
| Orchestration | Kubernetes | Latest | Container orchestration |
| CI/CD | GitHub Actions | - | Automated testing/deployment |

## Next Steps

### Completed Development Tasks

1. ✅ **Authentication System**
   - ✅ User registration and login
   - ✅ API key generation and management
   - ✅ Argon2id password hashing
   - ✅ Authentication middleware
   - ✅ Database schema (users, api_keys)
   - ✅ Input validation
   - ✅ Unit tests with 100% coverage

2. ✅ **Database Layer**
   - ✅ Database connection pool with TimescaleDB
   - ✅ Migration system (SQL files)
   - ✅ Repository layer for users and API keys
   - ✅ Proper error handling

3. ✅ **API Development**
   - ✅ REST endpoints for authentication
   - ✅ Chi router with middleware
   - ✅ CORS configuration
   - ✅ Request/response handlers
   - ✅ Health and readiness endpoints

### Immediate Next Tasks

1. **Program Management API**
   - Define bug bounty program models
   - Create program CRUD endpoints
   - Implement scope management

2. **Worker Implementation**
   - Integrate job queue system
   - Implement worker pool
   - Add tool integrations (subfinder, httpx)
   - Create scan job handlers

3. **Scan Management**
   - Implement scan service
   - Create asset discovery logic
   - Build diff engine
   - Implement anomaly detection
   - Create scoring algorithm

4. **Testing**
   - Add integration tests for auth endpoints
   - Create test fixtures
   - Add end-to-end tests

### Environment Setup for Development

1. **Install required tools:**
   ```bash
   make install-tools
   ```

2. **Start development environment:**
   ```bash
   make docker-up
   ```

3. **Configure environment:**
   ```bash
   cp configs/.env.example configs/.env
   # Edit configs/.env
   ```

4. **Run the API:**
   ```bash
   make dev
   ```

5. **Run tests:**
   ```bash
   make test
   ```

### Deployment Checklist

**Before deploying to production:**

- [ ] Update all secrets in `deployments/k8s/base/secret.yaml`
- [ ] Configure production database credentials
- [ ] Set up proper API keys
- [ ] Configure alert webhooks (Discord/Slack)
- [ ] Review resource limits in Kubernetes manifests
- [ ] Set up monitoring (Prometheus/Grafana)
- [ ] Configure logging aggregation
- [ ] Set up backup strategy
- [ ] Review security settings
- [ ] Enable TLS/SSL certificates
- [ ] Configure firewall rules
- [ ] Test disaster recovery

## Architecture Decisions

### Why This Structure?

1. **Standard Go Layout**: Follows community best practices
2. **Clean Architecture**: Separation of concerns (handlers, services, repository)
3. **Testability**: Easy to mock and test components
4. **Scalability**: Microservices-ready structure
5. **Cloud-Native**: Kubernetes-first approach

### Why These Tools?

1. **Chi v5**: Lightweight, composable HTTP router with excellent middleware support
2. **Viper**: Industry standard for Go configuration with environment override support
3. **Argon2id**: Latest recommended password hashing algorithm (more secure than bcrypt)
4. **TimescaleDB**: PostgreSQL + time-series optimization (perfect for recon data)
5. **go-playground/validator**: Powerful struct validation with custom rules
6. **Docker**: Standard containerization
7. **Kubernetes**: Production-grade orchestration
8. **GitHub Actions**: Integrated CI/CD with no external dependencies
9. **Native Go Tools**: go vet, gofmt, staticcheck (Go 1.25 compatible)

## Best Practices Implemented

✅ Multi-stage Docker builds (smaller images)
✅ Non-root containers (security)
✅ Health checks and readiness probes
✅ Graceful shutdown handling
✅ Proper error handling patterns
✅ Configuration via environment variables
✅ Secrets management (Kubernetes Secrets)
✅ Resource limits (prevent resource exhaustion)
✅ Horizontal Pod Autoscaling
✅ Persistent storage for databases
✅ Service discovery (Kubernetes DNS)
✅ CI/CD automation
✅ Code quality enforcement (linting)
✅ Test coverage tracking

## Verification

To verify the setup is working:

```bash
# 1. Check Go version
go version  # Should show 1.25.3

# 2. Download dependencies
make deps

# 3. Build all binaries
make build

# 4. Run linter
make lint

# 5. Run tests (when written)
make test

# 6. Start Docker environment
make docker-up

# 7. Check services are running
docker ps
```

## Resources

- [Effective Go](https://golang.org/doc/effective_go)
- [Go Project Layout](https://github.com/golang-standards/project-layout)
- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [TimescaleDB Docs](https://docs.timescale.com/)

---

**Setup completed on**: 2025-10-26
**Go version**: 1.25.3
**Status**: Ready for development

