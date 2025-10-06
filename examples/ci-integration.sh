#!/usr/bin/env bash
# CI/CD Integration Example
#
# This script demonstrates how to integrate changie into
# continuous integration and deployment pipelines.

set -e

echo "=== CI/CD Integration Examples ==="
echo

echo "1. GitHub Actions Workflow"
echo "----------------------------"
cat <<'GITHUB_ACTIONS'
# .github/workflows/release.yml
name: Release

on:
  push:
    branches: [main]
  workflow_dispatch:
    inputs:
      bump_type:
        description: 'Version bump type'
        required: true
        default: 'patch'
        type: choice
        options:
          - major
          - minor
          - patch

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Need full history for git tags

      - name: Install changie
        run: |
          go install github.com/peiman/changie@latest

      - name: Configure git
        run: |
          git config user.name "github-actions[bot]"
          git config user.email "github-actions[bot]@users.noreply.github.com"

      - name: Bump version
        env:
          BUMP_TYPE: ${{ inputs.bump_type || 'patch' }}
        run: |
          changie bump $BUMP_TYPE --json > release.json
          cat release.json

      - name: Extract version
        id: version
        run: |
          VERSION=$(jq -r '.new_version' release.json)
          echo "version=$VERSION" >> $GITHUB_OUTPUT

      - name: Push changes
        run: |
          git push origin main
          git push origin --tags

      - name: Create GitHub Release
        uses: actions/create-release@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
        with:
          tag_name: ${{ steps.version.outputs.version }}
          release_name: Release ${{ steps.version.outputs.version }}
          draft: false
          prerelease: false
GITHUB_ACTIONS

echo
echo "2. GitLab CI Pipeline"
echo "---------------------"
cat <<'GITLAB_CI'
# .gitlab-ci.yml
stages:
  - release

release:
  stage: release
  image: golang:1.21
  only:
    - main
  before_script:
    - go install github.com/peiman/changie@latest
    - git config user.name "GitLab CI"
    - git config user.email "ci@gitlab.com"
  script:
    - changie bump patch --json > release.json
    - cat release.json
    - git push origin main
    - git push origin --tags
  artifacts:
    paths:
      - release.json
    expire_in: 1 week
GITLAB_CI

echo
echo "3. Automated Release with Changelog Entries"
echo "--------------------------------------------"
cat <<'AUTO_RELEASE'
#!/bin/bash
# auto-release.sh

# Check if there are unreleased entries
if ! grep -q "## \[Unreleased\]" CHANGELOG.md; then
  echo "No unreleased changes found"
  exit 0
fi

# Determine bump type based on commits or changelog sections
if grep -q "### Added" CHANGELOG.md && grep -A5 "## \[Unreleased\]" CHANGELOG.md | grep -q "^- "; then
  BUMP_TYPE="minor"
elif grep -q "### Fixed" CHANGELOG.md; then
  BUMP_TYPE="patch"
else
  echo "No changes to release"
  exit 0
fi

# Perform release
changie bump $BUMP_TYPE --auto-push --json
AUTO_RELEASE

echo
echo "4. Pre-release Check Script"
echo "----------------------------"
cat <<'PRERELEASE'
#!/bin/bash
# prerelease-check.sh
# Run this before releasing to validate everything is ready

set -e

echo "🔍 Pre-release validation"

# Check for uncommitted changes
if [[ -n $(git status -s) ]]; then
  echo "❌ Uncommitted changes detected"
  git status -s
  exit 1
fi

# Check if on main branch
BRANCH=$(git branch --show-current)
if [[ "$BRANCH" != "main" && "$BRANCH" != "master" ]]; then
  echo "❌ Not on main/master branch (current: $BRANCH)"
  exit 1
fi

# Check if changelog has unreleased entries
if ! grep -A5 "## \[Unreleased\]" CHANGELOG.md | grep -q "^- "; then
  echo "❌ No unreleased entries in CHANGELOG.md"
  exit 1
fi

# Run tests
echo "🧪 Running tests..."
go test ./...

echo "✅ All pre-release checks passed"
echo "Ready to run: changie bump [major|minor|patch]"
PRERELEASE

echo
echo "5. JSON Output Parsing"
echo "----------------------"
cat <<'JSON_PARSING'
#!/bin/bash
# Parse JSON output for automation

RESULT=$(changie bump patch --json)

# Extract values using jq
SUCCESS=$(echo "$RESULT" | jq -r '.success')
NEW_VERSION=$(echo "$RESULT" | jq -r '.new_version')
OLD_VERSION=$(echo "$RESULT" | jq -r '.old_version')

if [[ "$SUCCESS" == "true" ]]; then
  echo "✅ Successfully released version $NEW_VERSION (was $OLD_VERSION)"

  # Use in subsequent steps
  echo "NEW_VERSION=$NEW_VERSION" >> $GITHUB_ENV

  # Send notification
  curl -X POST https://hooks.slack.com/services/YOUR/WEBHOOK/URL \
    -H 'Content-Type: application/json' \
    -d "{\"text\":\"🚀 Released $NEW_VERSION\"}"
else
  ERROR=$(echo "$RESULT" | jq -r '.error')
  echo "❌ Release failed: $ERROR"
  exit 1
fi
JSON_PARSING

echo
echo "=== Environment Variables ==="
echo "You can also configure changie via environment variables:"
echo "  export APP_VERSION_AUTO_PUSH=true"
echo "  export APP_CHANGELOG_FILE=HISTORY.md"
echo "  export APP_LOG_LEVEL=debug"
