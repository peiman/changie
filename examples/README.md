# Changie Examples

This directory contains practical examples and scripts demonstrating how to use changie in various scenarios.

## Quick Links

- [Basic Workflow](#basic-workflow) - Simple day-to-day usage
- [CI/CD Integration](#cicd-integration) - Automation in pipelines
- [Release Workflows](#release-workflows) - Different release strategies

## Files

### `basic-workflow.sh`

Demonstrates the fundamental changie workflow:
- Initializing a project
- Adding changelog entries during development
- Releasing versions
- Complete real-world example

**Best for:** New users learning changie basics

### `ci-integration.sh`

Shows how to integrate changie into CI/CD pipelines:
- GitHub Actions workflows
- GitLab CI configuration
- JSON output parsing for automation
- Pre-release validation scripts
- Automated release based on changelog content

**Best for:** DevOps engineers setting up automated releases

### `release-workflow.sh`

Covers various release strategies:
- Simple main branch releases
- Release branch workflow
- Hotfix workflow
- Pre-release (alpha/beta) handling
- Monorepo strategies
- Dry run testing
- Rollback procedures

**Best for:** Teams establishing release processes

## Usage

All scripts are educational and contain embedded examples. They're not meant to be executed directly but rather to demonstrate patterns and commands.

### Running Examples Interactively

```bash
# View an example
cat examples/basic-workflow.sh

# Or make it executable and run it (displays instructions)
chmod +x examples/basic-workflow.sh
./examples/basic-workflow.sh
```

## Common Patterns

### Initialize a New Project

```bash
changie init
git add CHANGELOG.md
git commit -m "Initialize changelog"
```

### Add Changelog Entries

```bash
# During development
changie changelog added "New user authentication feature"
changie changelog fixed "Memory leak in background processor"
changie changelog security "Updated dependencies with security fixes"
```

### Release a Version

```bash
# Bug fixes: 1.2.3 → 1.2.4
changie bump patch

# New features: 1.2.3 → 1.3.0
changie bump minor

# Breaking changes: 1.2.3 → 2.0.0
changie bump major

# With automatic push
changie bump minor --auto-push
```

### Using JSON Output

```bash
# Get machine-readable output
changie bump patch --json > release.json

# Parse with jq
VERSION=$(changie bump patch --json | jq -r '.new_version')
echo "Released version: $VERSION"
```

## Environment Configuration

Configure changie behavior via environment variables:

```bash
# Auto-push by default
export APP_VERSION_AUTO_PUSH=true

# Custom changelog file
export APP_CHANGELOG_FILE=HISTORY.md

# Verbose logging
export APP_LOG_LEVEL=debug

# JSON logs
export APP_LOG_FORMAT=json
```

## Integration Tips

### For AI Agents / LLMs

When using changie programmatically:

1. **Use `--json` flag** for machine-readable output
2. **Check `.success` field** in JSON response
3. **Parse `.new_version`** to get the released version
4. **Read `.error`** field if success is false

Example:
```bash
RESULT=$(changie bump patch --json)
if [[ $(echo "$RESULT" | jq -r '.success') == "true" ]]; then
  VERSION=$(echo "$RESULT" | jq -r '.new_version')
  echo "Success: $VERSION"
else
  ERROR=$(echo "$RESULT" | jq -r '.error')
  echo "Failed: $ERROR"
fi
```

### For CI/CD Pipelines

1. **Install changie** in your pipeline:
   ```bash
   go install github.com/peiman/changie@latest
   ```

2. **Configure git** for commits:
   ```bash
   git config user.name "CI Bot"
   git config user.email "ci@example.com"
   ```

3. **Use `--auto-push`** to automatically push changes

4. **Parse JSON output** for subsequent steps

### Error Handling

Changie returns non-zero exit codes on failure. Always check:

```bash
if changie bump patch; then
  echo "Release successful"
else
  echo "Release failed"
  exit 1
fi
```

## Best Practices

1. **Commit changes before bumping** - Changie requires a clean working directory
2. **Review unreleased section** before releasing
3. **Use semantic versioning correctly**:
   - Patch: Bug fixes only
   - Minor: New features, backward compatible
   - Major: Breaking changes
4. **Add changelog entries as you work** - Don't wait until release time
5. **Use `--auto-push` in CI** - But review manually in development

## Troubleshooting

### "uncommitted changes detected"

```bash
# Option 1: Commit your changes
git add .
git commit -m "your message"

# Option 2: Stash changes temporarily
git stash
changie bump patch
git stash pop
```

### "not on main or master branch"

```bash
# Option 1: Switch to main
git checkout main

# Option 2: Use --allow-any-branch
changie bump patch --allow-any-branch
```

### "no git tags found"

```bash
# First release will start from 0.1.0 by default
changie bump patch
# Result: v0.1.0
```

## Further Reading

- [Keep a Changelog](https://keepachangelog.com/) - Changelog format specification
- [Semantic Versioning](https://semver.org/) - Versioning specification
- [Main Documentation](../README.md) - Full changie documentation
- [llms.txt](../llms.txt) - LLM-optimized documentation

## Contributing Examples

Have a useful workflow or integration? Please contribute!

1. Create a new `.sh` file with clear comments
2. Include practical, real-world examples
3. Document any prerequisites
4. Submit a pull request

Example template:

```bash
#!/usr/bin/env bash
# Description of what this example demonstrates
#
# Prerequisites:
# - Requirement 1
# - Requirement 2

set -e

echo "=== Example Name ==="
echo

# Your example code here
```
