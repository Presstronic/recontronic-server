# Contributing to Recontronic Server

Thank you for your interest in contributing to Recontronic Server! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Requirements](#testing-requirements)
- [Submitting Changes](#submitting-changes)

## Code of Conduct

### Our Pledge

We are committed to providing a welcoming and inspiring community for all. Please be respectful and professional in all interactions.

### Expected Behavior

- Be respectful and inclusive
- Welcome newcomers and help them get started
- Focus on what is best for the community
- Show empathy towards other community members

### Unacceptable Behavior

- Harassment, discrimination, or offensive comments
- Trolling, insulting/derogatory comments
- Public or private harassment
- Publishing others' private information without permission

## Getting Started

### Prerequisites

- Go 1.25.3 or higher
- Docker and Docker Compose
- Git
- Make
- PostgreSQL client (psql) for running migrations
- Understanding of bug bounty hunting and security research ethics

### Technology Stack

Familiarize yourself with these technologies used in the project:

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Language** | Go 1.25.3 | Primary development language |
| **HTTP Router** | Chi v5 | Lightweight, composable HTTP routing |
| **Configuration** | Viper | YAML config with environment overrides |
| **Database** | TimescaleDB (PostgreSQL 16) | Time-series optimized storage |
| **Database Driver** | lib/pq | PostgreSQL driver for Go |
| **Password Hashing** | Argon2id (golang.org/x/crypto) | Secure password hashing |
| **Validation** | go-playground/validator | Struct and input validation |
| **Testing** | Standard Go testing | Unit and integration tests |
| **Linting** | go vet, gofmt, staticcheck | Code quality tools |
| **CI/CD** | GitHub Actions | Automated testing and builds |

### Setting Up Development Environment

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/presstronic/recontronic-server.git
   cd recontronic-server
   ```

2. **Install dependencies**
   ```bash
   make deps
   ```

3. **Set up configuration**
   ```bash
   cp configs/config.yaml configs/config.local.yaml
   # Edit configs/config.local.yaml with your settings
   # Or use environment variables with RECONTRONIC_ prefix
   ```

4. **Start development environment**
   ```bash
   make docker-up
   ```

5. **Run tests to verify setup**
   ```bash
   make test
   ```

## Development Workflow

### Branch Naming Convention

- `feature/` - New features (e.g., `feature/add-nuclei-integration`)
- `bugfix/` - Bug fixes (e.g., `bugfix/fix-scan-timeout`)
- `hotfix/` - Urgent production fixes (e.g., `hotfix/security-patch`)
- `refactor/` - Code refactoring (e.g., `refactor/simplify-diff-engine`)
- `docs/` - Documentation updates (e.g., `docs/update-api-guide`)

### Workflow Steps

1. **Create a new branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Write clean, readable code
   - Follow Go best practices
   - Add tests for new functionality
   - Update documentation as needed

3. **Run quality checks**
   ```bash
   make lint      # Run linter
   make fmt       # Format code
   make vet       # Run go vet
   make test      # Run tests
   ```

4. **Commit your changes**
   ```bash
   git add .
   git commit -m "feat: add nuclei integration for vulnerability scanning"
   ```

5. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Open a Pull Request**
   - Provide a clear description of changes
   - Reference any related issues
   - Ensure CI/CD checks pass

## Coding Standards

### Go Style Guide

Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines and these additional rules:

**File Organization**
- Group imports: standard library, third-party, internal
- Order: constants, variables, types, functions
- Use meaningful package names

**Naming Conventions**
- Use camelCase for unexported identifiers
- Use PascalCase for exported identifiers
- Use descriptive names (avoid single-letter variables except in short scopes)
- Interface names should end with `er` when appropriate (e.g., `Scanner`, `Alerter`)

**Code Structure**
```go
package service

import (
    "context"
    "fmt"

    "github.com/external/package"

    "github.com/presstronic/recontronic-server/internal/models"
)

// Constants
const (
    MaxRetries = 3
    DefaultTimeout = 30 * time.Second
)

// Types
type ScanService struct {
    db     Database
    logger Logger
}

// Constructor
func NewScanService(db Database, logger Logger) *ScanService {
    return &ScanService{
        db:     db,
        logger: logger,
    }
}

// Methods
func (s *ScanService) ExecuteScan(ctx context.Context, programID int) error {
    // Implementation
}
```

**Error Handling**
- Always handle errors explicitly
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Return errors instead of panicking
- Use custom error types for business logic errors

**Comments**
- Write package documentation
- Document exported functions, types, and constants
- Use complete sentences
- Explain WHY, not just WHAT

```go
// ScanService handles reconnaissance scanning operations.
// It coordinates between workers, manages scan state, and
// processes results for anomaly detection.
type ScanService struct {
    // ... fields
}

// ExecuteScan performs a reconnaissance scan for the specified program.
// It returns an error if the program is not found or if the scan fails.
func (s *ScanService) ExecuteScan(ctx context.Context, programID int) error {
    // Implementation
}
```

### Project Structure

Follow the standard Go project layout:

```
internal/       - Private application code
cmd/           - Application entry points
pkg/           - Public libraries (reusable)
api/           - API definitions (proto, OpenAPI)
deployments/   - Deployment configurations
scripts/       - Build and dev scripts
```

## Testing Requirements

### Unit Tests

- Write unit tests for all business logic
- Aim for >80% code coverage
- Use table-driven tests when appropriate
- Mock external dependencies

```go
func TestScanService_ExecuteScan(t *testing.T) {
    tests := []struct {
        name      string
        programID int
        wantErr   bool
    }{
        {
            name:      "valid program",
            programID: 1,
            wantErr:   false,
        },
        {
            name:      "invalid program",
            programID: -1,
            wantErr:   true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Integration Tests

- Place in `test/integration/`
- Use real database (Docker Compose test environment)
- Clean up resources after tests

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run specific package
go test ./internal/services/...

# Run integration tests
make test-integration
```

## Submitting Changes

### Commit Message Format

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, no logic change)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(worker): add nuclei integration for vulnerability scanning

Implements nuclei scanner as a new worker type. This enables
automated vulnerability detection during reconnaissance scans.

Closes #123
```

```
fix(api): resolve race condition in scan endpoint

The scan endpoint had a race condition when multiple requests
were made simultaneously. Added proper mutex locking.

Fixes #456
```

### Pull Request Guidelines

**Before Submitting:**
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex logic
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] All tests passing
- [ ] Linter passing
- [ ] No merge conflicts

**PR Description Should Include:**
1. **Summary**: What does this PR do?
2. **Motivation**: Why is this change needed?
3. **Changes**: What specific changes were made?
4. **Testing**: How was this tested?
5. **Screenshots**: If UI changes (not applicable for backend)
6. **Related Issues**: Link to related issues

**PR Template:**
```markdown
## Summary
Brief description of changes

## Motivation
Why this change is necessary

## Changes
- Change 1
- Change 2
- Change 3

## Testing
How this was tested

## Checklist
- [ ] Tests pass
- [ ] Linter passes
- [ ] Documentation updated
- [ ] Ready for review
```

### Code Review Process

1. **Automated Checks**: CI/CD must pass
2. **Peer Review**: At least one approval required
3. **Maintainer Review**: Final approval from maintainer
4. **Merge**: Squash and merge or rebase

## Additional Guidelines

### Security

- Never commit secrets, API keys, or credentials
- Use environment variables for configuration
- Follow security best practices
- Report security vulnerabilities privately

### Performance

- Profile code for performance-critical sections
- Avoid premature optimization
- Use benchmarks to validate improvements
- Consider memory usage and allocations

### Dependencies

- Minimize external dependencies
- Use well-maintained, reputable packages
- Check licenses for compatibility
- Keep dependencies updated

### Documentation

- Update README.md for user-facing changes
- Add/update API documentation
- Document configuration changes
- Include code examples

## Getting Help

- **Questions**: Open a discussion on GitHub
- **Bugs**: Open an issue with reproduction steps
- **Features**: Open an issue with detailed description
- **Security**: Email security@yourproject.com

## Recognition

Contributors will be recognized in:
- CONTRIBUTORS.md file
- Release notes
- Project documentation

Thank you for contributing to Recontronic Server!
