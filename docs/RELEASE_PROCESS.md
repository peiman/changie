# Release Process Documentation

This document describes the release process for changie, including local testing, manual releases, and automated CI/CD releases.

## Table of Contents

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Phase 1 & 2: Local Release Tasks (Implemented)](#phase-1--2-local-release-tasks-implemented)
- [Phase 3: CI Integration (Planned)](#phase-3-ci-integration-planned)
- [Package Distribution](#package-distribution)
- [Troubleshooting](#troubleshooting)

---

## Overview

The changie project uses [GoReleaser](https://goreleaser.com) for multi-platform releases. The release process is divided into three phases:

1. **Phase 1**: Local testing tasks (snapshot builds, validation)
2. **Phase 2**: Manual release tasks (tagging, publishing)
3. **Phase 3**: Automated CI/CD integration (GitHub Actions)

**Current Status**: ✅ Phase 1 & 2 complete | ⏸️ Phase 3 planned

---

## Prerequisites

### 1. Install GoReleaser

```bash
# Install via task
task release:setup

# Or manually (macOS):
brew install goreleaser/tap/goreleaser

# Or manually (other platforms):
go install github.com/goreleaser/goreleaser/v2@latest
```

### 2. GitHub Personal Access Token

Create a fine-grained personal access token:

1. Go to: [GitHub Settings → Developer Settings](https://github.com/settings/tokens?type=beta)
2. Click "Generate new token" → **Fine-grained tokens**
3. Configure:
   - **Token name**: `changie-releases`
   - **Expiration**: 90 days (or custom)
   - **Resource owner**: peiman
   - **Repository access**: Only select repositories → `changie`
   - **Permissions**:
     - ✅ Contents: Read and write
     - ✅ Metadata: Read-only (auto-selected)
4. Generate and copy the token

**For local releases:**
```bash
export GITHUB_TOKEN="github_pat_xxxxxxxxxxxxx"
```

**For CI (add to GitHub repository secrets):**
- Name: `CHANGIE_GITHUB_TOKEN`
- Value: `github_pat_xxxxxxxxxxxxx`

### 3. Homebrew Tap (Optional)

For Homebrew distribution:

1. Create repository: `https://github.com/peiman/homebrew-tap`
2. Add `Formula/` directory
3. Create a token (same steps as above) named `homebrew-tap-bot`
4. Add to GitHub secrets as `HOMEBREW_TAP_GITHUB_TOKEN`

---

## Phase 1 & 2: Local Release Tasks (Implemented)

### Phase 1: Testing & Validation

#### `task release:setup`
Install goreleaser (cross-platform).

```bash
task release:setup
```

#### `task release:check`
Validate `.goreleaser.yml` configuration.

```bash
task release:check
```

#### `task release:snapshot`
Build a local snapshot release (no publishing).

```bash
task release:snapshot
```

**Output**:
- Binaries in `./dist/changie_<os>_<arch>/`
- Archives: `./dist/changie_<version>_<OS>_<ARCH>.tar.gz`
- Packages: `.deb`, `.rpm`, `.apk` files

**Platforms built**:
- Linux: amd64, 386, arm64, arm (v6, v7)
- macOS: amd64, arm64 (Universal binary)
- Windows: amd64, 386, arm64
- FreeBSD, OpenBSD, NetBSD: various architectures

#### `task release:test`
Test the snapshot build binary for your platform.

```bash
task release:test
```

Runs: `./dist/changie_<os>_<arch>/changie --version`

#### `task release:dry-run`
Perform a complete release dry run with all validations.

```bash
# Without GITHUB_TOKEN (local-only)
task release:dry-run

# With GITHUB_TOKEN (full validation)
export GITHUB_TOKEN="..."
task release:dry-run
```

Includes: `task check` (format, lint, test, vuln scan)

#### `task release:clean`
Clean release artifacts.

```bash
task release:clean
```

Removes: `./dist/` directory

---

### Phase 2: Manual Releases

#### `task release:tag`
Create and push a version tag.

```bash
task release:tag TAG=v1.0.0
```

**Validations**:
- ❌ Fails if TAG not provided
- ❌ Fails if tag already exists
- ❌ Fails if working directory is dirty

**Actions**:
1. Creates annotated tag locally
2. Pushes tag to origin
3. Prints GitHub Actions tracking URL

#### `task release`
Build and publish a full release (requires GITHUB_TOKEN).

```bash
export GITHUB_TOKEN="github_pat_xxxxxxxxxxxxx"
task release
```

**Validations**:
- ❌ Fails if GITHUB_TOKEN not set
- ❌ Fails if working directory is dirty
- ❌ Fails if current commit is not tagged

**Actions**:
1. Runs `goreleaser release --clean`
2. Publishes to GitHub Releases
3. Creates all distribution packages
4. Prints release URL

---

## Recommended Workflow

### Using changie (Preferred)

The recommended workflow uses changie's built-in version bumping:

```bash
# 1. Bump version and update changelog
changie minor              # or: major, patch
# This creates: CHANGELOG.md update, git commit, git tag

# 2. (Optional) Test the release locally
task release:dry-run

# 3. Push to trigger CI release
git push --follow-tags

# GitHub Actions automatically creates the release
```

### Manual Workflow

If you prefer manual control:

```bash
# 1. Ensure working directory is clean
git status

# 2. Test locally
task release:snapshot
task release:test

# 3. Create and push tag
task release:tag TAG=v1.0.0

# 4. (Optional) Manual release instead of CI
export GITHUB_TOKEN="..."
task release
```

---

## Phase 3: CI Integration (Planned)

### Current CI Behavior

**File**: `.github/workflows/ci.yml`

**Current Release Job**:
- Triggers on: tag push
- Builds: Single binary via `task build`
- Publishes: GitHub release with one binary

**Limitations**:
- ❌ Only one platform (GitHub runner)
- ❌ No multi-platform builds
- ❌ No package distribution
- ❌ GoReleaser config unused

### Proposed Enhancement

Replace the current release job with GoReleaser integration.

#### Updated Release Job

```yaml
release:
  name: Release
  runs-on: ubuntu-latest
  if: startsWith(github.ref, 'refs/tags/')
  needs: build

  steps:
    - name: Checkout
      uses: actions/checkout@v4
      with:
        fetch-depth: 0  # Full history for changelog

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.24'
        cache: true

    - name: Run GoReleaser
      uses: goreleaser/goreleaser-action@v6
      with:
        distribution: goreleaser
        version: '~> v2'
        args: release --clean
      env:
        GITHUB_TOKEN: ${{ secrets.CHANGIE_GITHUB_TOKEN }}
        HOMEBREW_TAP_GITHUB_TOKEN: ${{ secrets.HOMEBREW_TAP_GITHUB_TOKEN }}
```

#### Implementation Steps

1. **Update `.github/workflows/ci.yml`**:
   - Replace current release job with GoReleaser action
   - Add required secrets to environment

2. **Configure GitHub Secrets**:
   - `CHANGIE_GITHUB_TOKEN` (already exists)
   - `HOMEBREW_TAP_GITHUB_TOKEN` (new - for brew tap)
   - `AUR_PRIVATE_KEY` (optional - for Arch packages)

3. **Create Homebrew Tap Repository**:
   - Repository: `github.com/peiman/homebrew-tap`
   - Directory: `Formula/`
   - README with installation instructions

4. **Update `.goreleaser.yml` (if needed)**:
   - Optionally change `brews` to `homebrews` (current config works, just has deprecation warning)
   - Disable snapcraft if not needed: `publish: false` (already set)

5. **Test CI Release**:
   ```bash
   # Create test tag
   git tag -a v0.9.2-rc1 -m "Test release"
   git push origin v0.9.2-rc1

   # Monitor at: https://github.com/peiman/changie/actions
   ```

6. **Rollback Plan**:
   - Keep old release job commented out in CI file
   - Can revert if issues occur

---

## Package Distribution

### Automatic (via GoReleaser)

**GitHub Releases**:
- All platforms and architectures
- Archives (tar.gz, zip)
- Checksums (SHA256)
- Installation: Download from [Releases](https://github.com/peiman/changie/releases)

**Homebrew** (when tap is configured):
```bash
brew install peiman/tap/changie
```

**Linux Packages**:

DEB (Ubuntu/Debian):
```bash
# Download .deb from releases
sudo dpkg -i changie_*_linux_amd64.deb
```

RPM (RedHat/Fedora/CentOS):
```bash
# Download .rpm from releases
sudo rpm -i changie_*_linux_amd64.rpm
```

APK (Alpine):
```bash
# Download .apk from releases
apk add --allow-untrusted changie_*_linux_amd64.apk
```

### Manual Installation

**Go Install**:
```bash
go install github.com/peiman/changie@latest
```

**Direct Binary Download**:
```bash
# Download from releases
curl -LO https://github.com/peiman/changie/releases/download/v1.0.0/changie_1.0.0_Linux_x86_64.tar.gz

# Extract
tar -xzf changie_1.0.0_Linux_x86_64.tar.gz

# Move to PATH
sudo mv changie /usr/local/bin/
```

---

## Troubleshooting

### Snapshot Build Fails

**Error**: `snapcraft not present in $PATH`

**Solution**: This is expected. Snapcraft is optional.

To disable snapcraft completely:
```yaml
# In .goreleaser.yml, set:
snapcrafts:
  - publish: false  # Already set
```

Or install snapcraft:
```bash
brew install snapcraft  # macOS
sudo snap install snapcraft --classic  # Linux
```

### Universal Binary Skipped

**Warning**: `no darwin binaries found with ids: changie-universal`

**Cause**: The universal binary config references a build ID that doesn't match.

**Solution**: Update `.goreleaser.yml`:
```yaml
builds:
  - id: changie  # Make sure this matches

universal_binaries:
  - id: changie-universal
    ids:
      - changie  # Reference the build ID
```

### Homebrew Formula Fails

**Error**: Failed to push to homebrew-tap

**Checklist**:
- ✅ Repository `peiman/homebrew-tap` exists
- ✅ `HOMEBREW_TAP_GITHUB_TOKEN` is set with `repo` scope
- ✅ Token has write access to the tap repository

### GitHub Release Fails (403 Forbidden)

**Error**: `creating release: POST 403`

**Cause**: Token lacks permissions

**Solution**:
1. Recreate token with `Contents: Read and write` permission
2. For classic tokens: ensure `repo` scope is checked
3. Update secret in GitHub repository settings

### Git State Dirty

**Error**: `git is in a dirty state`

**Cause**: Uncommitted changes

**Solution**:
```bash
git status
git add .
git commit -m "chore: prepare release"
# Or: git stash
```

---

## Release Checklist

Before releasing:

- [ ] All tests pass: `task test`
- [ ] Linters pass: `task lint`
- [ ] Dependencies verified: `task deps:check`
- [ ] Vulnerability scan clean: `task vuln`
- [ ] Working directory clean: `git status`
- [ ] CHANGELOG.md updated
- [ ] Version bumped (via `changie` or manually)
- [ ] Local snapshot builds successfully: `task release:snapshot`
- [ ] Binary tested: `task release:test`

---

## Additional Resources

- **GoReleaser Docs**: https://goreleaser.com/
- **GitHub Actions**: https://github.com/peiman/changie/actions
- **Releases**: https://github.com/peiman/changie/releases
- **Homebrew Tap**: https://github.com/peiman/homebrew-tap (when created)
- **Keep a Changelog**: https://keepachangelog.com/
- **Semantic Versioning**: https://semver.org/

---

**Last Updated**: 2025-10-06
**Implemented**: Phase 1 & 2
**Next**: Phase 3 (CI Integration)
