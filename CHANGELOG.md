# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com),
and this project adheres to [Semantic Versioning (SemVer)](https://semver.org).

## [Unreleased]

### Added

- llms.txt file for LLM-optimized documentation and MCP integration
- JSON output support with --json flag for machine-readable command results
- examples/ directory with comprehensive usage scripts for basic workflows, CI/CD integration, and release strategies
- .ai/ directory with AI agent resources, prompts, workflows, and MCP integration guidance
- feat: add MCP server for AI agent integration with official Go SDK v1.0.0
- chore: add Task commands for MCP Docker operations (build, run, test, clean)

### Fixed

- Table of contents in README.md now correctly references 'bump' command
- fix: update Dockerfile.mcp to use golang:1.24-alpine for Go 1.24.0 compatibility

### Changed

- Enhanced help text for bump commands with detailed examples, use cases, and step-by-step explanations


## [v1.1.0] - 2025-10-06

### Changed

- Restructured version bump commands under 'bump' parent command (changie major → changie bump major)


## [v1.0.2] - 2025-10-06

### Fixed

- Changelog comparison links now include previous version (e.g., [1.0.0]: .../v0.9.0...v1.0.0 instead of [1.0.0]: .../...v1.0.0)


### Changed

- Refactored AddChangelogSection to reduce cyclomatic complexity from 19 to below 15


## [v1.0.1] - 2025-10-06

### Fixed

- Commit message format changed from 'Update changelog for version X.Y.Z' to 'Release X.Y.Z'

## [v1.0.0] - 2025-10-06

### Added

- Ported to [ckeletin-go](https://github.com/peiman/ckeletin-go) framework for better structure and testability
- Branch protection with `--allow-any-branch` flag for version bump commands
- Interactive version prefix prompt when no git tags exist
- `--use-v-prefix` flag for explicit control over version tag format (v1.0.0 vs 1.0.0)
- `--log-format` flag (auto, json, console) following zerolog best practices
- `--log-caller` flag for debugging with file:line information
- Component-specific sub-loggers (logger.Version, logger.Changelog, logger.Git, etc.)
- golangci-lint configuration with 30+ linters
- CLAUDE.md architecture guide for contributors
- Automatic repository detection from git remote for changelog links
- Comprehensive test coverage across all packages (325 tests total)
- GoReleaser configuration for automated multi-platform releases
- Release task commands (setup, check, snapshot, test, dry-run, clean, tag, release)
- GitHub Actions CI integration with GoReleaser
- Automated Homebrew tap updates
- Linux package generation (DEB, RPM, APK)
- .env file support in Taskfile

### Changed

- Migrated to `blang/semver/v4` library for semantic versioning
- Refactored business logic from cmd/ to internal/ packages following Go best practices
- Init command auto-detects and respects existing git tag conventions
- Logger uses TTY detection for automatic JSON/console format selection
- Updated golang.org/x/text to v0.29.0
- Test coverage: overall 80.1%, cmd 82.3%, logger 97.0%, internal/version 74.6%

### Fixed

- Init command now creates v0.0.0 tag when no tags exist (regression from v0.3.0)
- Version prefix handling now uses user preference consistently
- Semver bump operations clear prerelease and build metadata per SemVer specification
- Config key for repository_provider now uses correct `app.changelog.repository_provider`
- Changelog comparison links respect user's v-prefix preference

## [0.9.1] - 2024-07-01

### Added

- Implemented dynamic version determination based on git tags and commits

### Changed

- Enhanced testing infrastructure, particularly for git command mocking and output capturing

## [0.9.0] - 2024-06-30

### Changed

- Updated CI/CD pipeline configuration

## [0.8.0] - 2024-06-28

### Fixed

- Fixed --auto-push flag

## [0.7.0] - 2024-06-28

### Added

- Added --auto-push flag to automatically push changes and tags after version bump

## [0.6.0] - 2024-06-28

### Added

- Reject version bumps when there are uncommitted git changes

## [0.5.0] - 2024-06-27

### Added

- Reject version bumps when there are uncommitted git changes

## [0.4.1] - 2024-06-27

### Fixed

- Comparison links fixed in this changelog

## [0.4.0] - 2024-06-27

### Added

- GetLatestChangelogVersion function to extract version from changelog content
- Version mismatch checking between git tags and changelog
- Test mode flag to suppress warnings during tests
- More comprehensive tests for version bumping and changelog updates

### Changed

- Improved UpdateChangelog function to handle existing content more robustly
- Enhanced updateDiffLinks function to correctly maintain comparison links
- Updated handleVersionBump to check for version mismatches before proceeding
- Refactored mock implementations in tests for better version handling

## [0.3.0] - 2024-06-26

### Added

- New feature: Implemented automatic changelog formatting
- Improved changie to automatically add 0.0.0 tag during initialization

### Changed

- Improved error handling in version bumping process

### Fixed

- Resolved issue with newline handling in changelog entries
- Fixed the bumping

## [0.2.2] - 2024-06-23

### Added

- New feature: Improved changelog management

## [0.2.1] - 2024-06-23

### Fixed

- Fixed issue with comparison URLs in changelog

## [0.2.0] - 2024-06-23

### Added

- CHANGELOG.md
- Updated README.md with improved documentation
- Added CI/CD configuration with ci.yml
- Implemented unit tests and improved test coverage

### Changed

- Modified main.go to improve functionality and error handling

## [0.1.0] - 2024-06-23

### Added

- Initial release with basic functionality

## 0.0.0 - 2016-09-07

### Added

- Initial project setup

[Unreleased]: https://github.com/peiman/changie/compare/v1.1.0...HEAD
[v1.1.0]: https://github.com/peiman/changie/compare/v1.0.2...v1.1.0
[v1.0.2]: https://github.com/peiman/changie/compare/v1.0.1...v1.0.2
[v1.0.1]: https://github.com/peiman/changie/compare/v1.0.0...v1.0.1
[v1.0.0]: https://github.com/peiman/changie/compare/v0.9.1...v1.0.0
[0.9.1]: https://github.com/peiman/changie/compare/v0.9.0...v0.9.1
[0.9.0]: https://github.com/peiman/changie/compare/v0.8.0...v0.9.0
[0.8.0]: https://github.com/peiman/changie/compare/v0.7.0...v0.8.0
[0.7.0]: https://github.com/peiman/changie/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/peiman/changie/compare/v0.5.0...v0.6.0
[0.5.0]: https://github.com/peiman/changie/compare/v0.4.1...v0.5.0
[0.4.1]: https://github.com/peiman/changie/compare/v0.4.0...v0.4.1
[0.4.0]: https://github.com/peiman/changie/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/peiman/changie/compare/v0.2.2...v0.3.0
[0.2.2]: https://github.com/peiman/changie/compare/v0.2.1...v0.2.2
[0.2.1]: https://github.com/peiman/changie/compare/v0.2.0...v0.2.1
[0.2.0]: https://github.com/peiman/changie/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/peiman/changie/releases/tag/v0.1.0