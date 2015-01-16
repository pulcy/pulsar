package get

import (
	"os"
	"path/filepath"

	log "github.com/op/go-logging"

	"arvika.subliminl.com/developers/subliminl/git"
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
			if err := git.Pull(log, "origin"); err != nil {
				return err
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
			log.Info("Found version %s, wanted %s", localVersion, flags.Version)
			// Fetch latest changes
			if err := git.Fetch(log, "origin"); err != nil {
				return err
			}
			// Checkout intended version
			if err := git.Checkout(log, flags.Version); err != nil {
				return err
			}
		} else {
			//log.Info("Found correct version of %s", flags.Folder)
		}
		// Get latest remote version
		remoteVersion, err := git.GetLatestRemoteTag(nil, flags.RepoUrl)
		if err != nil {
			return err
		}
		if remoteVersion != flags.Version {
			log.Warning("Latest remote version '%s' is different from requested version '%s'", remoteVersion, flags.Version)
		}
	}
	return nil
}
