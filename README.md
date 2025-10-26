# Recontronic Server

> 24/7 automated reconnaissance platform for bug bounty hunters. Never miss that weekend deployment or 2 AM emergency fix again.

[![Go Version](https://img.shields.io/badge/Go-1.25.3+-00ADD8?style=flat&logo=go)](https://go.dev/)
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

## ⚠️ Legal Disclaimer

**READ THIS BEFORE USING THIS SOFTWARE**

This tool is designed exclusively for **authorized security testing** as part of legitimate bug bounty programs. By using this software, you agree to the following:

### Authorized Use Only
- ✅ **DO** use only on bug bounty programs where you are enrolled
- ✅ **DO** respect all program rules, scope limitations, and rate limits
- ✅ **DO** obtain explicit authorization before testing any systems
- ✅ **DO** follow responsible disclosure practices
- ❌ **DO NOT** use on systems without explicit permission
- ❌ **DO NOT** exceed the scope defined by bug bounty programs
- ❌ **DO NOT** use for any illegal or unauthorized activities

### Your Responsibility
You are solely responsible for:
- Ensuring you have proper authorization for all testing activities
- Complying with all applicable local, state, national, and international laws
- Following bug bounty program terms of service
- Any consequences of misuse of this software

### No Warranty
THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

**Unauthorized access to computer systems is illegal.** This tool is for security researchers operating within legal boundaries only.

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

**Why TimescaleDB?** It's PostgreSQL with time-series superpowers. One database handles everything:
- Regular tables (programs, users, scan_jobs)
- Time-series optimized tables (assets, anomalies) with automatic compression
- River job queue tables
- 1000x faster time-range queries, 90% storage savings with compression

## 🚀 Quick Start

### Prerequisites

- Go 1.25.3+
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
- **Redis**: localhost:6379 (optional - River uses Postgres)

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
| **Language** | Go 1.25.3+ | High-performance, concurrent processing |
| **Database** | TimescaleDB | Time-series optimized PostgreSQL (one DB for everything!) |
| **Queue** | River | Postgres-backed job queue (no Redis needed) |
| **API** | Chi Router | Lightweight HTTP routing |
| **Orchestration** | Kubernetes (k3s) | Container orchestration |
| **IaC** | Terraform | Infrastructure as Code |
| **Recon Tools** | subfinder, httpx, amass | Asset discovery and probing |

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

**Planned for open source release:** Q1 2026

If you're interested in contributing once the project is public, please watch this repository for updates.

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
make help             # Show all available commands
```

## 🔧 Environment Variables

Key configuration variables (see `.env.example` for complete list):

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=recon_platform
DB_USER=postgres
DB_PASSWORD=your_secure_password

# API
REST_API_PORT=8080
API_KEY=your_api_key_here

# Alerts
DISCORD_ENABLED=true
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/...
SLACK_ENABLED=false
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/...

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Scanning
DEFAULT_SCAN_FREQUENCY=1h
MAX_CONCURRENT_SCANS=5
SCAN_TIMEOUT=30m
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

**Rate limiting issues:**
```bash
# Adjust scan frequency in program configuration
# Reduce MAX_CONCURRENT_SCANS in environment variables
# Add delays between requests (configure in worker)
```

See [docs/troubleshooting.md](docs/troubleshooting.md) for more solutions.

## 🔒 Security Best Practices

When deploying Recontronic:

1. **API Keys**: Use strong, randomly generated API keys (32+ characters)
2. **Secrets Management**: Store credentials in Kubernetes Secrets, never in code
3. **Network Security**: Use firewalls, restrict API access to known IPs
4. **Rate Limiting**: Configure appropriate rate limits to respect target systems
5. **Logging**: Enable audit logging for all scan activities
6. **Updates**: Keep recon tools and dependencies up to date

## 💰 Cost Optimization

Running costs for Recontronic Server:

| Deployment | Monthly Cost | Best For |
|------------|--------------|----------|
| **Local/Dev** | $0 | Development, testing |
| **VPS (Contabo)** | $7-24 | 1-10 programs, personal use |
| **Cloud (DigitalOcean)** | $24-100 | 10-50 programs, team use |
| **Production** | $100-500 | 50+ programs, enterprise |

**Tips to reduce costs:**
- Use a single VPS with k3s instead of managed Kubernetes
- Enable TimescaleDB compression (saves 90% storage)
- Set appropriate data retention policies
- Use spot/preemptible instances for workers

## 📝 License

License: TBD (planning to open source in the future)

Currently all rights reserved. When open sourced, this project will likely use the MIT License to maximize accessibility for the security research community.

## 🙏 Acknowledgments

Built with and inspired by:
- [ProjectDiscovery](https://github.com/projectdiscovery) tools (subfinder, httpx, nuclei) - MIT License
- [TimescaleDB](https://www.timescale.com/) for time-series data - Apache 2.0 License
- [River](https://github.com/riverqueue/river) for job processing - MIT License
- [Chi](https://github.com/go-chi/chi) for HTTP routing - MIT License
- The bug bounty and security research community

Special thanks to all bug bounty hunters who inspired this project.

## 🎓 Learning Resources

If you're new to bug bounty hunting or reconnaissance:

- [HackerOne University](https://www.hacker101.com/)
- [Bugcrowd University](https://www.bugcrowd.com/hackers/bugcrowd-university/)
- [OWASP Testing Guide](https://owasp.org/www-project-web-security-testing-guide/)
- [ProjectDiscovery Blog](https://blog.projectdiscovery.io/)

## 📬 Contact

For questions, issues, or security concerns:
- Open an issue on GitHub
- Email: security@yourproject.com (for security vulnerabilities only)

**Security Disclosure**: If you find a security vulnerability in Recontronic itself, please report it responsibly. Do not open a public issue. Email security@yourproject.com with details.

---

**Built with ❤️ for the bug bounty community**

*Remember: With great automation comes great responsibility. Always hunt ethically and legally.*
