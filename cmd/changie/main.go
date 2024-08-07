package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	HasUncommittedChanges() (bool, error)
	PushChanges() error
	GetVersion() (string, error)
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
func (m DefaultGitManager) TagVersion(version string) error { return git.TagVersion(version) }
func (m DefaultGitManager) GetVersion() (string, error)     { return git.GetVersion() }
func (m DefaultGitManager) HasUncommittedChanges() (bool, error) {
	return git.HasUncommittedChanges()
}
func (m DefaultGitManager) PushChanges() error {
	return git.PushChanges()
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
	app                        = kingpin.New("changie", "A version and change log manager for releases. Made for projects using Git, SemVer and Keep a Changelog.")
	initCommand                = app.Command("init", "Initiate project directory for SemVer and Keep a Changelog.")
	majorCommand               = app.Command("major", "Release a major version. Bump the first version number.")
	minorCommand               = app.Command("minor", "Release a minor version. Bump the second version number.")
	patchCommand               = app.Command("patch", "Release a patch version. Bump the third version number.")
	remoteRepositoryProvider   = app.Flag("rrp", "Remote repository provider, github or bitbucket.").Short('r').Default("github").Enum("github", "bitbucket")
	autoPush                   = app.Flag("auto-push", "Automatically push changes and tags after version bump").Bool()
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
var exitFunction = os.Exit

func handleError(err error) {
	if err != nil {
		fmt.Printf("Debug: handleError called with error: %v\n", err)
		fmt.Fprintln(os.Stderr, err)
		exitFunction(1)
	}
}

func checkVersionMismatch(gitManager GitManager, changelogManager ChangelogManager, printWarning bool) error {
	gitVersion, err := gitManager.GetVersion()
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
		err := fmt.Errorf("Version mismatch: Git tag version %s does not match changelog version %s", gitVersion, changelogVersion)
		if printWarning {
			fmt.Println("Warning:", err)
		}
		return err
	}

	return nil
}

func handleVersionBump(bumpType string, changelogManager ChangelogManager, gitManager GitManager, semverManager SemverManager) error {
	hasUncommittedChanges, err := gitManager.HasUncommittedChanges()
	if err != nil {
		return fmt.Errorf("Error checking for uncommitted changes: %v", err)
	}
	if hasUncommittedChanges {
		return fmt.Errorf("Error: Uncommitted changes found. Please commit or stash your changes before bumping the version.")
	}

	if err := checkVersionMismatch(gitManager, changelogManager, !isTestMode); err != nil {
		return err
	}

	gitVersion, err := gitManager.GetVersion()
	if err != nil {
		return fmt.Errorf("Error getting project version: %v", err)
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

	fmt.Printf("Tagging version: %s\n", newVersion)
	if err := gitManager.TagVersion(newVersion); err != nil {
		return fmt.Errorf("Error tagging version: %v", err)
	}

	fmt.Printf("%s release %s done.\n", bumpType, newVersion)

	if *autoPush {
		fmt.Println("Pushing changes and tags...")
		if err := gitManager.PushChanges(); err != nil {
			return fmt.Errorf("Error pushing changes: %v", err)
		}
		fmt.Println("Automatically pushed changes and tags to remote repository.")
	} else {
		fmt.Println("Don't forget to git push and git push --tags.")
	}

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

func run(changelogManager ChangelogManager, gitManager GitManager, semverManager SemverManager) error {
	fmt.Println("Debug: Entering run function")

	if !isGitInstalled() {
		return fmt.Errorf("Error: Git is not installed.")
	}

	// Get the git tag
	version, err := gitManager.GetVersion()
	fmt.Printf("Debug: GetVersion result: version=%s, err=%v\n", version, err)
	if err != nil {
		return fmt.Errorf("Error getting project version: %w", err)
	}
	app.Version(version)

	command, err := app.Parse(os.Args[1:])
	if err != nil {
		return fmt.Errorf("Error parsing command: %w", err)
	}

	switch command {
	case initCommand.FullCommand():
		log.Printf("Initializing project with changelog file: %s", *changeLogFile)
		handleError(changelogManager.InitProject(*changeLogFile))
		fmt.Println("Project initialized for SemVer and Keep a Changelog.")

	case majorCommand.FullCommand():
		return handleVersionBump("major", changelogManager, gitManager, semverManager)
	case minorCommand.FullCommand():
		return handleVersionBump("minor", changelogManager, gitManager, semverManager)
	case patchCommand.FullCommand():
		return handleVersionBump("patch", changelogManager, gitManager, semverManager)

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
	handleError(run(changelogManager, gitManager, semverManager))
}
