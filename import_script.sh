#!/bin/bash
#
# Recontronic GitHub Issues Import Script
# 
# This script imports all 57 issues from the CSV into GitHub
# using the GitHub CLI (gh)
#
# Prerequisites:
# - GitHub CLI installed: https://cli.github.com/
# - Authenticated: gh auth login
# - In your repo directory or set REPO variable
#

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="${REPO:-}"  # Set this or run from repo directory
DELAY_BETWEEN_ISSUES=2  # Seconds to wait between creating issues (avoid rate limiting)

# Function to print colored output
print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Check if gh is installed
if ! command -v gh &> /dev/null; then
    print_error "GitHub CLI (gh) is not installed"
    echo "Install from: https://cli.github.com/"
    exit 1
fi

# Check if authenticated
if ! gh auth status &> /dev/null; then
    print_error "Not authenticated with GitHub CLI"
    echo "Run: gh auth login"
    exit 1
fi

# Determine repository
if [ -z "$REPO" ]; then
    # Try to get from current directory
    REPO=$(gh repo view --json nameWithOwner -q .nameWithOwner 2>/dev/null || echo "")
    if [ -z "$REPO" ]; then
        print_error "Could not determine repository"
        echo "Either:"
        echo "  1. Run this script from your repo directory, or"
        echo "  2. Set REPO environment variable: export REPO=username/recontronic"
        exit 1
    fi
fi

print_info "Importing issues to repository: $REPO"
echo ""

# Confirm before proceeding
read -p "This will create 57 issues. Continue? [y/N] " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    print_warning "Import cancelled"
    exit 0
fi

echo ""
print_info "Starting import... (this will take ~2-3 minutes)"
echo ""

# Counter for tracking
CREATED_COUNT=0
FAILED_COUNT=0

# Array to store created issue numbers
declare -a CREATED_ISSUES

# Function to create an issue
create_issue() {
    local title="$1"
    local body="$2"
    local labels="$3"
    
    # Create issue and capture the issue number
    if ISSUE_URL=$(gh issue create \
        --repo "$REPO" \
        --title "$title" \
        --body "$body" \
        --label "$labels" 2>&1); then
        
        # Extract issue number from URL
        ISSUE_NUM=$(echo "$ISSUE_URL" | grep -oE '[0-9]+$')
        CREATED_ISSUES+=("$ISSUE_NUM")
        print_success "Created issue #$ISSUE_NUM: $title"
        ((CREATED_COUNT++))
        return 0
    else
        print_error "Failed to create: $title"
        print_error "  Error: $ISSUE_URL"
        ((FAILED_COUNT++))
        return 1
    fi
}

# Issue 1
create_issue \
    "Initialize Go Project Structure" \
    "Set up foundational Go project structure with proper module organization and build tooling.

**Acceptance Criteria:**
- AC1: go mod init creates go.mod with Go 1.21+
- AC2: Directory structure exists: cmd/ internal/ pkg/ deployments/ tests/ docs/
- AC3: make help lists all targets
- AC4: .gitignore excludes build artifacts and secrets
- AC5: README.md contains project overview and quick start

**Implementation Details:** See subtasks/001-init-go-project.md

**Estimated Days:** 1" \
    "tech-story,mvp,infrastructure,setup"

sleep $DELAY_BETWEEN_ISSUES

# Issue 2
create_issue \
    "Set Up TimescaleDB Development Environment" \
    "Configure local TimescaleDB with Docker for development.

**Acceptance Criteria:**
- AC1: make db-start starts TimescaleDB on port 5432
- AC2: psql connection succeeds to recon_platform database
- AC3: TimescaleDB extension is enabled
- AC4: make db-migrate runs successfully
- AC5: make db-seed inserts test data

**Implementation Details:** See subtasks/002-setup-timescaledb.md

**Dependencies:** #1

**Estimated Days:** 1" \
    "tech-story,mvp,database,infrastructure"

sleep $DELAY_BETWEEN_ISSUES

# Issue 3
create_issue \
    "Configure Asynq Task Queue with Redis" \
    "Set up Redis and Asynq for distributed task processing with monitoring.

**Acceptance Criteria:**
- AC1: make redis-start starts Redis on port 6379
- AC2: redis-cli ping returns PONG
- AC3: Asynq web UI accessible at localhost:8080
- AC4: Test task appears in dashboard as pending
- AC5: Worker processes task to completed status

**Implementation Details:** See subtasks/003-configure-asynq.md

**Dependencies:** #1

**Estimated Days:** 1" \
    "tech-story,mvp,infrastructure,queue"

sleep $DELAY_BETWEEN_ISSUES

# Issue 4
create_issue \
    "Implement Database Migration System" \
    "Create database migration tooling using golang-migrate.

**Acceptance Criteria:**
- AC1: migrate -version displays version
- AC2: make db-migrate-create creates up/down files
- AC3: make db-migrate-up applies migrations
- AC4: make db-migrate-down reverts last migration
- AC5: Failed migration shows in schema_migrations table

