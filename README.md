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

### Bumping versions

To bump the version, use one of the following commands:

```bash
changie major  # Bump major version (e.g., 1.0.0 -> 2.0.0)
changie minor  # Bump minor version (e.g., 1.0.0 -> 1.1.0)
changie patch  # Bump patch version (e.g., 1.0.0 -> 1.0.1)
```

### Managing the changelog

To add a new section to the changelog, use one of the following commands:

```bash
changie changelog added      # Add an "Added" section
changie changelog changed    # Add a "Changed" section
changie changelog deprecated # Add a "Deprecated" section
changie changelog removed    # Add a "Removed" section
changie changelog fixed      # Add a "Fixed" section
changie changelog security   # Add a "Security" section
```

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

To check test coverage:

```bash
go test ./... -cover
```

### Contributing

Contributions to changie are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.