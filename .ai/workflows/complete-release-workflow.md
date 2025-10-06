# Workflow: Complete Release Process

This workflow guides you through the complete process of releasing a new version using changie.

## Overview

This multi-step workflow takes you from development through release and deployment.

## Prerequisites Check

Before starting, verify:
- [ ] Git is installed and repository initialized
- [ ] Working directory is clean (no uncommitted changes)
- [ ] On main or master branch (or have --allow-any-branch ready)
- [ ] CHANGELOG.md exists and has unreleased entries
- [ ] All tests passing

## Workflow Steps

### Step 1: Review Unreleased Changes

**Purpose:** Understand what's being released

```bash
# View unreleased section
sed -n '/## \[Unreleased\]/,/## \[/p' CHANGELOG.md | head -n -1
```

**Review:**
- Are all changes documented?
- Are entries in correct sections?
- Is anything missing?

### Step 2: Determine Version Bump Type

**Decision Matrix:**

| Changes Include | Bump Type |
|----------------|-----------|
| Breaking changes, removed features | `major` |
| New features, enhancements | `minor` |
| Only bug fixes, patches | `patch` |

**Analysis Questions:**
1. Does this break existing API/functionality? → major
2. Does this add new capabilities? → minor
3. Does this only fix bugs? → patch

### Step 3: Pre-Release Validation

**Run checks:**

```bash
# Check git status
git status

# Verify tests pass
go test ./...  # or your test command

# Verify current version
git describe --tags --abbrev=0
```

**All green?** → Proceed
**Any failures?** → Fix first

### Step 4: Execute Version Bump

**Command:**
```bash
changie bump <type> --json > release-result.json
```

**Parse result:**
```bash
SUCCESS=$(jq -r '.success' release-result.json)
NEW_VERSION=$(jq -r '.new_version' release-result.json)
```

**Verify:**
- `SUCCESS` is `true`
- `NEW_VERSION` looks correct

### Step 5: Review Git Changes

**Check what was changed:**

```bash
# View last commit
git log -1 --stat

# View new tag
git tag -l | tail -1

# View changelog update
git diff HEAD~1 CHANGELOG.md
```

**Verify:**
- Unreleased section moved to new version
- Commit message is "Release vX.Y.Z"
- Tag created with correct version

### Step 6: Push to Remote

**Option A: Manual push (more control)**
```bash
git push origin main
git push origin --tags
```

**Option B: Auto-push (in Step 4)**
```bash
changie bump <type> --auto-push --json > release-result.json
```

### Step 7: Post-Release Actions

**Typical follow-ups:**

1. **Create GitHub/GitLab Release**
   ```bash
   gh release create $NEW_VERSION --title "Release $NEW_VERSION" --notes-from-tag
   ```

2. **Trigger Deployment**
   ```bash
   # Trigger CI/CD pipeline
   gh workflow run deploy.yml -f version=$NEW_VERSION
   ```

3. **Update Documentation**
   ```bash
   # Update docs site with new version
   curl -X POST $DOCS_WEBHOOK -d "{\"version\":\"$NEW_VERSION\"}"
   ```

4. **Notify Team**
   ```bash
   # Slack notification
   curl -X POST $SLACK_WEBHOOK \
     -d "{\"text\":\"🚀 Released $NEW_VERSION\"}"
   ```

5. **Update Package Managers**
   - npm: Update package.json version
   - Docker: Build and push new image
   - etc.

### Step 8: Verify Release

**Checks:**

- [ ] GitHub/GitLab shows new tag
- [ ] Release notes published
- [ ] CI/CD pipeline triggered
- [ ] Deployment successful
- [ ] New version accessible to users

## Error Recovery

### If Step 4 Fails

**Uncommitted changes:**
```bash
git add .
git commit -m "chore: prepare for release"
# Retry step 4
```

**Wrong branch:**
```bash
git checkout main
# Retry step 4
```

**No unreleased entries:**
```bash
changie changelog <section> "Entry text"
git add CHANGELOG.md
git commit -m "docs: add missing changelog entry"
# Retry step 4
```

### If Push Fails (Step 6)

**Network issues:**
```bash
# Retry push
git push origin main --tags
```

**Conflicts:**
```bash
# Pull first
git pull origin main --rebase
git push origin main --tags
```

## Rollback Procedure

**If you need to undo a release:**

```bash
# Delete tag locally
git tag -d v1.2.3

# Delete tag remotely
git push origin :refs/tags/v1.2.3

# Revert commit (if not already pushed elsewhere)
git reset --hard HEAD~1

# Or create revert commit (if pushed)
git revert HEAD
```

**⚠️ Warning:** Only rollback if version hasn't been deployed/used

## Complete Example

```bash
# 1. Review changes
sed -n '/## \[Unreleased\]/,/## \[/p' CHANGELOG.md

# 2. Decide: has new features → minor bump

# 3. Pre-check
git status  # clean
go test ./...  # pass

# 4. Bump version
changie bump minor --json > release.json

# 5. Verify
jq . release.json
# {
#   "success": true,
#   "new_version": "v1.3.0",
#   ...
# }

# 6. Push
git push origin main --tags

# 7. Create release
VERSION=$(jq -r '.new_version' release.json)
gh release create $VERSION --generate-notes

# 8. Notify
echo "✅ Released $VERSION"
```

## Automation Script

Save this as `release.sh`:

```bash
#!/bin/bash
set -e

# Parse bump type
BUMP_TYPE=${1:-patch}

# Pre-checks
echo "🔍 Running pre-release checks..."
git diff --exit-code || { echo "❌ Uncommitted changes"; exit 1; }
go test ./... || { echo "❌ Tests failed"; exit 1; }

# Release
echo "🚀 Releasing $BUMP_TYPE version..."
changie bump $BUMP_TYPE --json > release.json

# Parse result
SUCCESS=$(jq -r '.success' release.json)
VERSION=$(jq -r '.new_version' release.json)

if [[ "$SUCCESS" != "true" ]]; then
  echo "❌ Release failed"
  jq -r '.error' release.json
  exit 1
fi

# Push
echo "📤 Pushing to remote..."
git push origin main --tags

# Create GitHub release
echo "📝 Creating GitHub release..."
gh release create $VERSION --generate-notes

echo "✅ Released $VERSION successfully!"
```

## Usage in CI/CD

See `examples/ci-integration.sh` for CI/CD pipeline examples.

## Next Steps

After successful release:
- Monitor deployment
- Watch for issues
- Be ready for hotfix if needed
- Start next development cycle