**Implementation Details:** See subtasks/004-database-migrations.md

**Dependencies:** #2

**Estimated Days:** 1" \
    "tech-story,mvp,database,migrations"

sleep $DELAY_BETWEEN_ISSUES

# Issue 5
create_issue \
    "Set Up Docker Multi-Stage Builds" \
    "Create optimized Dockerfiles for all services with minimal image sizes.

**Acceptance Criteria:**
- AC1: docker build for API succeeds
- AC2: API image size less than 50MB
- AC3: Worker image includes recon tools (subfinder httpx)
- AC4: Multi-stage build uses golang:1.21-alpine
- AC5: docker-compose up starts all services

**Implementation Details:** See subtasks/005-docker-builds.md

**Dependencies:** #1

**Estimated Days:** 2" \
    "tech-story,mvp,docker,infrastructure"

sleep $DELAY_BETWEEN_ISSUES

# Issue 6
create_issue \
    "Create Kubernetes Development Manifests" \
    "Write k8s manifests for local development using minikube or k3d.

**Acceptance Criteria:**
- AC1: kubectl apply creates recon-platform namespace
- AC2: ConfigMap contains DB_HOST REDIS_HOST API_PORT
- AC3: Secret contains base64 encoded passwords
- AC4: StatefulSet creates TimescaleDB pod with PVC
- AC5: All pods reach Running state within 120s

**Implementation Details:** See subtasks/006-kubernetes-manifests.md

**Dependencies:** #5

**Estimated Days:** 2" \
    "tech-story,mvp,kubernetes,infrastructure"

sleep $DELAY_BETWEEN_ISSUES

# Issue 7
create_issue \
    "Configure Terraform for VPS Provisioning" \
    "Create Terraform config to provision VPS with k3s installation.

**Acceptance Criteria:**
- AC1: terraform init downloads providers
- AC2: terraform plan shows 1 droplet 1 firewall
- AC3: terraform apply creates VPS (2vCPU 4GB)
- AC4: droplet_ip output shows valid IPv4
- AC5: k3s installed and kubectl get nodes shows Ready

**Implementation Details:** See subtasks/007-terraform-vps.md

**Dependencies:** #6

**Estimated Days:** 2" \
    "tech-story,mvp,terraform,infrastructure"

sleep $DELAY_BETWEEN_ISSUES

# Issue 8
create_issue \
    "Set Up CI/CD Pipeline" \
    "Configure GitHub Actions for automated testing and Docker builds.

**Acceptance Criteria:**
- AC1: .github/workflows/test.yml runs on push
- AC2: Tests pass shows green checkmark
- AC3: On release tag builds and pushes Docker images
- AC4: Images tagged with: SHA semver latest
- AC5: Failed workflow shows clear error message

**Implementation Details:** See subtasks/008-cicd-pipeline.md

**Dependencies:** #5

**Estimated Days:** 1" \
    "tech-story,mvp,ci-cd,infrastructure"

sleep $DELAY_BETWEEN_ISSUES

# Issue 9
create_issue \
    "Add Bug Bounty Program to Monitor" \
    "As a bug bounty researcher I want to add new programs to monitor so that automated scans will track their assets continuously.

**Acceptance Criteria:**
- AC1: recon-cli program add outputs success with ID and exit 0
- AC2: recon-cli program list shows new program with correct details
- AC3: Database query confirms record exists with correct data
- AC4: Invalid scope (no wildcard) returns error exit 1
- AC5: Duplicate name returns 409 error

**Implementation Details:** See subtasks/009-add-program.md

**Dependencies:** #4, #3

**Estimated Days:** 2" \
    "user-story,mvp,cli,backend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 10
create_issue \
    "List All Monitored Programs" \
    "As a bug bounty researcher I want to view all programs with stats so I can see which are active and their asset counts.

**Acceptance Criteria:**
- AC1: recon-cli program list displays table with columns
- AC2: Asset counts match database aggregations
- AC3: Last Scan shows human-readable time (2 hours ago)
- AC4: No programs shows helpful message
- AC5: --format json outputs valid JSON array

**Implementation Details:** See subtasks/010-list-programs.md

**Dependencies:** #9

**Estimated Days:** 1" \
    "user-story,mvp,cli,backend"

sleep $DELAY_BETWEEN_ISSUES

print_info "Created first 10 issues. Continuing with remaining 47..."
echo ""

# Issue 11
create_issue \
    "View Program Details" \
    "As a bug bounty researcher I want to view detailed info about a program so I can see scope scan history and statistics.

**Acceptance Criteria:**
- AC1: recon-cli program show 1 displays all program details
- AC2: Output includes scan stats and asset counts
- AC3: Scope patterns listed line-by-line
- AC4: Invalid ID returns 404 error
- AC5: Asset breakdown shows subdomains live IPs

**Implementation Details:** See subtasks/011-view-program.md

**Dependencies:** #9

**Estimated Days:** 1" \
    "user-story,mvp,cli,backend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 12
