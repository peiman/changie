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
	AddChangelogSection(string, string, string) error
}

type GitManager interface {
	CommitChangelog(string, string) error
	TagVersion(string) error
	GetProjectVersion() (string, error)
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
func (m DefaultChangelogManager) AddChangelogSection(file, section, content string) error {
	return changelog.AddChangelogSection(file, section, content)
}

type DefaultGitManager struct{}

func (m DefaultGitManager) CommitChangelog(file, version string) error {
	return git.CommitChangelog(file, version)
}
func (m DefaultGitManager) TagVersion(version string) error    { return git.TagVersion(version) }
func (m DefaultGitManager) GetProjectVersion() (string, error) { return git.GetProjectVersion() }

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

func handleVersionBump(bumpType string, changelogManager ChangelogManager, gitManager GitManager, semverManager SemverManager) error {
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

	currentVersion, err := gitManager.GetProjectVersion()
	if err != nil {
		return fmt.Errorf("Error getting project version: %v", err)
	}

	newVersion, err := bumpFunc(currentVersion)
	if err != nil {
		return fmt.Errorf("Error bumping version: %v", err)
	}

	changelogFilePath := filepath.Join(".", *changeLogFile)
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
	if err := changelogManager.AddChangelogSection(*changeLogFile, section, content); err != nil {
		return fmt.Errorf("Error adding changelog section: %v", err)
	}
	fmt.Printf("Added %s section to changelog: %s\n", section, content)
	return nil
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
	if err := run(changelogManager, gitManager, semverManager); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
