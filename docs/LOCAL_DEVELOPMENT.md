# Local Development Setup

This guide helps you set up your local development environment for changie.

## Quick Start

### 1. Clone the Repository

```bash
git clone git@github.com:peiman/changie.git
cd changie
```

### 2. Install Development Tools

```bash
task setup
```

This installs:
- goimports (code formatting)
- golangci-lint (linting)
- gotestsum (test runner)
- govulncheck (vulnerability scanner)
- lefthook (git hooks)
- go-mod-outdated (dependency checker)

### 3. Set Up Environment Variables

```bash
# Copy the example file
cp .env.example .env

# Edit .env with your tokens
nano .env  # or vim, code, etc.
```

Add your GitHub token:
```bash
GITHUB_TOKEN=github_pat_xxxxxxxxxxxxx
```

**How to get the token**: See [Token Setup](#github-token-setup) below.

### 4. Load Environment Variables

Before running release tasks, load the .env file:

```bash
source .env
```

Or add to your shell profile (~/.zshrc or ~/.bashrc):
```bash
# Auto-load .env if present
if [ -f .env ]; then
  set -a
  source .env
  set +a
fi
```

### 5. Run Quality Checks

```bash
task check
```

This runs:
- Code formatting (goimports, gofmt)
- Linting (golangci-lint)
- Config validation (check-defaults)
- Dependency verification
- Full test suite with coverage
- Vulnerability scan

---

## GitHub Token Setup

### Create Fine-Grained Token

1. Go to: https://github.com/settings/tokens?type=beta
2. Click **"Generate new token"** → **Fine-grained tokens** tab
3. Configure:
   - **Token name**: `changie-releases`
   - **Expiration**: 90 days (or custom)
   - **Resource owner**: peiman
   - **Repository access**: Only select repositories → `changie`
   - **Permissions**:
     - ✅ **Contents**: Read and write
     - ✅ **Metadata**: Read-only (auto-selected)
4. Click **"Generate token"**
5. **Copy the token** (you won't see it again!)

### Add to .env

```bash
# In .env file:
GITHUB_TOKEN=github_pat_xxxxxxxxxxxxx
```

### Test It

```bash
source .env
task release:dry-run
```

---

## Homebrew Tap Token (Optional)

Only needed if you want to test Homebrew formula updates locally.

### Create Token

1. Go to: https://github.com/settings/tokens?type=beta
2. Click **"Generate new token"** → **Fine-grained tokens**
3. Configure:
   - **Token name**: `homebrew-tap-bot`
   - **Expiration**: 90 days
   - **Resource owner**: peiman
   - **Repository access**: Only select repositories → `homebrew-tap`
   - **Permissions**:
     - ✅ **Contents**: Read and write
     - ✅ **Metadata**: Read-only
4. Generate and copy

### Add to .env

```bash
# In .env file:
HOMEBREW_TAP_GITHUB_TOKEN=github_pat_xxxxxxxxxxxxx
```

---

## Common Development Tasks

### Build and Run

```bash
# Build the binary
task build

# Run directly
./changie --help

# Or build and run
task run
```

### Testing

```bash
# Run all tests
task test

# Run tests with race detection
task test:race

# Watch mode (auto-rerun on changes)
task test:watch

# Coverage report (text)
task test:coverage-text

# Coverage report (HTML in browser)
task test:coverage-html
```

### Code Quality

```bash
# Format code
task format

# Run linters
task lint

# Check for vulnerabilities
task vuln

# Verify dependencies
task deps:verify

# Check for outdated deps
task deps:outdated

# Run all checks
task check
```

### Dependencies

```bash
# Update go.mod and go.sum
task tidy

# Verify dependencies unchanged
task deps:verify

# Check for outdated dependencies
task deps:outdated

# Run all dependency checks
task deps:check
```

### Documentation

```bash
# Generate config documentation (markdown)
task docs:config

# Generate YAML config template
task docs:config-yaml
```

### Release Tasks (Local Testing)

```bash
# Install goreleaser
task release:setup

# Validate .goreleaser.yml
task release:check

# Build snapshot (local, no publishing)
task release:snapshot

# Test the built binary
task release:test

# Full dry run
source .env
task release:dry-run

# Clean release artifacts
task release:clean
```

### Release Tasks (Publishing)

⚠️ **Requires GITHUB_TOKEN in .env**

```bash
# Create and push a tag
task release:tag TAG=v1.0.0

# Or publish release manually
source .env
task release
```

---

## Pre-commit Hooks

Lefthook runs checks before each commit:

- ✅ Format code (goimports, gofmt)
- ✅ Check for unauthorized viper.SetDefault()
- ✅ Lint code (go vet, golangci-lint)
- ✅ Verify dependencies (go mod verify)
- ✅ Run tests

**To skip hooks** (not recommended):
```bash
git commit --no-verify -m "message"
```

**To update hooks**:
```bash
lefthook install
```

---

## Environment Variables

The project uses `.env` for local development:

### Required for Releases

```bash
GITHUB_TOKEN=github_pat_xxxxxxxxxxxxx
```

### Optional

```bash
# For Homebrew tap updates
HOMEBREW_TAP_GITHUB_TOKEN=github_pat_xxxxxxxxxxxxx

# For Arch Linux (AUR) packages
AUR_PRIVATE_KEY=xxxxxxxxxxxxx
```

### Using .env

**Load once per shell session:**
```bash
source .env
task release:dry-run
```

**Or auto-load in shell profile** (~/.zshrc):
```bash
# Add this to auto-load when entering the directory
if [ -f .env ]; then
  set -a
  source .env
  set +a
fi
```

**Or use direnv** (automatic):
```bash
brew install direnv
echo "dotenv" > .envrc
direnv allow
```

---

## Project Structure

```
changie/
├── cmd/                    # CLI commands (thin wrappers)
│   ├── root.go            # Root command + config init
│   ├── version.go         # Version bump commands
│   └── ...
├── internal/              # Business logic (private packages)
│   ├── version/           # Version bump logic
│   ├── changelog/         # Changelog operations
│   ├── git/               # Git wrapper
│   ├── semver/            # SemVer operations
│   ├── config/            # Config registry
│   └── ...
├── docs/                  # Documentation
├── scripts/               # Build scripts
├── .github/workflows/     # CI/CD
├── Taskfile.yml          # Task automation
├── .goreleaser.yml       # Release config
├── .env                  # Local secrets (gitignored)
└── .env.example          # Template
```

---

## Troubleshooting

### "goreleaser: command not found"

```bash
task release:setup
```

### "GITHUB_TOKEN not set"

```bash
# Make sure .env exists and has your token
cat .env | grep GITHUB_TOKEN

# Load it
source .env

# Verify
echo $GITHUB_TOKEN
```

### "git is in a dirty state"

```bash
git status
git add .
git commit -m "chore: prepare release"
```

### Tests failing

```bash
# Run with verbose output
go test -v ./...

# Check specific package
go test -v ./internal/version/
```

### Linter errors

```bash
# Fix formatting first
task format

# Then check linter
task lint

# See what's wrong
golangci-lint run --verbose
```

---

## Quick Reference

### Essential Commands

```bash
# Setup
task setup

# Development
task build
task test
task check

# Release (local test)
source .env
task release:snapshot
task release:test

# Release (publish)
source .env
task release:tag TAG=v1.0.0
```

### File Locations

- **Environment**: `.env` (your secrets)
- **Config**: `internal/config/registry.go`
- **Tasks**: `Taskfile.yml`
- **CI**: `.github/workflows/ci.yml`
- **Hooks**: `.lefthook.yml`
- **Release**: `.goreleaser.yml`

---

## Getting Help

- **Project Docs**: See `docs/` directory
- **Task List**: `task --list`
- **Command Help**: `./changie --help`
- **Issues**: https://github.com/peiman/changie/issues

---

**Last Updated**: 2025-10-06