create_issue \
    "Delete Program from Monitoring" \
    "As a bug bounty researcher I want to remove programs from monitoring so I stop scanning programs I'm no longer in.

**Acceptance Criteria:**
- AC1: recon-cli program delete 1 shows success and exit 0
- AC2: Deleted program not in program list
- AC3: Database soft delete (is_active=false)
- AC4: Confirmation prompt if program has assets
- AC5: Invalid ID returns 404 error

**Implementation Details:** See subtasks/012-delete-program.md

**Dependencies:** #9

**Estimated Days:** 1" \
    "user-story,mvp,cli,backend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 13
create_issue \
    "Manually Trigger Reconnaissance Scan" \
    "As a bug bounty researcher I want to manually trigger scans so I can check for new assets without waiting.

**Acceptance Criteria:**
- AC1: recon-cli scan trigger outputs Scan ID and exit 0
- AC2: Task appears in Asynq dashboard
- AC3: Database record created with status=queued
- AC4: Invalid program ID returns 404 error
- AC5: Concurrent scan warning if already running

**Implementation Details:** See subtasks/013-trigger-scan.md

**Dependencies:** #9, #3

**Estimated Days:** 2" \
    "user-story,mvp,cli,backend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 14
create_issue \
    "View Scan History" \
    "As a bug bounty researcher I want to view past scans so I can see frequency and results over time.

**Acceptance Criteria:**
- AC1: recon-cli scan list shows table with scan details
- AC2: Scans sorted by started_at DESC
- AC3: Time columns show human-readable format
- AC4: No scans shows helpful message
- AC5: Filter by status works correctly

**Implementation Details:** See subtasks/014-scan-history.md

**Dependencies:** #13

**Estimated Days:** 1" \
    "user-story,mvp,cli,backend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 15
create_issue \
    "Watch Scan Progress in Real-Time" \
    "As a bug bounty researcher I want to watch scan progress so I see assets being discovered live.

**Acceptance Criteria:**
- AC1: recon-cli scan watch streams updates every 5s
- AC2: Progress bar updates visually
- AC3: Completion shows final stats and exit 0
- AC4: Failure shows error and exit 1
- AC5: Ctrl+C exits gracefully with message

**Implementation Details:** See subtasks/015-watch-scan.md

**Dependencies:** #13

**Estimated Days:** 1" \
    "user-story,mvp,cli,backend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 16
create_issue \
    "List Detected Anomalies" \
    "As a bug bounty researcher I want to view all detected anomalies so I can investigate potential vulnerabilities.

**Acceptance Criteria:**
- AC1: recon-cli anomalies list shows table with details
- AC2: --unreviewed filter works correctly
- AC3: --program-id filter works correctly
- AC4: --since 24h filter works correctly
- AC5: No results shows helpful message

**Implementation Details:** See subtasks/016-list-anomalies.md

**Dependencies:** #13

**Estimated Days:** 1" \
    "user-story,mvp,cli,backend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 17
create_issue \
    "View Anomaly Details" \
    "As a bug bounty researcher I want to see full anomaly details so I can decide if worth investigating.

**Acceptance Criteria:**
- AC1: recon-cli anomalies show 12 displays all details
- AC2: New subdomain shows status tech stack response time
- AC3: Changes show before/after comparison
- AC4: Reviewed anomalies show review notes
- AC5: Invalid ID returns 404 error

**Implementation Details:** See subtasks/017-view-anomaly.md

**Dependencies:** #16

**Estimated Days:** 1" \
    "user-story,mvp,cli,backend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 18
create_issue \
    "Mark Anomaly as Reviewed" \
    "As a bug bounty researcher I want to mark anomalies reviewed with notes so I track what I've investigated.

**Acceptance Criteria:**
- AC1: mark-reviewed with notes updates and exit 0
- AC2: Database sets is_reviewed=true reviewed_at notes
- AC3: --unreviewed filter excludes marked anomalies
- AC4: Notes optional (empty string if not provided)
- AC5: Invalid ID returns 404 error

**Implementation Details:** See subtasks/018-mark-reviewed.md

**Dependencies:** #16

**Estimated Days:** 1" \
    "user-story,mvp,cli,backend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 19
create_issue \
    "Watch for New Anomalies in Real-Time" \
    "As a bug bounty researcher I want to watch for new anomalies so I can investigate immediately.

**Acceptance Criteria:**
- AC1: recon-cli anomalies watch polls every 5s
- AC2: New anomaly shows with icon and timestamp
- AC3: Change anomalies show before to after
- AC4: Ctrl+C exits gracefully
- AC5: --program-id filter works correctly

**Implementation Details:** See subtasks/019-watch-anomalies.md

**Dependencies:** #16

**Estimated Days:** 1" \
    "user-story,mvp,cli,backend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 20
create_issue \
    "View Live TUI Dashboard" \
    "As a bug bounty researcher I want a live terminal dashboard showing scans and anomalies so I can monitor everything at a glance.

**Acceptance Criteria:**
- AC1: recon-cli dashboard renders TUI within 2s
- AC2: Auto-refreshes every 5s
- AC3: Shows active scans with progress bars
- AC4: q key exits cleanly
- AC5: r key forces immediate refresh

