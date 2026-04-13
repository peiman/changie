#!/usr/bin/env bash
# CI/CD Integration
#
# changie is a developer-side tool. The CI role is limited to two things:
#   1. Validate changelog quality on PRs (changie changelog validate)
#   2. Extract release notes from changelog (changie diff)
#
# Version bumping happens locally — the developer runs changie bump,
# which commits and tags. CI picks up the tag and publishes via goreleaser.

set -e

echo "=== CI Integration: PR Validation ==="
echo
cat <<'GITHUB_ACTIONS_PR'
# .github/workflows/pr.yml
name: PR Checks

on:
  pull_request:
    branches: [main]

jobs:
  changelog:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: stable

      - name: Install changie
        run: go install github.com/peiman/changie@latest

      - name: Validate changelog
        run: changie changelog validate
GITHUB_ACTIONS_PR

echo
echo "=== CI Integration: Release Notes Extraction ==="
echo
cat <<'RELEASE_NOTES'
# Extract what changed between two versions for notifications
# Useful in post-release workflows

changie diff v0.9.0 v0.9.1
changie diff v0.9.0 v0.9.1 --output json
RELEASE_NOTES

echo
echo "=== Release Flow (developer-side, not CI) ==="
echo
cat <<'RELEASE_FLOW'
# 1. Developer bumps version locally
changie bump minor --auto-push

# 2. Tag push triggers goreleaser in CI
# .github/workflows/release.yml (triggered by tag push)
# goreleaser builds binaries, creates GitHub release, publishes to Homebrew

# This separation is intentional:
# - changie bump requires a clean working directory and git write access
# - goreleaser requires build toolchain and publish credentials
# - Keeping them separate avoids CI pushing commits back to the repo
RELEASE_FLOW
