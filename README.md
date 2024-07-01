# changie

Changie is a version and change log manager for releases. It's designed for projects using Git, [Semantic Versioning 2.0.0](https://semver.org), and [Keep a Changelog 1.0.0](https://keepachangelog.com/en/1.0.0/).

## Features

- Semantic versioning support (major, minor, patch)
- Automatic CHANGELOG.md management
- Git integration for version tagging
- Support for different remote repository providers (GitHub, Bitbucket)

## Quick Start

### Installation

To install changie, use the following command:

```bash
go get -u github.com/peiman/changie
```

### Basic Usage

1. Initialize your project:

```bash
changie init
```

2. Add a changelog entry:

```bash
changie changelog added "New feature: Improved error handling"
```

3. Bump the version:

```bash
changie minor
```

## Detailed Usage

### Managing the changelog

To add a new entry to the changelog, use one of the following commands:

```bash
changie changelog added "Description of new feature"
changie changelog changed "Description of changes in existing functionality"
changie changelog deprecated "Description of soon-to-be removed features"
changie changelog removed "Description of removed features"
changie changelog fixed "Description of any bug fixes"
changie changelog security "Description of security vulnerabilities fixed"
```

### Bumping versions

To bump the version, use one of the following commands:

```bash
changie major  # Bump major version (e.g., 1.3.2 -> 2.0.0)
changie minor  # Bump minor version (e.g., 1.3.2 -> 1.4.0)
changie patch  # Bump patch version (e.g., 1.3.2 -> 1.3.3)
```

### Automatic pushing

To bump the version and automatically push changes and tags, use the `--auto-push` flag:

```bash
changie minor --auto-push
```

### Specifying the remote repository provider

By default, changie assumes you're using GitHub. To specify a different provider, use the `--rrp` flag:

```bash
changie --rrp bitbucket major
```

## Configuration

Changie doesn't require any configuration files. It uses command-line flags for customization.

## Troubleshooting

### Version mismatch between Git tag and Changelog

If you encounter a warning about version mismatch, ensure that your Git tags and CHANGELOG.md are in sync. You may need to manually edit the changelog or create a new Git tag.

### Git is not installed

Changie requires Git to be installed and available in your system's PATH. Ensure Git is properly installed and accessible from the command line.

## Contributing

Contributions to changie are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
