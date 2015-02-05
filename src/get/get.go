package get

import (
	"os"
	"path/filepath"

	log "github.com/op/go-logging"

	"arvika.subliminl.com/developers/subliminl/git"
)

const (
	defaultGetBranch = "master"
)

type Flags struct {
	Folder  string
	RepoUrl string
	Version string
}

// Get ensures that flags.Folder contains an up to date copy of flags.RepoUrl checked out to flags.Version.
func Get(log *log.Logger, flags *Flags) error {
	// Expand folder
	var err error
	flags.Folder, err = filepath.Abs(flags.Folder)
	if err != nil {
		return err
	}

	// Get current folder
	wd, _ := os.Getwd()

	// Make sure a clone exists
	_, err = os.Stat(flags.Folder)
	cloned := false
	if os.IsNotExist(err) {
		// Clone repo into folder
		if err := git.Clone(log, flags.RepoUrl, flags.Folder); err != nil {
			return err
		}
		cloned = true
	}
	// Change dir to folder
	if err := os.Chdir(flags.Folder); err != nil {
		return err
	}
	// Specific version needed?
	if flags.Version == "" {
		// Get latest version
		if !cloned {
			localCommit, err := git.GetLatestLocalCommit(nil, flags.Folder, defaultGetBranch)
			if err != nil {
				return err
			}
			remoteCommit, err := git.GetLatestRemoteCommit(nil, flags.RepoUrl, defaultGetBranch)
			if err != nil {
				return err
			}
			if localCommit != remoteCommit {
				if err := git.Pull(log, "origin"); err != nil {
					return err
				}
			} else {
				log.Info("%s is up to date\n", makeRel(wd, flags.Folder))
			}
		}
	} else {
		// Get latest (local) version
		localVersion, err := git.GetLatestTag(nil, flags.Folder)
		if err != nil {
			return err
		}
		if localVersion != flags.Version {
			// Checkout requested version
			if cloned {
				log.Info("Checking out version %s in %s.\n", flags.Version, makeRel(wd, flags.Folder))
			} else {
				log.Info("Found version %s, wanted %s. Updating %s now.\n", localVersion, flags.Version, makeRel(wd, flags.Folder))
			}
			// Fetch latest changes
			if err := git.Fetch(log, "origin"); err != nil {
				return err
			}
			if err := git.FetchTags(log, "origin"); err != nil {
				return err
			}
			// Checkout intended version
			if err := git.Checkout(log, flags.Version); err != nil {
				return err
			}
		} else {
			log.Info("Found correct version. No changes needed in %s\n", makeRel(wd, flags.Folder))
		}
		// Get latest remote version
		remoteVersion, err := git.GetLatestRemoteTag(nil, flags.RepoUrl)
		if err != nil {
			return err
		}
		if remoteVersion != flags.Version {
			log.Warning("Update available for %s: '%s' => '%s'\n", makeRel(wd, flags.Folder), flags.Version, remoteVersion)
		}
	}
	return nil
}

// makeRel tries to make the given path relative to the current directory.
// Returns a full path in case of errors.
func makeRel(wd, path string) string {
	rel, err := filepath.Rel(wd, path)
	if err != nil {
		return path
	}
	return rel
}