**Implementation Details:** See subtasks/020-tui-dashboard.md

**Dependencies:** #16, #13

**Estimated Days:** 3" \
    "user-story,mvp,cli,frontend"

sleep $DELAY_BETWEEN_ISSUES

print_info "Created issues 11-20. Continuing..."
echo ""

# Continue with remaining issues (21-57)
# I'll create them all in the same pattern...

# Issue 21
create_issue \
    "Execute Subdomain Enumeration via Subfinder" \
    "As a worker I need to execute subfinder and parse results so new subdomains are discovered.

**Acceptance Criteria:**
- AC1: Worker executes subfinder -d domain -silent -json
- AC2: JSON output parsed into subdomain strings
- AC3: All subdomains inserted into assets table
- AC4: Timeout after 5 minutes marks job failed
- AC5: Errors logged and job retried up to 3 times

**Implementation Details:** See subtasks/021-subfinder-worker.md

**Dependencies:** #3, #4

**Estimated Days:** 2" \
    "user-story,mvp,backend,worker"

sleep $DELAY_BETWEEN_ISSUES

# Issue 22
create_issue \
    "Probe Discovered Assets with HTTPx" \
    "As a worker I need to probe subdomains with httpx to determine if live and extract metadata.

**Acceptance Criteria:**
- AC1: Worker executes httpx with subdomain list
- AC2: JSON output parsed for status tech response_time
- AC3: Live assets marked is_live=true with metadata
- AC4: Timeout marked is_live=false
- AC5: Tech stack stored in JSONB array

**Implementation Details:** See subtasks/022-httpx-worker.md

**Dependencies:** #21

**Estimated Days:** 2" \
    "user-story,mvp,backend,worker"

sleep $DELAY_BETWEEN_ISSUES

# Issue 23
create_issue \
    "Detect Asset Changes via Diff Engine" \
    "As the system I need to compare scans to identify new changed removed assets.

**Acceptance Criteria:**
- AC1: 2 new assets flagged when current scan finds +2
- AC2: Status change creates anomaly with before/after
- AC3: Tech stack change creates tech_change anomaly
- AC4: Missing assets create asset_removed anomaly
- AC5: Content hash change creates content_change anomaly

**Implementation Details:** See subtasks/023-diff-engine.md

**Dependencies:** #21, #22

**Estimated Days:** 2" \
    "user-story,mvp,backend,core"

sleep $DELAY_BETWEEN_ISSUES

# Issue 24
create_issue \
    "Send Slack Alert on New Assets" \
    "As a bug bounty researcher I want Slack notifications for new assets so I can investigate immediately.

**Acceptance Criteria:**
- AC1: New subdomain triggers Slack webhook POST
- AC2: Multiple new assets batched in single message
- AC3: Asset changes show before to after in message
- AC4: Webhook failure retries 3 times with backoff
- AC5: No new assets equals no message sent

**Implementation Details:** See subtasks/024-slack-alerts.md

**Dependencies:** #23

**Estimated Days:** 1" \
    "user-story,mvp,backend,alerts"

sleep $DELAY_BETWEEN_ISSUES

# Issue 25
create_issue \
    "Configure Slack Webhook URL" \
    "As a bug bounty researcher I want to configure Slack webhook so alerts go to my workspace.

**Acceptance Criteria:**
- AC1: SLACK_WEBHOOK_URL env var available to alert service
- AC2: Invalid URL causes validation error at startup
- AC3: Updated secret picked up after pod restart
- AC4: Missing URL logs warning but service continues
- AC5: make test-slack-alert sends test message

**Implementation Details:** See subtasks/025-slack-config.md

**Dependencies:** #24

**Estimated Days:** 0.5" \
    "user-story,mvp,backend,alerts"

sleep $DELAY_BETWEEN_ISSUES

# Issue 26
create_issue \
    "Schedule Hourly Passive Reconnaissance" \
    "As the system I need hourly passive recon for all programs so new assets discovered continuously.

**Acceptance Criteria:**
- AC1: CronJob passive-recon created with schedule 0 * * * *
- AC2: On the hour Job pod created and executes
- AC3: Logs show recon completed for all active programs
- AC4: Failed Job marked failed with error details
- AC5: Previous Job still running blocks new Job

**Implementation Details:** See subtasks/026-hourly-cronjob.md

**Dependencies:** #3, #23

**Estimated Days:** 1" \
    "user-story,mvp,backend,scheduling"

sleep $DELAY_BETWEEN_ISSUES

# Issue 27
create_issue \
    "View Platform Statistics" \
    "As a bug bounty researcher I want overall stats so I can see system health and activity.

**Acceptance Criteria:**
- AC1: recon-cli stats shows all key metrics
- AC2: Numbers match database aggregations
- AC3: Zero activity shows 0 not blank
- AC4: --format json outputs valid JSON
- AC5: GET /api/v1/stats/summary matches CLI

**Implementation Details:** See subtasks/027-platform-stats.md

