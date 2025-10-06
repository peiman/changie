#!/usr/bin/env bash
# Basic changie workflow example
#
# This script demonstrates the typical workflow for using changie
# to manage a project's changelog and version releases.

set -e  # Exit on error

echo "=== Changie Basic Workflow Example ==="
echo

# 1. Initialize changie (only needed once per project)
echo "Step 1: Initialize changelog"
echo "  $ changie init"
echo "  Creates CHANGELOG.md following Keep a Changelog format"
echo

# 2. During development, add changelog entries as you work
echo "Step 2: Add changelog entries during development"
echo "  $ changie changelog added \"New user authentication feature\""
echo "  $ changie changelog fixed \"Memory leak in background processor\""
echo "  $ changie changelog changed \"Improved API response times\""
echo

# 3. When ready to release, bump the version
echo "Step 3: Release a new version"
echo "  For bug fixes:     changie bump patch   # 1.2.3 → 1.2.4"
echo "  For new features:  changie bump minor   # 1.2.3 → 1.3.0"
echo "  For breaking changes: changie bump major # 1.2.3 → 2.0.0"
echo

# 4. Push to remote
echo "Step 4: Push changes to remote"
echo "  $ git push && git push --tags"
echo "  Or use --auto-push: changie bump minor --auto-push"
echo

echo "=== Complete Workflow ==="
echo
cat <<'WORKFLOW'
# Real-world example
changie init
git add CHANGELOG.md
git commit -m "Initialize changelog"

# Work on features
git checkout -b feature/user-auth
# ... make changes ...
changie changelog added "OAuth2 authentication support"
git commit -am "feat: add OAuth2 authentication"

# Merge to main
git checkout main
git merge feature/user-auth

# Release new version
changie bump minor --auto-push
# Result: 1.2.3 → 1.3.0, changelog updated, git tag created, pushed to remote
WORKFLOW
