# Changie Examples

Practical examples for using changie.

## Files

- **`basic-workflow.sh`** — Day-to-day workflow: init, add entries, validate, bump, push.
- **`ci-integration.sh`** — CI patterns: changelog validation on PRs, release notes extraction.
- **`release-workflow.sh`** — Release patterns: standard, hotfix, pre-release review.

## Quick Reference

```bash
# Initialize (once)
changie init

# Add entries during development
changie changelog added "New feature"
changie changelog fixed "Bug fix"

# Validate before merging
changie changelog validate

# Release (from main, clean working directory)
changie bump patch              # 1.2.3 → 1.2.4
changie bump minor              # 1.2.3 → 1.3.0
changie bump major              # 1.2.3 → 2.0.0
changie bump minor --auto-push  # bump + push in one step

# Compare versions
changie diff v1.0.0 v1.1.0
```

## How Releases Work

```
Developer: changie bump minor --auto-push
    → edits CHANGELOG.md, commits, tags, pushes
    
CI: detects new tag
    → goreleaser builds binaries, publishes GitHub release
```

Version bumping is always developer-side. CI validates and publishes.
