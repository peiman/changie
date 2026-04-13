#!/usr/bin/env bash
# Release Workflows
#
# All version bumping is developer-side. These are the common patterns.

set -e

echo "=== Standard Release ==="
echo
cat <<'STANDARD'
git checkout main
git pull origin main
changie bump patch --auto-push
# CI picks up the tag and publishes
STANDARD

echo
echo "=== Hotfix ==="
echo
cat <<'HOTFIX'
git checkout -b hotfix/critical-fix v1.2.3
# ... fix the issue ...
changie changelog security "Fixed SQL injection vulnerability"
git add .
git commit -m "security: fix SQL injection in user input"
changie bump patch --allow-any-branch
git push origin hotfix/critical-fix --tags
# Merge back to main
git checkout main
git merge hotfix/critical-fix
git push origin main
HOTFIX

echo
echo "=== Pre-release Review ==="
echo
cat <<'REVIEW'
#!/bin/bash
# Review unreleased changes before bumping

echo "Unreleased changes:"
sed -n '/## \[Unreleased\]/,/## \[/p' CHANGELOG.md | head -n -1

changie changelog validate

read -p "Bump type (patch/minor/major): " bump_type
changie bump "$bump_type" --auto-push
REVIEW
