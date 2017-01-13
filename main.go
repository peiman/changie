package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Changies prereqs
// Project supports semver v2.0.0
// Project supports "Keep a Changelog" v0.3.0
// Git is installed in host environment
//
// Changies required input
// Version bump: Major | Minor | Patch
// Remote Repository Provider: Github, Bitbucket, default: Github
// Change log file, default: CHANGELOG.md
//
// Changies algorithm
// Check version git tag to determine current version, git describe --tags --abbrev=0 --match "[0-9]*.[0-9]*.[0-9]*"
// Update CHANGELOG.md,
// 	- Check that unreleased is not empty.
// 	- Move everything in unreleased to the new version number.
//	- Add diff link for github or bitbucket
// Commit CHANGELOG.md
// Git tag with the bumped version
// Remind user not to forget to git push

var (
	app                        = kingpin.New("changie", "A version and change log manager for releases. Made for projects using Git, semver v2.0.0 and Keep a Changelog v0.3.0.")
	initCommand                = app.Command("init", "Initiate project directory for semver and Keep a Changelog.")
	majorCommand               = app.Command("major", "Release a major version. Bump the first version number.")
	minorCommand               = app.Command("minor", "Release a minor version. Bump the second version number.")
	patchCommand               = app.Command("patch", "Release a patch version. Bump the third version number.")
	remoteRepositoryProvider   = app.Flag("rrp", "Remote repository provider, github or bitbucket.").Short('r').Default("github").Enum("github", "bitbucket")
	changelogCommand           = app.Command("changelog", "Change log commands.")
	changeLogFile              = changelogCommand.Flag("file", "Change log file name.").Short('f').Default("CHANGELOG.md").File()
	changelogAddCommand        = changelogCommand.Command("added", "Add an added section to changelog.")
	changelogChangedCommand    = changelogCommand.Command("changed", "Add an changed section to changelog.")
	changelogDeprecatedCommand = changelogCommand.Command("deprecated", "Add an deprecated section to changelog.")
	changelogRemovedCommand    = changelogCommand.Command("removed", "Add an removed section to changelog.")
	changelogFixedCommand      = changelogCommand.Command("fixed", "Add an fixed section to changelog.")
	changelogSecurityCommand   = changelogCommand.Command("security", "Add an security section to changelog.")
)

type semverT struct {
	str string
}

func (s *semverT) explode() [3]int {
	var re = regexp.MustCompile(`(\d+).(\d+).(\d+)`)
	allMatches := re.FindAllStringSubmatch(s.str, -1)
	major, err := strconv.Atoi(allMatches[0][1])
	if err != nil {
		log.Fatal(err)
	}
	minor, err := strconv.Atoi(allMatches[0][2])
	if err != nil {
		log.Fatal(err)
	}
	patch, err := strconv.Atoi(allMatches[0][3])
	if err != nil {
		log.Fatal(err)
	}
	return [3]int{major, minor, patch}
}

func (s *semverT) bumpMajor() string {
	exploded := s.explode()
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(exploded[0] + 1))
	buffer.WriteString(".0.0")
	s.str = buffer.String()
	return s.str
}

func (s *semverT) bumpMinor() string {
	exploded := s.explode()
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(exploded[0]))
	buffer.WriteString(".")
	buffer.WriteString(strconv.Itoa(exploded[1] + 1))
	buffer.WriteString(".0")
	s.str = buffer.String()
	return s.str
}

func (s *semverT) bumpPatch() string {
	exploded := s.explode()
	var buffer bytes.Buffer
	buffer.WriteString(strconv.Itoa(exploded[0]))
	buffer.WriteString(".")
	buffer.WriteString(strconv.Itoa(exploded[1]))
	buffer.WriteString(".")
	buffer.WriteString(strconv.Itoa(exploded[2] + 1))
	s.str = buffer.String()
	return s.str
}

func getProjectVersion() string {
	var (
		cmdOut []byte
		err    error
	)
	cmdName := "git"
	cmdArgs := []string{"describe", "--tags", "--abbrev=0", "--match", "[[:digit:]]*.[[:digit:]]*.[[:digit:]]*"}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).CombinedOutput(); err != nil {
		fmt.Print(string(cmdOut[:]))
		fmt.Fprintln(os.Stderr, "There was an error running git describe command ", err)
		os.Exit(1)
	}
	return string(cmdOut)
}

func gitTag(semver string) {
	var (
		cmdOut []byte
		err    error
	)
	cmdName := "git"
	cmdArgs := []string{"tag", semver}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).CombinedOutput(); err != nil {
		fmt.Print(string(cmdOut[:]))
		fmt.Fprintln(os.Stderr, "There was an error running git describe command ", err)
		os.Exit(1)
	}
}

func main() {
	out, err := exec.Command("git", "--version").CombinedOutput()
	gitExists := regexp.MustCompile(`git version`)
	s := string(out[:])
	if err != nil {
		log.Fatal(err)
	} else if !gitExists.MatchString(s) {
		log.Fatal("Git is not installed.")
	}

	app.Version("0.0.0")
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case "init":
		log.Print("Initiate project for semver and Keep a Changelog.")
	case "major":
		semver := semverT{getProjectVersion()}
		semver.bumpMajor()
		//gitTag(semver)
		log.Print("Major release " + semver.str + " done.")
	case "minor":
		semver := semverT{getProjectVersion()}
		semver.bumpMinor()
		//gitTag(semver)
		log.Print("Minor release " + semver.str + " done.")
	case "patch":
		semver := semverT{getProjectVersion()}
		semver.bumpPatch()
		//gitTag(semver)
		log.Print("Patch release " + semver.str + " done.")
	}
	log.Print("Don't forget to git push.")
}