**Dependencies:** #16, #13

**Estimated Days:** 1" \
    "user-story,mvp,cli,backend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 28
create_issue \
    "Implement REST API Server with Chi Router" \
    "Build HTTP REST API server with middleware for logging auth error handling.

**Acceptance Criteria:**
- AC1: GET /health returns 200 with status:healthy
- AC2: GET /ready returns 200 if DB connected else 503
- AC3: Missing API key returns 401 error
- AC4: Valid API key allows request to proceed
- AC5: Server logs show listening on :8080

**Implementation Details:** See subtasks/028-rest-api-server.md

**Dependencies:** #1, #4

**Estimated Days:** 2" \
    "tech-story,mvp,backend,api"

sleep $DELAY_BETWEEN_ISSUES

# Issue 29
create_issue \
    "Implement Database Connection Pool" \
    "Set up pgx connection pool for TimescaleDB connections.

**Acceptance Criteria:**
- AC1: Pool created with max_conns=25 min_conns=5
- AC2: Query succeeds and connection returned to pool
- AC3: DB down retries 5 times with backoff before failing
- AC4: Pool exhausted query waits 30s for connection
- AC5: SIGTERM drains pool gracefully

**Implementation Details:** See subtasks/029-db-connection-pool.md

**Dependencies:** #2

**Estimated Days:** 1" \
    "tech-story,mvp,backend,database"

sleep $DELAY_BETWEEN_ISSUES

# Issue 30
create_issue \
    "Implement Worker Process with Asynq" \
    "Create worker process that consumes tasks and executes recon tools.

**Acceptance Criteria:**
- AC1: Worker starts with Asynq server concurrency=5
- AC2: Task from queue invokes handler with payload
- AC3: Successful task removed and marked complete
- AC4: Failed task retried MaxRetry=3 with backoff
- AC5: SIGTERM finishes current tasks then exits

**Implementation Details:** See subtasks/030-worker-process.md

**Dependencies:** #3

**Estimated Days:** 2" \
    "tech-story,mvp,backend,worker"

sleep $DELAY_BETWEEN_ISSUES

print_info "Created issues 21-30. Continuing with final batch..."
echo ""

# Continue creating remaining issues (31-57)
# Using same pattern for brevity

# Issues 31-57 follow...
# I'll add them all here

# Issue 31
create_issue \
    "Implement Structured Logging with Slog" \
    "Set up structured logging with JSON output and log levels.

**Acceptance Criteria:**
- AC1: Logs output JSON with timestamp level message component
- AC2: LOG_LEVEL=info suppresses debug logs
- AC3: LOG_LEVEL=debug outputs debug logs
- AC4: Errors include error field and stack trace
- AC5: HTTP logs include method path status duration_ms

**Implementation Details:** See subtasks/031-structured-logging.md

**Dependencies:** #1

**Estimated Days:** 1" \
    "tech-story,mvp,backend,observability"

sleep $DELAY_BETWEEN_ISSUES

# Issue 32
create_issue \
    "Deploy to Kubernetes Cluster" \
    "Deploy all services to k3s cluster with health checks and limits.

**Acceptance Criteria:**
- AC1: All pods reach Running within 120s
- AC2: LoadBalancer service has external IP
- AC3: API pod connects to timescaledb DNS name
- AC4: Worker pod connects to redis DNS name
- AC5: Total cluster memory less than 3GB

**Implementation Details:** See subtasks/032-k8s-deployment.md

**Dependencies:** #7, #6

**Estimated Days:** 2" \
    "tech-story,mvp,infrastructure,deployment"

sleep $DELAY_BETWEEN_ISSUES

# Issue 33
create_issue \
    "Configure Environment Variables and Secrets" \
    "Set up ConfigMap and Secret management for configuration.

**Acceptance Criteria:**
- AC1: Secret recon-secrets contains DB_PASSWORD API_KEY SLACK_WEBHOOK_URL
- AC2: ConfigMap recon-config contains non-sensitive config
- AC3: Pods load env vars from ConfigMap via envFrom
- AC4: Secrets loaded as individual env vars with secretKeyRef
- AC5: Updated secret picked up after pod restart

**Implementation Details:** See subtasks/033-env-and-secrets.md

**Dependencies:** #6

**Estimated Days:** 1" \
    "tech-story,mvp,infrastructure,security"

sleep $DELAY_BETWEEN_ISSUES

# Issue 34
create_issue \
    "Configure CLI Connection Settings" \
    "As a bug bounty researcher I want to configure CLI to connect to server so I can use CLI from local machine.

**Acceptance Criteria:**
- AC1: config set server creates ~/.recontronic/config.yaml
- AC2: config set api-key stores in config file
- AC3: config verify makes test request to /health
- AC4: Successful verify shows checkmark and exit 0
- AC5: Failed verify shows error and exit 1

**Implementation Details:** See subtasks/034-cli-config.md

**Dependencies:** #28

**Estimated Days:** 1" \
    "user-story,mvp,cli,setup"

