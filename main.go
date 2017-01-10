package main

import (
	"log"
	"os"
	"os/exec"
	"regexp"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Changies prereqs
// Project supports semver v2.0.0
// Project supports "Keep a Changelog" v0.3.0
// Git is installed in host environment
//
// Changies required input
// Version bump: Major | Minor | Patch
// Remote Repository Provider: Github, Bitbucket
//
// Changies optional input
// Change log file, default CHANGELOG.md
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
	app                      = kingpin.New("changie", "A version and change log manager for releases. Made for projects using Git, semver v2.0.0 and Keep a Changelog v0.3.0.")
	majorCommand             = app.Command("major", "Release a major version. Bump the first version number.")
	minorCommand             = app.Command("minor", "Release a minor version. Bump the second version number.")
	patchCommand             = app.Command("patch", "Release a patch version. Bump the third version number.")
	remoteRepositoryProvider = app.Flag("rrp", "Remote repository provider, github or bitbucket.").Short('r').Default("github").Enum("github", "bitbucket")
	changeLogFile            = app.Flag("changelog", "Change log file name.").Short('c').Default("CHANGELOG.md").ExistingFile()
)

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
	case "major":
		log.Print("major release done")
	case "minor":
		log.Print("major release done")
	case "patch":
		log.Print("major release done")
	}
}
