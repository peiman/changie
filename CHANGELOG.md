# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com),
and this project adheres to [Semantic Versioning (SemVer)](https://semver.org).

## [Unreleased]

### Added

- Ported to the [ckeletin-go](https://github.com/peiman/ckeletin-go) framework for improved structure, testability, and maintainability

### Changed

- Simplified the way Changie retrieves the current version from Git, making it more reliable.
- Improved error messages for better clarity when Git operations fail.
- Enhanced debug messages to help users troubleshoot issues more effectively.
- Enhanced code documentation across all packages for better maintainability
- Added comprehensive package-level documentation to improve developer onboarding
- Improved function documentation with detailed parameter and return value descriptions

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

[Unreleased]: https://github.com/peiman/changie/compare/v0.9.1...HEAD
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
