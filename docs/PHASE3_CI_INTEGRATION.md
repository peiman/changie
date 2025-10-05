# Phase 3: CI Integration with GoReleaser

This document contains the exact changes needed to integrate GoReleaser with the existing GitHub Actions CI workflow.

## Quick Start

When ready to implement Phase 3, follow these steps:

1. [Update GitHub Actions workflow](#1-update-github-actions-workflow)
2. [Configure GitHub Secrets](#2-configure-github-secrets)
3. [Create Homebrew Tap Repository](#3-create-homebrew-tap-repository-optional)
4. [Test the Integration](#4-test-the-integration)

---

## 1. Update GitHub Actions Workflow

### File: `.github/workflows/ci.yml`

**Current Release Job** (lines ~40-70):

```yaml
  release:
    name: Release
    runs-on: ubuntu-latest
    if: startsWith(github.ref, 'refs/tags/')
    needs: build

    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Validate tag format
        run: |
          if ! echo "${{ github.ref_name }}" | grep -qE '^v?[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$'; then
            echo "Error: Tag must follow semantic versioning (e.g., v1.0.0 or 1.0.0)"
            exit 1
          fi

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.24'
          cache: true

      - name: Install yq
        run: |
          wget https://github.com/mikefarah/yq/releases/download/v4.34.1/yq_linux_amd64 -O /usr/local/bin/yq
          chmod +x /usr/local/bin/yq

      - name: Extract BINARY_NAME from Taskfile.yml
        id: get_binary_name
        run: |
          APP_NAME=$(yq '.vars.BINARY_NAME' Taskfile.yml)
          echo "app_name=$APP_NAME" >> $GITHUB_OUTPUT

      - name: Install Task
        uses: arduino/setup-task@v2
        with:
          version: 3.x
          repo-token: ${{ secrets.GITHUB_TOKEN }}

      - name: Build
        run: task build

      - name: Create Release
        uses: softprops/action-gh-release@v2
        with:
          files: ./${{ steps.get_binary_name.outputs.app_name }}
          token: ${{ secrets.CHANGIE_GITHUB_TOKEN }}
          name: Release ${{ github.ref_name }}
          body: |
            Release ${{ github.ref_name }}
          draft: false
          prerelease: false
```

**Replace with** (GoReleaser Integration):

```yaml
  release:
    name: Release with GoReleaser
    runs-on: ubuntu-latest
    if: startsWith(github.ref, 'refs/tags/')
    needs: build

    steps:
      - name: Checkout code
        uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Required for changelog generation

      - name: Validate tag format
        run: |
          if ! echo "${{ github.ref_name }}" | grep -qE '^v?[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$'; then
            echo "Error: Tag must follow semantic versioning (e.g., v1.0.0 or 1.0.0)"
            exit 1
          fi

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

### Key Changes

1. **Removed**:
   - Manual binary building with `task build`
   - `yq` installation
   - Binary name extraction
   - `softprops/action-gh-release`

2. **Added**:
   - `goreleaser/goreleaser-action@v6`
   - `fetch-depth: 0` for full git history
   - `HOMEBREW_TAP_GITHUB_TOKEN` environment variable

3. **Kept**:
   - Tag validation
   - Go setup
   - Conditional execution on tags

---

## 2. Configure GitHub Secrets

### Required Secrets

Navigate to: `https://github.com/peiman/changie/settings/secrets/actions`

#### A. CHANGIE_GITHUB_TOKEN (Already Exists)

If it doesn't exist or needs recreation:

1. Create token at: https://github.com/settings/tokens?type=beta
2. Configure:
   - Name: `changie-releases`
   - Expiration: 90 days
   - Repository access: Only `peiman/changie`
   - Permissions:
     - ✅ Contents: Read and write
     - ✅ Metadata: Read-only
3. Copy token
4. Add to repository secrets:
   - Name: `CHANGIE_GITHUB_TOKEN`
   - Value: `github_pat_xxxxxxxxxxxxx`

#### B. HOMEBREW_TAP_GITHUB_TOKEN (New)

1. Create token at: https://github.com/settings/tokens?type=beta
2. Configure:
   - Name: `homebrew-tap-bot`
   - Expiration: 90 days
   - Repository access: Only `peiman/homebrew-tap`
   - Permissions:
     - ✅ Contents: Read and write
     - ✅ Metadata: Read-only
3. Copy token
4. Add to repository secrets:
   - Name: `HOMEBREW_TAP_GITHUB_TOKEN`
   - Value: `github_pat_xxxxxxxxxxxxx`

---

## 3. Create Homebrew Tap Repository (Optional)

### Skip if You Don't Want Homebrew Distribution

If you don't want Homebrew, update `.goreleaser.yml`:

```yaml
# Comment out or remove the brews section:
# brews:
#   - name: changie
#     ...
```

Then skip to step 4.

### Create Tap Repository

1. **Create Repository**:
   - Go to: https://github.com/new
   - Repository name: `homebrew-tap` (must be exact)
   - Visibility: Public
   - Initialize with README

2. **Add Formula Directory**:
   ```bash
   git clone git@github.com:peiman/homebrew-tap.git
   cd homebrew-tap
   mkdir -p Formula
   git add Formula/.gitkeep
   git commit -m "chore: initialize Formula directory"
   git push
   ```

3. **Update README**:
   ```markdown
   # Homebrew Tap for changie

   ## Installation

   ```bash
   brew install peiman/tap/changie
   ```

   ## Formulas

   - **changie**: Professional changelog management CLI for SemVer projects
   ```

4. **Configure Token** (see step 2B above)

---

## 4. Test the Integration

### A. Create Test Tag Locally

```bash
# Ensure working directory is clean
git status

# Create RC tag for testing
git tag -a v0.9.2-rc1 -m "Test GoReleaser CI integration"

# Don't push yet - test locally first
task release:dry-run
```

### B. Push Test Tag

```bash
git push origin v0.9.2-rc1
```

### C. Monitor CI

1. Go to: https://github.com/peiman/changie/actions
2. Watch the "Release with GoReleaser" workflow
3. Check for errors

### D. Verify Release

1. Go to: https://github.com/peiman/changie/releases/tag/v0.9.2-rc1
2. Verify:
   - ✅ Multiple platform archives (Linux, macOS, Windows)
   - ✅ Checksums file
   - ✅ DEB, RPM, APK packages
   - ✅ Release notes generated from commits

### E. Test Homebrew (if configured)

```bash
# Check tap was updated
open https://github.com/peiman/homebrew-tap/tree/main/Formula

# Should see: changie.rb formula

# Test installation
brew tap peiman/tap
brew install changie
changie --version
```

### F. Cleanup Test Release (Optional)

```bash
# Delete test tag and release
git tag -d v0.9.2-rc1
git push origin :refs/tags/v0.9.2-rc1

# Delete release on GitHub:
# https://github.com/peiman/changie/releases/tag/v0.9.2-rc1 → Delete
```

---

## 5. Rollback Plan

If Phase 3 causes issues, revert the changes:

### Revert GitHub Actions

```bash
# In .github/workflows/ci.yml, restore the old release job
git diff HEAD~1 .github/workflows/ci.yml  # See changes
git checkout HEAD~1 -- .github/workflows/ci.yml
git commit -m "revert: rollback to old release process"
git push
```

### Keep Old Job Commented

Alternatively, keep both jobs and comment out the new one:

```yaml
  # release:  # Old job (backup)
  #   name: Release
  #   ...

  release:  # New job (active)
    name: Release with GoReleaser
    ...
```

---

## Benefits After Phase 3

### Before (Current)
- ✅ Single binary for one platform
- ✅ GitHub release created
- ❌ Manual multi-platform builds
- ❌ No package distribution
- ❌ No Homebrew support

### After (Phase 3)
- ✅ Multi-platform binaries (20+ platforms/architectures)
- ✅ Automated GitHub releases
- ✅ Homebrew tap updated automatically
- ✅ DEB/RPM/APK packages
- ✅ Checksums for verification
- ✅ Better release notes from commits
- ✅ Archive files (tar.gz, zip)

---

## Estimated Time to Implement

- Update CI workflow: **5 minutes**
- Configure GitHub secrets: **10 minutes**
- Create Homebrew tap: **15 minutes**
- Test and verify: **20 minutes**

**Total: ~50 minutes**

---

## Questions & Troubleshooting

### Q: Will this break existing releases?

No. Old releases remain unchanged. New releases will use GoReleaser.

### Q: Can I still use `task release` locally?

Yes! Local tasks (Phase 1 & 2) are independent of CI.

### Q: What if Homebrew fails?

Set `HOMEBREW_TAP_GITHUB_TOKEN` in CI. If still failing, comment out the `brews` section in `.goreleaser.yml`.

### Q: How do I know it's working?

Check GitHub Actions logs. Look for:
```
• building binaries
  • building     binary=dist/changie_linux_amd64/changie
  ...
• archives
  • archiving    name=dist/changie_1.0.0_Linux_x86_64.tar.gz
  ...
✅ Release published successfully
```

---

## Checklist

Before implementing Phase 3:

- [ ] Phase 1 & 2 tasks work locally (`task release:snapshot`)
- [ ] `.goreleaser.yml` validates (`task release:check`)
- [ ] `CHANGIE_GITHUB_TOKEN` secret exists in GitHub
- [ ] Homebrew tap repository created (or brews section disabled)
- [ ] `HOMEBREW_TAP_GITHUB_TOKEN` secret created (if using Homebrew)
- [ ] Test tag ready (e.g., `v0.9.2-rc1`)
- [ ] Backup of current CI workflow saved

After implementing:

- [ ] CI workflow updated
- [ ] Test tag pushed
- [ ] Release created successfully
- [ ] Multiple platform binaries present
- [ ] Homebrew formula updated (if enabled)
- [ ] Documentation updated

---

## Next Steps

1. ✅ Read this document
2. ⏸️ When ready, follow steps 1-4
3. ⏸️ Create first real release with GoReleaser
4. ⏸️ Update main README.md with installation instructions
5. ⏸️ Announce multi-platform support

---

**Implementation Status**: 📋 Planned (Ready to Execute)

**Last Updated**: 2025-10-06
