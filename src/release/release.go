package release

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/blang/semver"
	"github.com/juju/errgo"
	log "github.com/op/go-logging"

	"arvika.subliminl.com/developers/subliminl/docker"
	"arvika.subliminl.com/developers/subliminl/git"
	"arvika.subliminl.com/developers/subliminl/tunnel"
	"arvika.subliminl.com/developers/subliminl/util"
)

const (
	packageJsonFile = "package.json"
	versionFile     = "VERSION"
	nameKey         = "name"
	versionKey      = "version"
	makefileFile    = "Makefile"
	gruntfileFile   = "Gruntfile.js"
	dockerfileFile  = "Dockerfile"
	defaultPerm     = 0664
)

type Flags struct {
	ReleaseType    string
	DockerRegistry string
}

type packageJson map[string]interface{}

func Release(log *log.Logger, flags *Flags) error {
	// Detect environment
	hasMakefile := false
	if _, err := os.Stat(makefileFile); err == nil {
		hasMakefile = true
		log.Info("Found %s", makefileFile)
	}

	hasGruntfile := false
	if _, err := os.Stat(gruntfileFile); err == nil {
		hasGruntfile = true
		log.Info("Found %s", gruntfileFile)
	}

	hasDockerfile := false
	if _, err := os.Stat(dockerfileFile); err == nil {
		hasDockerfile = true
		log.Info("Found %s", dockerfileFile)
	}

	// Read the current version and name
	info, err := GetProjectInfo()
	if err != nil {
		return err
	}

	log.Info("Found old version %s", info.Version)
	version, err := semver.New(info.Version)
	if err != nil {
		return err
	}

	// Check repository state
	if err := checkRepoClean(log); err != nil {
		return err
	}

	// Bump version
	switch flags.ReleaseType {
	case "major":
		version.Major++
		version.Minor = 0
		version.Patch = 0
	case "minor":
		version.Minor++
		version.Patch = 0
	case "patch":
		version.Patch++
	default:
		return errgo.Newf("Unknown release type %s", flags.ReleaseType)
	}
	version.Build = []string{}

	// Write new release version
	if err := writeVersion(log, version.String(), info.pkg, false); err != nil {
		return err
	}

	// Open SSH tunnel
	if err := tunnel.OpenTunnel(log); err != nil {
		return err
	}

	// Build project
	if hasGruntfile {
		if err := util.ExecPrintError(log, "grunt", "build-release"); err != nil {
			return err
		}
	}
	if hasMakefile {
		// Clean first
		if err := util.ExecPrintError(log, "make", "clean"); err != nil {
			return err
		}
		// Now build
		if err := util.ExecPrintError(log, "make"); err != nil {
			return err
		}
	}

	if hasDockerfile {
		// Build docker images
		tag := fmt.Sprintf("%s:%s", info.Name, version.String())
		if err := util.ExecPrintError(log, "docker", "build", "--tag", tag, "."); err != nil {
			return err
		}
		if flags.DockerRegistry != "" {
			// Push image to registry
			if err := docker.Push(log, tag, flags.DockerRegistry); err != nil {
				return err
			}
		}
	}

	// Build succeeded, re-write new release version and commit
	if err := writeVersion(log, version.String(), info.pkg, true); err != nil {
		return err
	}

	// Tag version
	if err := git.Tag(log, version.String()); err != nil {
		return err
	}

	// Update version to "+git" working version
	version.Build = []string{"git"}

	// Write new release version
	if err := writeVersion(log, version.String(), info.pkg, true); err != nil {
		return err
	}

	// Push changes
	if err := git.Push(log, "", false); err != nil {
		return err
	}

	// Push tags
	if err := git.Push(log, "", true); err != nil {
		return err
	}

	return nil
}

// Update the version of the given package (if any) and an existing VERSION file (if any)
// Commit changes afterwards
func writeVersion(log *log.Logger, version string, pkg packageJson, commit bool) error {
	files := []string{}
	if pkg != nil {
		pkg[versionKey] = version
		data, err := json.MarshalIndent(pkg, "", "  ")
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(packageJsonFile, data, defaultPerm); err != nil {
			return err
		}
		files = append(files, packageJsonFile)
	}
	if _, err := os.Stat(versionFile); err == nil {
		if err := ioutil.WriteFile(versionFile, []byte(version), defaultPerm); err != nil {
			return err
		}
		files = append(files, versionFile)
	}

	if commit {
		if err := git.Add(log, files...); err != nil {
			return err
		}
		msg := fmt.Sprintf("Updated version to %s", version)
		if err := git.Commit(log, msg); err != nil {
			return err
		}
	}

	return nil
}

// Are the no uncommited changes in this repo?
func checkRepoClean(log *log.Logger) error {
	if st, err := git.Status(log, true); err != nil {
		return err
	} else if st != "" {
		return errgo.New("There are uncommited changes")
	}
	if err := git.Fetch(log, "origin"); err != nil {
		return err
	}
	if diff, err := git.Diff(log, "master", "origin/master"); err != nil {
		return err
	} else if diff != "" {
		return errgo.New("Master is not in sync with origin")
	}

	return nil
}
