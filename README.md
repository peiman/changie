# changie

Changie is a powerful version and changelog management tool for software projects. It simplifies the process of updating versions and maintaining changelogs, integrating seamlessly with Git, Semantic Versioning (SemVer), and the Keep a Changelog format.

## Features

- **Semantic Versioning Support**: Easily bump major, minor, and patch versions
- **Automated Changelog Management**: Automatically update your CHANGELOG.md file
- **Git Integration**: Commit changes and tag releases with built-in Git commands
- **Multiple Repository Support**: Works with both GitHub and Bitbucket
- **Customizable**: Flexible configuration to fit your project's needs
- **Command-line Interface**: Easy-to-use CLI for all operations

## Installation

To install changie, use the following command:

```bash
go get -u github.com/peiman/changie
```

Ensure you have Go installed on your system (version 1.16 or later).

## Quick Start

1. Initialize your project:

   ```bash
   changie init
   ```

2. Bump the version (choose one):

   ```bash
   changie major  # For major version bump (e.g., 1.0.0 -> 2.0.0)
   changie minor  # For minor version bump (e.g., 1.0.0 -> 1.1.0)
   changie patch  # For patch version bump (e.g., 1.0.0 -> 1.0.1)
   ```

3. Add a new changelog entry:

   ```bash
   changie changelog added "New feature description"
   ```

## Usage

### Initializing a Project

```bash
changie init
```

This command creates a new CHANGELOG.md file in your project directory if it doesn't exist.

### Bumping Versions

```bash
changie major
changie minor
changie patch
```

These commands will:
1. Update the version number
2. Update the CHANGELOG.md file
3. Commit the changes
4. Create a new Git tag

### Managing the Changelog

Add new sections to the changelog:

```bash
changie changelog added "New feature description"
changie changelog changed "Description of changes in existing functionality"
changie changelog deprecated "Description of soon-to-be removed features"
changie changelog removed "Description of removed features"
changie changelog fixed "Description of any bug fixes"
changie changelog security "Description of security vulnerabilities fixed"
```

### Specifying the Remote Repository Provider

By default, changie assumes you're using GitHub. To specify a different provider:

```bash
changie --rrp bitbucket major
```

### Custom Changelog File

If your changelog file is not named CHANGELOG.md or is in a different location:

```bash
changie --file docs/CHANGELOG.md major
```

## Contributing

Contributions to changie are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Support

If you encounter any problems or have any questions, please open an issue on the GitHub repository.

---

Don't forget to star the repository if you find changie useful!