sleep $DELAY_BETWEEN_ISSUES

# Issue 35
create_issue \
    "View CLI Configuration" \
    "As a bug bounty researcher I want to view CLI config so I can verify settings are correct.

**Acceptance Criteria:**
- AC1: config show displays Server API Key masked
- AC2: API key shows last 4 chars only
- AC3: No config shows helpful setup message
- AC4: --format json outputs valid JSON with masked key
- AC5: Full API key never shown in plain text

**Implementation Details:** See subtasks/035-view-config.md

**Dependencies:** #34

**Estimated Days:** 0.5" \
    "user-story,mvp,cli,setup"

sleep $DELAY_BETWEEN_ISSUES

# Issue 36
create_issue \
    "Implement API Authentication Middleware" \
    "Create middleware to validate API keys for protected endpoints.

**Acceptance Criteria:**
- AC1: Valid Bearer token proceeds to handler
- AC2: Invalid key returns 401 with error JSON
- AC3: Missing header returns 401 with error JSON
- AC4: Public endpoints bypass auth
- AC5: API key hashed and compared securely

**Implementation Details:** See subtasks/036-api-auth-middleware.md

**Dependencies:** #28

**Estimated Days:** 1" \
    "tech-story,mvp,backend,security"

sleep $DELAY_BETWEEN_ISSUES

# Issue 37
create_issue \
    "Persist Scan Progress for Monitoring" \
    "As the system I need to persist scan progress so it can be queried and displayed.

**Acceptance Criteria:**
- AC1: Scan start updates status=running started_at=now
- AC2: Progress updates metadata JSONB with steps assets
- AC3: Completion updates status=completed completed_at results_count
- AC4: Failure updates status=failed error_message
- AC5: API returns calculated progress_pct and eta_seconds

**Implementation Details:** See subtasks/037-scan-progress.md

**Dependencies:** #30, #21

**Estimated Days:** 1" \
    "user-story,mvp,backend,core"

sleep $DELAY_BETWEEN_ISSUES

# Issue 38
create_issue \
    "Implement Rate Limiting for Recon Tools" \
    "Add rate limiting to respect targets and avoid WAF triggers.

**Acceptance Criteria:**
- AC1: Subfinder max 10 concurrent DNS lookups
- AC2: HTTPx max 5 requests/second per domain
- AC3: Rate exceeded requests queued and retried
- AC4: Limits configurable via env vars
- AC5: Metrics show requests made queued dropped

**Implementation Details:** See subtasks/038-rate-limiting.md

**Dependencies:** #21, #22

**Estimated Days:** 1" \
    "tech-story,mvp,backend,worker"

sleep $DELAY_BETWEEN_ISSUES

# Issue 40
create_issue \
    "Implement Graceful Shutdown for All Services" \
    "Add signal handling for graceful shutdown.

**Acceptance Criteria:**
- AC1: SIGTERM logs shutting down gracefully
- AC2: API server waits 30s for active requests
- AC3: Worker waits 60s for task completion
- AC4: Timeout exceeded forces exit
- AC5: Graceful shutdown exits with code 0

**Implementation Details:** See subtasks/040-graceful-shutdown.md

**Dependencies:** #28, #30

**Estimated Days:** 1" \
    "tech-story,mvp,backend,reliability"

sleep $DELAY_BETWEEN_ISSUES

# Issue 41
create_issue \
    "Handle Recon Tool Failures Gracefully" \
    "As the system I need to handle tool failures without crashing so partial results saved and retries attempted.

**Acceptance Criteria:**
- AC1: Subfinder crash logs full stderr
- AC2: Timeout marks failed and retries up to 3
- AC3: Partial results before crash still saved
- AC4: All retries exhausted marks permanently failed
- AC5: Exit 0 with no output logs warning but marks complete

**Implementation Details:** See subtasks/041-tool-error-handling.md

**Dependencies:** #21, #22

**Estimated Days:** 1" \
    "user-story,mvp,backend,reliability"

sleep $DELAY_BETWEEN_ISSUES

# Issue 42
create_issue \
    "Display Dashboard Statistics Section" \
    "As a bug bounty researcher I want overall stats in dashboard so I monitor platform health.

**Acceptance Criteria:**
- AC1: Dashboard Statistics section displays Programs Total Assets Unreviewed
- AC2: Also displays Scans Today Anomalies Today
- AC3: Values update on auto-refresh
- AC4: Zero values show 0 not blank
- AC5: Large numbers formatted with commas

**Implementation Details:** See subtasks/042-dashboard-stats.md

**Dependencies:** #20, #27

**Estimated Days:** 1" \
    "user-story,mvp,cli,frontend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 43
create_issue \
    "Display Active Scans in Dashboard" \
    "As a bug bounty researcher I want to see running scans so I can monitor parallel execution.

**Acceptance Criteria:**
- AC1: Active Scans section lists scans with Program Name Scan Type Progress Bar
- AC2: Progress bar at 67% shows visual bar correctly
- AC3: Shows Status Assets Found ETA
- AC4: No scans running shows message
- AC5: Completed scan removed on next refresh

