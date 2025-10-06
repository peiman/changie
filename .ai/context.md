# Changie - AI Agent Context

This document provides essential context for AI agents working with changie.

## What is Changie?

Changie is a CLI tool for managing changelogs following the Keep a Changelog format and Semantic Versioning principles. It automates version bumping, changelog updates, and git operations.

## Core Capabilities

### 1. Initialize Projects
```bash
changie init
```
Creates CHANGELOG.md with proper Keep a Changelog structure.

### 2. Add Changelog Entries
```bash
changie changelog <section> "entry text"
```
Sections: `added`, `changed`, `deprecated`, `removed`, `fixed`, `security`

### 3. Bump Versions
```bash
changie bump <type>
```
Types: `major` (X.0.0), `minor` (x.Y.0), `patch` (x.y.Z)

## When to Use Each Bump Type

**Major (X.0.0)** - Breaking changes:
- Removing features
- Changing API contracts
- Incompatible changes

**Minor (x.Y.0)** - New features:
- Adding new functionality
- New API endpoints
- Backward-compatible enhancements

**Patch (x.y.Z)** - Bug fixes:
- Fixing bugs
- Security patches
- Documentation updates

## Important Constraints

1. **Clean Working Directory Required**
   - All changes must be committed before bumping version
   - Use: `git add . && git commit` first

2. **Branch Restrictions (default)**
   - Only works on `main` or `master` branch
   - Override with: `--allow-any-branch` flag

3. **Git Tags Determine Version**
   - Current version read from git tags
   - No tags = starts at v0.1.0

4. **Unreleased Section Required**
   - CHANGELOG.md must have `## [Unreleased]` section
   - Must contain at least one entry

## JSON Output for Agents

Always use `--json` flag for machine-readable output:

```bash
changie bump patch --json
```

Output structure:
```json
{
  "success": true,
  "old_version": "1.2.3",
  "new_version": "1.2.4",
  "tag": "v1.2.4",
  "changelog_file": "CHANGELOG.md",
  "pushed": false,
  "bump_type": "patch"
}
```

On error:
```json
{
  "success": false,
  "error": "detailed error message with solution hints",
  "bump_type": "patch"
}
```

## Common Workflows

### Standard Release Flow
1. Ensure working directory is clean
2. Verify on main/master branch
3. Check unreleased changelog entries exist
4. Run: `changie bump <type> --json`
5. Parse JSON response for success
6. Push with: `git push && git push --tags`

### CI/CD Automation Flow
1. Install: `go install github.com/peiman/changie@latest`
2. Configure git: `git config user.name/email`
3. Run: `changie bump <type> --auto-push --json`
4. Parse version from JSON for downstream jobs

### Hotfix Flow
1. Create hotfix branch from tag
2. Make fixes and add changelog entry
3. Run: `changie bump patch --allow-any-branch --json`
4. Merge back to main

## Error Handling

Common errors and solutions:

**"uncommitted changes detected"**
- Solution: Commit or stash changes first

**"not on main or master branch"**
- Solution 1: Switch to main
- Solution 2: Use `--allow-any-branch`

**"no git tags found"**
- Solution: Run any bump (will create v0.1.0)

**"CHANGELOG.md not found"**
- Solution: Run `changie init` first

**"no unreleased entries"**
- Solution: Add entry with `changie changelog <section> "text"`

## Configuration

### Via Flags
```bash
changie bump minor --file HISTORY.md --auto-push
```

### Via Environment Variables
```bash
export APP_VERSION_AUTO_PUSH=true
export APP_CHANGELOG_FILE=CHANGELOG.md
export APP_LOG_LEVEL=debug
changie bump minor
```

### Via Config File (~/.changie.yaml)
```yaml
app:
  log_level: info
  json_output: false
  version:
    use_v_prefix: true
    auto_push: false
  changelog:
    file: CHANGELOG.md
```

## Best Practices for Agents

1. **Always use JSON output** - Reliable parsing
2. **Check success field first** - Before accessing other fields
3. **Handle errors gracefully** - Error messages include solutions
4. **Verify prerequisites** - Git installed, clean working dir, etc.
5. **Use appropriate bump type** - Based on change nature
6. **Add changelog entries first** - Before bumping version

## Integration Points

### GitHub Actions
```yaml
- run: changie bump ${{ inputs.type }} --json > release.json
- id: version
  run: echo "version=$(jq -r '.new_version' release.json)" >> $GITHUB_OUTPUT
```

### GitLab CI
```yaml
script:
  - changie bump patch --json | tee release.json
  - export VERSION=$(jq -r '.new_version' release.json)
```

### Shell Scripts
```bash
RESULT=$(changie bump patch --json)
if [[ $(echo "$RESULT" | jq -r '.success') == "true" ]]; then
  VERSION=$(echo "$RESULT" | jq -r '.new_version')
  echo "Released: $VERSION"
fi
```

## Project Specifics

- **Language**: Go 1.20+
- **CLI Framework**: Cobra
- **Config**: Viper
- **Logging**: Zerolog (JSON format available)
- **SemVer Library**: Masterminds/semver

## File Locations

- **Binary**: `changie`
- **Config**: `~/.changie.yaml` (optional)
- **Changelog**: `./CHANGELOG.md` (default, configurable)
- **Documentation**: `./llms.txt` (LLM-optimized docs)
- **Examples**: `./examples/` (usage patterns)

## Decision Trees

### Which bump type should I use?
```
Does the change break existing functionality?
├─ Yes → major
└─ No → Does it add new features?
    ├─ Yes → minor
    └─ No → patch (bug fixes, docs, etc.)
```

### What if command fails?
```
Check error message:
├─ "uncommitted changes" → git commit/stash
├─ "not on main branch" → git checkout main OR --allow-any-branch
├─ "no git tags" → First release, will create v0.1.0
├─ "no unreleased entries" → changie changelog <section> "text"
└─ Other → Check error message (includes solution hints)
```

## Version History

- v1.1.0: Restructured commands (major → bump major)
- v1.0.0: Initial stable release with ckeletin-go framework
- Earlier: Legacy versions (pre-framework)

## Support Resources

- Main docs: README.md
- LLM docs: llms.txt
- Examples: examples/
- Issues: https://github.com/peiman/changie/issues
