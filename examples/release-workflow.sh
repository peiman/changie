#!/usr/bin/env bash
# Release Workflow Examples
#
# Different release strategies and workflows using changie

set -e

echo "=== Release Workflow Strategies ==="
echo

echo "1. Simple Main Branch Release"
echo "------------------------------"
cat <<'SIMPLE'
# Ensure you're on main with latest changes
git checkout main
git pull origin main

# Add any last-minute changelog entries
changie changelog fixed "Critical bug in parser"

# Commit changes
git add CHANGELOG.md
git commit -m "docs: add final changelog entries for release"

# Bump version and release
changie bump patch --auto-push

# Done! Version tagged and pushed
SIMPLE

echo
echo "2. Release Branch Workflow"
echo "--------------------------"
cat <<'RELEASE_BRANCH'
# Create release branch from main
git checkout main
git pull origin main
git checkout -b release/v2.0.0

# Make any release-specific changes
# Fix bugs, update docs, etc.
changie changelog fixed "Issue found during release testing"

# Commit all release changes
git add .
git commit -m "chore: prepare v2.0.0 release"

# Bump version on release branch
changie bump major --allow-any-branch

# Push release branch
git push origin release/v2.0.0

# Merge back to main
git checkout main
git merge release/v2.0.0
git push origin main
git push origin --tags

# Clean up release branch
git branch -d release/v2.0.0
git push origin --delete release/v2.0.0
RELEASE_BRANCH

echo
echo "3. Hotfix Workflow"
echo "------------------"
cat <<'HOTFIX'
# Create hotfix branch from latest tag
git checkout -b hotfix/critical-security-fix v1.2.3

# Fix the issue
# ... make changes ...

# Add changelog entry
changie changelog security "Fixed SQL injection vulnerability"

# Commit the fix
git add .
git commit -m "security: fix SQL injection in user input"

# Bump patch version
changie bump patch --allow-any-branch

# Push hotfix
git push origin hotfix/critical-security-fix
git push origin --tags

# Merge to main
git checkout main
git merge hotfix/critical-security-fix
git push origin main

# Merge to develop if you use git-flow
git checkout develop
git merge hotfix/critical-security-fix
git push origin develop

# Clean up
git branch -d hotfix/critical-security-fix
git push origin --delete hotfix/critical-security-fix
HOTFIX

echo
echo "4. Pre-release Workflow"
echo "-----------------------"
cat <<'PRERELEASE'
# For alpha/beta releases, manual tagging might be preferred
# Use changie for changelog, but manually create pre-release tags

# Add changelog entries under unreleased
changie changelog added "Experimental new API (beta)"

# Review changelog
cat CHANGELOG.md

# Create pre-release tag manually
git tag -a v2.0.0-beta.1 -m "Release v2.0.0-beta.1"
git push origin v2.0.0-beta.1

# When ready for actual release
changie bump major --auto-push
PRERELEASE

echo
echo "5. Monorepo Release Strategy"
echo "----------------------------"
cat <<'MONOREPO'
# For monorepos with multiple packages
# Use different changelog files for each package

# Package A release
changie bump minor --file packages/api/CHANGELOG.md

# Package B release
changie bump patch --file packages/web/CHANGELOG.md

# Or use a unified changelog with sections
changie changelog added "[API] New authentication endpoint"
changie changelog changed "[Web] Updated dashboard UI"

# Single version bump for the entire monorepo
changie bump minor --auto-push
MONOREPO

echo
echo "6. Dry Run Testing"
echo "------------------"
cat <<'DRYRUN'
#!/bin/bash
# Test release process without actually releasing

# Create a test branch
git checkout -b test-release

# Try the bump (without --auto-push)
changie bump minor

# Review changes
git log -1
git show HEAD
git tag -l

# If satisfied, reset and do it for real on main
git checkout main
git branch -D test-release

# Reset tags if created during test
git tag -d v1.3.0  # Replace with your test version

# Do the real release
changie bump minor --auto-push
DRYRUN

echo
echo "7. Rollback Strategy"
echo "--------------------"
cat <<'ROLLBACK'
#!/bin/bash
# If you need to undo a release

# Find the release commit and tag
git log --oneline -5
git tag -l

# Delete the tag locally and remotely
git tag -d v1.2.4
git push origin :refs/tags/v1.2.4

# Revert the changelog commit
git revert HEAD  # Or use git reset --hard HEAD~1 if not pushed

# If already pushed to main
git push origin main

# Note: In most cases, it's better to create a new patch release
# rather than rolling back a published version
ROLLBACK

echo
echo "8. Changelog Review Before Release"
echo "-----------------------------------"
cat <<'REVIEW'
#!/bin/bash
# Script to review unreleased changes before bumping

echo "📋 Unreleased changes:"
echo

# Extract unreleased section from changelog
sed -n '/## \[Unreleased\]/,/## \[/p' CHANGELOG.md | head -n -1

echo
read -p "Review looks good? (y/n) " -n 1 -r
echo

if [[ $REPLY =~ ^[Yy]$ ]]; then
  echo "What type of release?"
  echo "1) patch (bug fixes)"
  echo "2) minor (new features)"
  echo "3) major (breaking changes)"
  read -p "Enter choice (1-3): " choice

  case $choice in
    1) changie bump patch --auto-push ;;
    2) changie bump minor --auto-push ;;
    3) changie bump major --auto-push ;;
    *) echo "Invalid choice"; exit 1 ;;
  esac
else
  echo "Cancelled. Please update changelog and try again."
  exit 1
fi
REVIEW
