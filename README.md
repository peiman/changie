# changie

Changie is a version and change log manager for releases. It's designed for projects using Git, [Semantic Versioning 2.0.0](https://semver.org), and [Keep a Changelog 1.0.0](https://keepachangelog.com/en/1.0.0/).

## Features

- Semantic versioning support (major, minor, patch)
- Automatic CHANGELOG.md management
- Git integration for version tagging
- Support for different remote repository providers (GitHub, Bitbucket)

## Installation

To install changie, use the following command:

```bash
go get -u github.com/peiman/changie
```

## Usage

### Initializing a project

```bash
changie init
```

This command creates a new CHANGELOG.md file in your project directory.

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

For example:

```bash
changie changelog added "New feature: Improved changelog management"
```

This will add a new entry under the "Added" section in the [Unreleased] part of your CHANGELOG.md file.

### Bumping versions

To bump the version, use one of the following commands:

```bash
changie major  # Bump major version (e.g., 1.0.0 -> 2.0.0)
changie minor  # Bump minor version (e.g., 1.0.0 -> 1.1.0)
changie patch  # Bump patch version (e.g., 1.0.0 -> 1.0.1)
```

When you bump the version, all entries added using the `changie changelog` commands will be moved from the [Unreleased] section to the new version section in the CHANGELOG.md file.
Also a comparison link (the actual version number is the link) will be created in the new version section.

### Specifying the remote repository provider

By default, changie assumes you're using GitHub. To specify a different provider, use the `--rrp` flag:

```bash
changie --rrp bitbucket major
```

## Development

### Running Tests

To run tests:

```bash
go test ./...
```

### Contributing

Contributions to changie are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.