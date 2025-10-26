# Recontronic Server

> 24/7 automated reconnaissance platform for bug bounty hunters. Never miss that weekend deployment or 2 AM emergency fix again.

[![CI](https://github.com/presstronic/recontronic-server/actions/workflows/ci.yml/badge.svg)](https://github.com/presstronic/recontronic-server/actions/workflows/ci.yml)
[![Go Version](https://img.shields.io/badge/Go-1.25.3-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/presstronic/recontronic-server)](https://goreportcard.com/report/github.com/presstronic/recontronic-server)
[![License](https://img.shields.io/badge/license-TBD-blue.svg)](LICENSE)
[![Status](https://img.shields.io/badge/status-Authentication%20Complete-green)](https://github.com/presstronic/recontronic-server)

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
- Regular tables (users, api_keys, programs)
- Time-series optimized tables (assets, anomalies) with automatic compression
- Future: Job queue for background processing
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
git clone https://github.com/presstronic/recontronic-server.git
cd recontronic-server

# Install dependencies
make deps

# Start database (TimescaleDB + Redis)
make docker-up

# Run migrations
psql -h localhost -U postgres -d recontronic \
  -f migrations/20251026021554_create_users_and_api_keys.up.sql

# Start API server
make run-api

# In another terminal, run tests
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

See the Kubernetes deployment section above for production deployment.

## 📖 Documentation

- **[Authentication Guide](AUTH_IMPLEMENTATION.md)** - Complete authentication system documentation
- **[Project Setup](PROJECT_SETUP.md)** - Detailed setup and architecture information
- **[Contributing Guide](CONTRIBUTING.md)** - How to contribute and development workflow
- **[Quick Start](#-quick-start)** - Get started with development
- **[Usage Examples](#-usage-example)** - API usage examples

## 🛠️ Tech Stack

| Component | Technology | Purpose |
|-----------|------------|---------|
| **Language** | Go 1.25.3 | High-performance, concurrent processing |
| **Database** | TimescaleDB (PostgreSQL 16) | Time-series optimized PostgreSQL |
| **API Router** | Chi v5 | Lightweight, composable HTTP routing |
| **Authentication** | Argon2id + API Keys | Secure password hashing and token-based auth |
| **Validation** | go-playground/validator | Request input validation |
| **Configuration** | Viper | Environment and config management |
| **Containerization** | Docker + Docker Compose | Development and production deployment |
| **Orchestration** | Kubernetes (k3s) | Container orchestration |
| **CI/CD** | GitHub Actions | Automated testing and builds |
| **Recon Tools** | subfinder, httpx | Asset discovery and probing |

## 🎮 Usage Example

### Authentication

```bash
# Register a new user
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "SecureP@ssw0rd123"
  }'

# Login and receive API key
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "SecureP@ssw0rd123"
  }'

# Response:
# {
#   "user": { "id": 1, "username": "johndoe", "email": "john@example.com" },
#   "api_key": "rct_AbCdEf123456...",
#   "message": "Login successful. Save this API key securely."
# }

# Get current user info (protected endpoint)
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer rct_AbCdEf123456..."

# Generate additional API key
curl -X POST http://localhost:8080/api/v1/auth/keys \
  -H "Authorization: Bearer rct_AbCdEf123456..." \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Production Key"
  }'

# List all your API keys
curl -X GET http://localhost:8080/api/v1/auth/keys \
  -H "Authorization: Bearer rct_AbCdEf123456..."

# Revoke an API key
curl -X DELETE http://localhost:8080/api/v1/auth/keys/2 \
  -H "Authorization: Bearer rct_AbCdEf123456..."
```

### Coming Soon: Bug Bounty Program Management

The following features are planned for future releases:

```bash
# Add a bug bounty program (coming soon)
curl -X POST http://localhost:8080/api/v1/programs \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Example Corp",
    "platform": "hackerone",
    "scope": ["*.example.com", "*.example.io"],
    "scan_frequency": "1h"
  }'
```

A dedicated [Recontronic CLI](https://github.com/presstronic/recontronic-cli) is planned for a better developer experience.

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
- [x] Authentication system (API key-based with Argon2id)
- [x] Core API endpoints for user management
- [x] CI/CD pipeline with GitHub Actions
- [x] Docker containerization
- [ ] Worker implementation with subfinder/httpx
- [ ] Diff engine and anomaly detection
- [ ] Discord/Slack alerting
- [ ] Scheduled CronJobs
- [ ] Production deployment

See the [v1.0 MVP Milestone](https://github.com/presstronic/recontronic-server/milestone/1) for detailed progress.

### ✅ Completed: Authentication System

The authentication system is fully implemented and production-ready:

- **User Registration & Login**: Username/password with Argon2id hashing (64MB memory, 3 iterations)
- **API Key Management**: Generate, list, and revoke long-lived API keys
- **Secure Authentication**: Bearer token authentication via middleware
- **Input Validation**: Comprehensive validation with go-playground/validator
- **Database Schema**: Users and API keys tables with proper indexes
- **100% Test Coverage**: All auth utilities thoroughly tested

**Available Endpoints:**
- `POST /api/v1/auth/register` - Create new user account
- `POST /api/v1/auth/login` - Login and receive API key
- `GET /api/v1/auth/me` - Get current user info (protected)
- `POST /api/v1/auth/keys` - Generate additional API keys (protected)
- `GET /api/v1/auth/keys` - List all user API keys (protected)
- `DELETE /api/v1/auth/keys/{id}` - Revoke an API key (protected)

See [AUTH_IMPLEMENTATION.md](AUTH_IMPLEMENTATION.md) for detailed documentation.

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

## 🔧 Configuration

The server uses a YAML configuration file. Create `config.yaml` in the project root:

```yaml
server:
  restport: 8080
  grpcport: 9090  # Reserved for future use
  environment: development
  readtimeout: 15s
  writetimeout: 15s
  idletimeout: 60s

database:
  host: "localhost"
  port: 5432
  user: postgres
  password: postgres
  dbname: recontronic
  sslmode: disable
  maxopenconns: 25
  maxidleconns: 5
  connmaxlifetime: 5m

logging:
  level: "info"  # debug, info, warn, error
  format: "json" # json or text
  output: "stdout"

# Future: Worker and scanning configuration
worker:
  poolsize: 10
  queuetype: postgres

# Future: Alerting configuration
alerting:
  discord:
    enabled: false
    webhookurl: ""
  slack:
    enabled: false
    webhookurl: ""
```

Environment variables can override config values using the `RECONTRONIC_` prefix:
- `RECONTRONIC_DATABASE_PASSWORD=secret` overrides `database.password`
- `RECONTRONIC_SERVER_RESTPORT=9000` overrides `server.restport`
- `RECONTRONIC_LOGGING_LEVEL=debug` overrides `logging.level`

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

For more troubleshooting help, check the GitHub Issues or open a new issue with details about your problem.

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
