package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/alecthomas/kingpin/v2"
	"github.com/peiman/changie/internal/changelog"
	"github.com/peiman/changie/internal/git"
	"github.com/peiman/changie/internal/semver"
)

// Interfaces for dependency injection
type ChangelogManager interface {
	InitProject(string) error
	UpdateChangelog(string, string, string) error
	AddChangelogSection(string, string, string) (bool, error)
	GetChangelogContent() (string, error)
}

type GitManager interface {
	CommitChangelog(string, string) error
	TagVersion(string) error
	GetProjectVersion() (string, error)
	HasUncommittedChanges() (bool, error)
}

type SemverManager interface {
	BumpMajor(string) (string, error)
	BumpMinor(string) (string, error)
	BumpPatch(string) (string, error)
}

// Default implementations
type DefaultChangelogManager struct{}

func (m DefaultChangelogManager) InitProject(file string) error { return changelog.InitProject(file) }
func (m DefaultChangelogManager) UpdateChangelog(file, version, provider string) error {
	return changelog.UpdateChangelog(file, version, provider)
}
func (m DefaultChangelogManager) AddChangelogSection(file, section, content string) (bool, error) {
	return changelog.AddChangelogSection(file, section, content)
}

func (m DefaultChangelogManager) GetChangelogContent() (string, error) {
	content, err := os.ReadFile(*changeLogFile)
	if err != nil {
		return "", fmt.Errorf("failed to read changelog: %v", err)
	}
	return string(content), nil
}

type DefaultGitManager struct{}

func (m DefaultGitManager) CommitChangelog(file, version string) error {
	return git.CommitChangelog(file, version)
}
func (m DefaultGitManager) TagVersion(version string) error    { return git.TagVersion(version) }
func (m DefaultGitManager) GetProjectVersion() (string, error) { return git.GetProjectVersion() }
func (m DefaultGitManager) HasUncommittedChanges() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %v", err)
	}
	return len(output) > 0, nil
}

type DefaultSemverManager struct{}

func (m DefaultSemverManager) BumpMajor(version string) (string, error) {
	return semver.BumpMajor(version)
}
func (m DefaultSemverManager) BumpMinor(version string) (string, error) {
	return semver.BumpMinor(version)
}
func (m DefaultSemverManager) BumpPatch(version string) (string, error) {
	return semver.BumpPatch(version)
}

var (
	app                        = kingpin.New("changie", "A version and change log manager for releases. Made for projects using Git, SemVer v2.0.0 and Keep a Changelog v1.0.0.")
	initCommand                = app.Command("init", "Initiate project directory for semver and Keep a Changelog.")
	majorCommand               = app.Command("major", "Release a major version. Bump the first version number.")
	minorCommand               = app.Command("minor", "Release a minor version. Bump the second version number.")
	patchCommand               = app.Command("patch", "Release a patch version. Bump the third version number.")
	remoteRepositoryProvider   = app.Flag("rrp", "Remote repository provider, github or bitbucket.").Short('r').Default("github").Enum("github", "bitbucket")
	changelogCommand           = app.Command("changelog", "Change log commands.")
	changeLogFile              = app.Flag("file", "Change log file name.").Short('f').Default("CHANGELOG.md").String()
	changelogAddCommand        = changelogCommand.Command("added", "Add an added section to changelog.")
	changelogAddContent        = changelogAddCommand.Arg("content", "Content to add to the changelog").Required().String()
	changelogChangedCommand    = changelogCommand.Command("changed", "Add a changed section to changelog.")
	changelogChangedContent    = changelogChangedCommand.Arg("content", "Content to add to the changelog").Required().String()
	changelogDeprecatedCommand = changelogCommand.Command("deprecated", "Add a deprecated section to changelog.")
	changelogDeprecatedContent = changelogDeprecatedCommand.Arg("content", "Content to add to the changelog").Required().String()
	changelogRemovedCommand    = changelogCommand.Command("removed", "Add a removed section to changelog.")
	changelogRemovedContent    = changelogRemovedCommand.Arg("content", "Content to add to the changelog").Required().String()
	changelogFixedCommand      = changelogCommand.Command("fixed", "Add a fixed section to changelog.")
	changelogFixedContent      = changelogFixedCommand.Arg("content", "Content to add to the changelog").Required().String()
	changelogSecurityCommand   = changelogCommand.Command("security", "Add a security section to changelog.")
	changelogSecurityContent   = changelogSecurityCommand.Arg("content", "Content to add to the changelog").Required().String()
)

var isGitInstalled = git.IsInstalled
var isTestMode bool