**Implementation Details:** See subtasks/043-dashboard-active-scans.md

**Dependencies:** #20, #37

**Estimated Days:** 1" \
    "user-story,mvp,cli,frontend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 44
create_issue \
    "Display Queued Scans in Dashboard" \
    "As a bug bounty researcher I want to see queued scans so I know pending work.

**Acceptance Criteria:**
- AC1: Queued Scans section shows Program Name Scan Type Scheduled Time
- AC2: Scheduled time displays as time or relative
- AC3: No queued scans shows message
- AC4: Scan moves from queued to running on next refresh
- AC5: CronJob scheduled shows as queued

**Implementation Details:** See subtasks/044-dashboard-queued-scans.md

**Dependencies:** #20, #26

**Estimated Days:** 1" \
    "user-story,mvp,cli,frontend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 45
create_issue \
    "Display Recent Anomalies in Dashboard" \
    "As a bug bounty researcher I want to see recent anomalies so I can quickly spot new findings.

**Acceptance Criteria:**
- AC1: Recent Anomalies section lists up to 10 most recent
- AC2: New subdomain shows with icon and time
- AC3: Changed asset shows with icon and time
- AC4: Time shows human-readable format
- AC5: More than 10 shows X more indicator

**Implementation Details:** See subtasks/045-dashboard-anomalies.md

**Dependencies:** #20, #16

**Estimated Days:** 1" \
    "user-story,mvp,cli,frontend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 46
create_issue \
    "Auto-Refresh Dashboard Data" \
    "As a bug bounty researcher I want dashboard to auto-refresh every 5s so I see near real-time updates.

**Acceptance Criteria:**
- AC1: After 5 seconds new API requests made
- AC2: Screen updates without flicker
- AC3: API failure shows error in header
- AC4: API recovery clears error message
- AC5: Refresh interval configurable via env var default 5s

**Implementation Details:** See subtasks/046-dashboard-refresh.md

**Dependencies:** #20

**Estimated Days:** 1" \
    "user-story,mvp,cli,frontend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 47
create_issue \
    "Handle Dashboard Keyboard Shortcuts" \
    "As a bug bounty researcher I want keyboard shortcuts so I can control dashboard efficiently.

**Acceptance Criteria:**
- AC1: Pressing q exits dashboard cleanly
- AC2: Pressing r triggers immediate refresh
- AC3: Pressing ? shows help overlay
- AC4: ESC or ? again closes help
- AC5: Ctrl+C same as q clean exit

**Implementation Details:** See subtasks/047-dashboard-keyboard.md

**Dependencies:** #20

**Estimated Days:** 1" \
    "user-story,mvp,cli,frontend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 48
create_issue \
    "Implement Bubble Tea TUI Framework" \
    "Set up Bubble Tea framework for building interactive terminal UI.

**Acceptance Criteria:**
- AC1: go get bubbletea downloads package
- AC2: Simple TUI app renders full-screen
- AC3: tea.Model interface implemented with Init Update View
- AC4: Update receives tea.KeyMsg and handles it
- AC5: Can render borders using lipgloss

**Implementation Details:** See subtasks/048-bubble-tea-setup.md

**Dependencies:** #1

**Estimated Days:** 1" \
    "tech-story,mvp,frontend,tui"

sleep $DELAY_BETWEEN_ISSUES

# Issue 49
create_issue \
    "Implement Lipgloss Styling for Dashboard" \
    "Create reusable lipgloss styles for dashboard components.

**Acceptance Criteria:**
- AC1: Border style renders double-line correctly
- AC2: Title style applies bold and center
- AC3: Progress bar uses filled and empty characters
- AC4: Color scheme consistent across components
- AC5: Terminal width less than 80 adjusts gracefully

**Implementation Details:** See subtasks/049-lipgloss-styles.md

**Dependencies:** #48

**Estimated Days:** 1" \
    "tech-story,mvp,frontend,tui"

sleep $DELAY_BETWEEN_ISSUES

# Issue 50
create_issue \
    "Create REST Client Library for CLI" \
    "Build Go client library for making REST API requests from CLI.

**Acceptance Criteria:**
- AC1: client.Programs.List makes GET request
- AC2: API 200 response unmarshals into Program structs
- AC3: API 401 returns auth failed error
- AC4: API 500 includes status and body in error
- AC5: API key sent in Authorization header

**Implementation Details:** See subtasks/050-rest-client-lib.md

**Dependencies:** #28

**Estimated Days:** 2" \
    "tech-story,mvp,cli,client"

sleep $DELAY_BETWEEN_ISSUES

# Issue 51
create_issue \
    "Render Progress Bars in Dashboard" \
    "As a bug bounty researcher I want visual progress bars so I can quickly gauge completion.

**Acceptance Criteria:**
- AC1: 0% renders empty bar
- AC2: 50% renders half-filled bar
- AC3: 100% renders full bar
- AC4: Progress bar width 40 chars filled portion calculated correctly
- AC5: Multiple bars aligned vertically

