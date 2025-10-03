# Recontronic Server

> 24/7 automated reconnaissance platform for bug bounty hunters. Never miss that weekend deployment or 2 AM emergency fix again.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/license-TBD-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-MVP%20Development-yellow)](https://github.com/yourusername/recontronic-server/milestones)

## 🎯 Overview

Recontronic Server is an intelligent reconnaissance and anomaly detection platform that continuously monitors bug bounty programs, detects changes, scores anomalies using behavioral analysis, and sends real-time alerts. Focus on finding vulnerabilities instead of manual recon.

**Key Features:**
- 🔄 **Continuous Monitoring** - Automated hourly scans of all in-scope assets
- 🧠 **Intelligent Detection** - Behavioral pattern learning and anomaly scoring
- 🚨 **Real-Time Alerts** - Discord/Slack notifications within 30 seconds
- 📊 **Time-Series Analysis** - Track asset changes and deployment patterns over time
- 🐳 **Production Ready** - Kubernetes-native with proper observability

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                   Recontronic Server                    │
│                                                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐            │
│  │ REST API │  │  Worker  │  │Scheduler │            │
│  │  (Chi)   │  │   Pool   │  │(CronJob) │            │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘            │
│       │             │              │                   │
│  ┌────┴─────────────┴──────────────┴─────┐            │
│  │         TimescaleDB + River Queue      │            │
│  └────────────────────────────────────────┘            │
│                                                         │
│  ┌──────────────────────────────────────────┐          │
│  │  Recon Tools: subfinder, httpx, amass    │          │
│  └──────────────────────────────────────────┘          │
└─────────────────────────────────────────────────────────┘
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- kubectl (for k8s deployment)
- Make (optional, for convenience commands)

### Local Development

```bash
# Clone the repository
git clone https://github.com/yourusername/recontronic-server.git
cd recontronic-server

# Copy environment template
cp .env.example .env
# Edit .env with your configuration

# Start all services (DB, Redis, API, Workers)
make dev

# Run migrations
make migrate-up

# Run tests
make test

# Build binaries
make build
```

The API server will be available at `http://localhost:8080`

### Docker Compose (Recommended for Development)

```bash
docker-compose up -d
```

Services:
- **API Server**: http://localhost:8080
- **TimescaleDB**: localhost:5432
- **Redis**: localhost:6379

### Production Deployment (Kubernetes)

```bash
# Using Terraform to provision infrastructure
cd terraform
terraform init
terraform apply

# Deploy to k8s
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/

# Verify deployment
kubectl get pods -n recon-platform
```

See [docs/deployment.md](docs/deployment.md) for detailed production deployment guide.

## 📖 Documentation

- **[API Documentation](docs/api.md)** - REST API endpoints and usage
- **[Architecture Guide](docs/architecture.md)** - System design and components
- **[Development Guide](docs/development.md)** - How to contribute and develop
- **[Deployment Guide](docs/deployment.md)** - Production deployment instructions
- **[Configuration Reference](docs/configuration.md)** - Environment variables and settings

## 🛠️ Tech Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Language** | Go 1.21+ | High-performance, concurrent processing |
| **Database** | TimescaleDB | Time-series optimized PostgreSQL |
| **Queue** | River | Postgres-backed job queue |
| **API** | Chi Router | Lightweight HTTP routing |
| **Orchestration** | Kubernetes (k3s) | Container orchestration |
| **IaC** | Terraform | Infrastructure as Code |
| **Recon Tools** | subfinder, httpx, amass | Asset discovery |

## 🎮 Usage Example

```bash
# Add a bug bounty program
curl -X POST http://localhost:8080/api/v1/programs \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Example Corp",
    "platform": "hackerone",
    "scope": ["*.example.com", "*.example.io"],
    "scan_frequency": "1h"
  }'

# Trigger a manual scan
curl -X POST http://localhost:8080/api/v1/scans \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "program_id": 1,
    "scan_type": "passive"
  }'

# Query anomalies
curl -X GET "http://localhost:8080/api/v1/anomalies?min_priority=70&unreviewed=true" \
  -H "Authorization: Bearer YOUR_API_KEY"
```

Or use the [Recontronic CLI](https://github.com/yourusername/recontronic-cli) for a better experience.

## 🧪 Testing

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run integration tests
make test-integration

# Run specific package tests
go test ./internal/worker/...

# Lint code
make lint
```

## 📊 Project Status

**Current Phase:** MVP Development (v1.0)

- [x] Project setup and architecture design
- [ ] Core API endpoints (in progress)
- [ ] Worker implementation with subfinder/httpx
- [ ] Diff engine and anomaly detection
- [ ] Discord/Slack alerting
- [ ] Scheduled CronJobs
- [ ] Production deployment

See the [v1.0 MVP Milestone](https://github.com/yourusername/recontronic-server/milestone/1) for detailed progress.

## 🤝 Contributing

This project is currently in private development. Contribution guidelines will be published when the project is open sourced.

## 📋 Makefile Commands

```bash
make dev              # Start development environment
make build            # Build all binaries
make test             # Run tests
make test-coverage    # Run tests with coverage report
make lint             # Run linters
make migrate-up       # Apply database migrations
make migrate-down     # Rollback database migrations
make docker-build     # Build Docker images
make clean            # Clean build artifacts
```

## 🔧 Environment Variables

Key configuration variables (see `.env.example` for complete list):

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=recon_platform
DB_PASSWORD=your_secure_password

# API
REST_API_PORT=8080
API_KEY=your_api_key

# Alerts
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/...
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/...

# Logging
LOG_LEVEL=info
```

## 🐛 Troubleshooting

**Database connection fails:**
```bash
# Check if TimescaleDB is running
docker ps | grep timescale

# Verify connection details
psql -h localhost -U postgres -d recon_platform
```

**Workers not processing jobs:**
```bash
# Check worker logs
kubectl logs -l app=worker -n recon-platform

# Verify River jobs in database
psql -c "SELECT * FROM river_job WHERE state = 'available';"
```

**Scans failing:**
```bash
# Check if recon tools are installed
which subfinder httpx

# View worker logs for errors
kubectl logs -l app=worker -n recon-platform --tail=100
```

See [docs/troubleshooting.md](docs/troubleshooting.md) for more solutions.

## 📝 License

License: TBD (planning to open source in the future)

Currently all rights reserved.

## 🙏 Acknowledgments

Built with:
- [ProjectDiscovery](https://github.com/projectdiscovery) tools (subfinder, httpx, nuclei)
- [TimescaleDB](https://www.timescale.com/) for time-series data
- [River](https://github.com/riverqueue/river) for job processing
- [Chi](https://github.com/go-chi/chi) for HTTP routing

## 📬 Contact

For questions or issues, please open an issue on GitHub.

---

**⚠️ Disclaimer:** This tool is intended for authorized security testing only. Always obtain proper authorization before testing any systems. Follow responsible disclosure practices and bug bounty program rules.