func handleVersionBump(bumpType string, changelogManager ChangelogManager, gitManager GitManager, semverManager SemverManager) error {
	hasUncommittedChanges, err := gitManager.HasUncommittedChanges()
	if err != nil {
		return fmt.Errorf("Error checking for uncommitted changes: %v", err)
	}
	if hasUncommittedChanges {
		return fmt.Errorf("Error: Uncommitted changes found. Please commit or stash your changes before bumping the version.")
	}
	gitVersion, err := gitManager.GetProjectVersion()
	if err != nil {
		return fmt.Errorf("Error getting project version: %v", err)
	}

	changelogContent, err := changelogManager.GetChangelogContent()
	if err != nil {
		return fmt.Errorf("Error reading changelog: %v", err)
	}

	changelogVersion, err := changelog.GetLatestChangelogVersion(changelogContent)
	if err != nil {
		return fmt.Errorf("Error getting changelog version: %v", err)
	}

	if gitVersion != changelogVersion {
		if !isTestMode {
			fmt.Printf("Warning: Git tag version %s does not match changelog version %s.\n", gitVersion, changelogVersion)
		}
		return fmt.Errorf("Version mismatch: Git tag version %s does not match changelog version %s", gitVersion, changelogVersion)
	}

	fmt.Printf("Current version from git tags: %s\n", gitVersion)

	var bumpFunc func(string) (string, error)
	switch bumpType {
	case "major":
		bumpFunc = semverManager.BumpMajor
	case "minor":
		bumpFunc = semverManager.BumpMinor
	case "patch":
		bumpFunc = semverManager.BumpPatch
	default:
		return fmt.Errorf("Invalid bump type: %s", bumpType)
	}

	newVersion, err := bumpFunc(gitVersion)
	if err != nil {
		return fmt.Errorf("Error bumping version: %v", err)
	}

	fmt.Printf("New version: %s\n", newVersion)

	changelogFilePath := filepath.Join(".", *changeLogFile)
	fmt.Printf("Updating changelog file: %s\n", changelogFilePath)

	if err := changelogManager.UpdateChangelog(changelogFilePath, newVersion, *remoteRepositoryProvider); err != nil {
		return fmt.Errorf("Error updating changelog: %v", err)
	}

	if err := gitManager.CommitChangelog(changelogFilePath, newVersion); err != nil {
		return fmt.Errorf("Error committing changelog: %v", err)
	}

	if err := gitManager.TagVersion(newVersion); err != nil {
		return fmt.Errorf("Error tagging version: %v", err)
	}

	fmt.Printf("%s release %s done.\n", bumpType, newVersion)
	fmt.Println("Don't forget to git push and git push --tags.")
	return nil
}

func handleChangelogUpdate(section, content string, changelogManager ChangelogManager) error {
	isDuplicate, err := changelogManager.AddChangelogSection(*changeLogFile, section, content)
	if err != nil {
		return fmt.Errorf("Error adding changelog section: %v", err)
	}

	if isDuplicate {
		fmt.Printf("%s section: %s (duplicate entry, not added)\n", section, content)
	} else {
		fmt.Printf("%s section: %s\n", section, content)
	}

	return nil
}

func checkVersionMismatch(gitManager GitManager, changelogManager ChangelogManager) error {
	gitVersion, err := gitManager.GetProjectVersion()
	if err != nil {
		return fmt.Errorf("Error getting project version: %v", err) // Updated this line
	}

	changelogContent, err := changelogManager.GetChangelogContent()
	if err != nil {
		return fmt.Errorf("Error reading changelog: %v", err)
	}

	changelogVersion, err := changelog.GetLatestChangelogVersion(changelogContent)
	if err != nil {
		return fmt.Errorf("Error getting changelog version: %v", err)
	}

	if gitVersion != changelogVersion {
		return fmt.Errorf("Version mismatch: Git tag version %s does not match changelog version %s", gitVersion, changelogVersion)
	}

	return nil
}
func getLatestChangelogVersion(content string) (string, error) {
	re := regexp.MustCompile(`## \[(\d+\.\d+\.\d+)\]`)
	matches := re.FindStringSubmatch(content)
	if len(matches) < 2 {
		return "", fmt.Errorf("no version found in changelog")
	}
	return matches[1], nil
}

func run(changelogManager ChangelogManager, gitManager GitManager, semverManager SemverManager) error {
	app.Version("0.1.0")

	if !isGitInstalled() {
		return fmt.Errorf("Error: Git is not installed.")
	}

	command, err := app.Parse(os.Args[1:])
	if err != nil {
		return err // Return the error from kingpin
	}

	switch command {
	case initCommand.FullCommand():
		log.Printf("Initializing project with changelog file: %s", *changeLogFile)
		if err := changelogManager.InitProject(*changeLogFile); err != nil {
			log.Printf("Error initializing project: %v", err)
			fmt.Fprintf(os.Stderr, "Error initializing project: %v\n", err)
			return err
		}
		fmt.Println("Project initialized for semver and Keep a Changelog.")

	case majorCommand.FullCommand():
		err := handleVersionBump("major", changelogManager, gitManager, semverManager)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
	case minorCommand.FullCommand():
		err := handleVersionBump("minor", changelogManager, gitManager, semverManager)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}
	case patchCommand.FullCommand():
		err := handleVersionBump("patch", changelogManager, gitManager, semverManager)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}

	case changelogAddCommand.FullCommand():
		return handleChangelogUpdate("Added", *changelogAddContent, changelogManager)
	case changelogChangedCommand.FullCommand():
		return handleChangelogUpdate("Changed", *changelogChangedContent, changelogManager)
	case changelogDeprecatedCommand.FullCommand():
		return handleChangelogUpdate("Deprecated", *changelogDeprecatedContent, changelogManager)
	case changelogRemovedCommand.FullCommand():
		return handleChangelogUpdate("Removed", *changelogRemovedContent, changelogManager)
	case changelogFixedCommand.FullCommand():
		return handleChangelogUpdate("Fixed", *changelogFixedContent, changelogManager)
	case changelogSecurityCommand.FullCommand():
		return handleChangelogUpdate("Security", *changelogSecurityContent, changelogManager)

	default:
		return fmt.Errorf("Unknown command: %s", command)
	}

	return nil
}

func main() {
	// Enable verbose logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	changelogManager := DefaultChangelogManager{}
	gitManager := DefaultGitManager{}
	semverManager := DefaultSemverManager{}
	if err := run(changelogManager, gitManager, semverManager); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