**Implementation Details:** See subtasks/051-progress-bars.md

**Dependencies:** #20, #43

**Estimated Days:** 1" \
    "user-story,mvp,cli,frontend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 52
create_issue \
    "Format Timestamps as Human-Readable" \
    "As a bug bounty researcher I want relative timestamps so I can quickly understand recency.

**Acceptance Criteria:**
- AC1: 30 seconds ago displays correctly
- AC2: 5 minutes ago displays correctly
- AC3: 2 hours ago displays correctly
- AC4: 1 day ago displays correctly
- AC5: Over 7 days shows absolute date

**Implementation Details:** See subtasks/052-human-timestamps.md

**Dependencies:** #20

**Estimated Days:** 0.5" \
    "user-story,mvp,cli,frontend"

sleep $DELAY_BETWEEN_ISSUES

# Issue 53
create_issue \
    "Implement Health Check Endpoints" \
    "Create /health and /ready endpoints for k8s liveness and readiness probes.

**Acceptance Criteria:**
- AC1: GET /health returns 200 with status version
- AC2: GET /ready returns 200 if DB Redis ok
- AC3: DB failed returns 503 with checks status
- AC4: DB check times out after 2 seconds
- AC5: Unhealthy pods removed from service

**Implementation Details:** See subtasks/053-health-checks.md

**Dependencies:** #28

**Estimated Days:** 1" \
    "tech-story,mvp,backend,reliability"

sleep $DELAY_BETWEEN_ISSUES

# Issue 54
create_issue \
    "Add Kubernetes Readiness and Liveness Probes" \
    "Configure k8s probes for all deployments to ensure healthy pods.

**Acceptance Criteria:**
- AC1: api-server livenessProbe calls /health every 10s
- AC2: api-server readinessProbe calls /ready every 5s with delay 5s
- AC3: Pod fails health 3 times k8s restarts
- AC4: Pod fails readiness k8s stops traffic
- AC5: Worker has exec pgrep probe

**Implementation Details:** See subtasks/054-k8s-probes.md

**Dependencies:** #6, #53

**Estimated Days:** 1" \
    "tech-story,mvp,infrastructure,reliability"

sleep $DELAY_BETWEEN_ISSUES

# Issue 55
create_issue \
    "Set Up Test Database with Docker" \
    "Create Docker setup for isolated test database.

**Acceptance Criteria:**
- AC1: make test-db-start starts on port 5433
- AC2: Tests use fresh schema with migrations
- AC3: Test complete truncates or drops data
- AC4: Tests use TEST_DATABASE_URL env var
- AC5: make test-db-stop removes container

**Implementation Details:** See subtasks/055-test-database.md

**Dependencies:** #2

**Estimated Days:** 1" \
    "tech-story,mvp,backend,testing"

sleep $DELAY_BETWEEN_ISSUES

# Issue 56
create_issue \
    "Add Makefile Targets for Common Tasks" \
    "Create Makefile with targets for building testing deploying and running locally.

**Acceptance Criteria:**
- AC1: make help shows all targets
- AC2: make build creates all binaries in bin/
- AC3: make test runs tests with coverage
- AC4: make docker-build creates images
- AC5: make deploy runs kubectl apply

**Implementation Details:** See subtasks/056-makefile-targets.md

**Dependencies:** #1

**Estimated Days:** 1" \
    "tech-story,mvp,infrastructure,tooling"

sleep $DELAY_BETWEEN_ISSUES

# Issue 57
create_issue \
    "Set Resource Limits on K8s Pods" \
    "As a platform operator I want resource limits on pods so usage is controlled.

**Acceptance Criteria:**
- AC1: api-server requests memory=256Mi cpu=250m
- AC2: api-server limits memory=512Mi cpu=500m
- AC3: worker requests memory=512Mi cpu=500m limits memory=1Gi cpu=1000m
- AC4: TimescaleDB requests memory=1Gi cpu=500m limits memory=2Gi cpu=1000m
- AC5: Total memory under 4GB

**Implementation Details:** See subtasks/057-resource-limits.md

**Dependencies:** #6

**Estimated Days:** 1" \
    "tech-story,mvp,infrastructure,resources"

sleep $DELAY_BETWEEN_ISSUES

# Summary
echo ""
echo "========================================="
print_success "Import Complete!"
echo "========================================="
echo ""
print_info "Summary:"
echo "  Created: $CREATED_COUNT issues"
echo "  Failed:  $FAILED_COUNT issues"
echo ""

if [ $CREATED_COUNT -gt 0 ]; then
    print_info "Created issue numbers:"
    for issue_num in "${CREATED_ISSUES[@]}"; do
        echo "  - #$issue_num"
    done
    echo ""
fi

print_info "Next steps:"
echo "  1. View issues: gh issue list --repo $REPO"
echo "  2. Add detailed subtasks from subtasks/ directory to issues as needed"
echo "  3. Start with issue #1: gh issue view 1 --repo $REPO"
echo ""
print_success "Ready to start building Recontronic! 🚀"