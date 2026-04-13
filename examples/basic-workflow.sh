#!/usr/bin/env bash
# Basic changie workflow
#
# The typical developer workflow for managing changelogs and releases.

set -e

echo "=== Changie Workflow ==="
echo

# 1. Initialize (once per project)
echo "Step 1: Initialize changelog"
echo "  $ changie init"
echo

# 2. Add entries as you develop
echo "Step 2: Add changelog entries during development"
echo "  $ changie changelog added \"New user authentication feature\""
echo "  $ changie changelog fixed \"Memory leak in background processor\""
echo "  $ changie changelog changed \"Improved API response times\""
echo "  $ changie changelog security \"Updated dependencies with CVE fixes\""
echo

# 3. Validate before merging
echo "Step 3: Validate changelog in PR"
echo "  $ changie changelog validate"
echo

# 4. Release from main
echo "Step 4: Bump version (from main branch, clean working directory)"
echo "  For bug fixes:        changie bump patch   # 1.2.3 → 1.2.4"
echo "  For new features:     changie bump minor   # 1.2.3 → 1.3.0"
echo "  For breaking changes: changie bump major   # 1.2.3 → 2.0.0"
echo

# 5. Push
echo "Step 5: Push to trigger CI release"
echo "  $ git push && git push --tags"
echo "  Or: changie bump minor --auto-push"
echo

echo "=== Complete Example ==="
echo
cat <<'WORKFLOW'
# Feature development
git checkout -b feature/user-auth
# ... make changes ...
changie changelog added "OAuth2 authentication support"
git add .
git commit -m "feat: add OAuth2 authentication"

# Before merging PR
changie changelog validate

# Merge to main, then release
git checkout main
git merge feature/user-auth
changie bump minor --auto-push
# → Updates CHANGELOG.md, commits, tags v1.3.0, pushes
# → CI detects new tag, runs goreleaser, publishes release
WORKFLOW
