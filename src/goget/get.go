package get

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mgutz/ansi"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/op/go-logging"

	"arvika.subliminl.com/developers/devtool/util"
)

const (
	cacheDir        = "~/devtool-cache"
	srcDir          = "src"
	cacheValidHours = 12
)

var (
	allGood   = ansi.ColorFunc("")
	updating  = ansi.ColorFunc("cyan")
	attention = ansi.ColorFunc("yellow")
)

type Flags struct {
	Package string
}

// Get executes a `go get` with a cache support.
func Get(log *log.Logger, flags *Flags) error {
	// Get cache folder
	cachedirRoot, err := homedir.Expand(cacheDir)
	if err != nil {
		return err
	}

	// Get GOPATH
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return errors.New("Specify GOPATH")
	}
	gopathDir := strings.Split(gopath, string(os.PathListSeparator))[0]

	// Create hash of package
	hashBytes := sha1.Sum([]byte(flags.Package))
	hash := fmt.Sprintf("%x", hashBytes)
	cachedir := filepath.Join(cachedirRoot, hash)

	// Check if cache directory exists
	s, err := os.Stat(cachedir)
	goGetNeeded := false
	if err == nil {
		// Package cache directory exists, check age.
		if s.ModTime().Add(time.Hour * cacheValidHours).Before(time.Now()) {
			// Cache has become invalid
			log.Info(updating("Refreshing cache of %s"), flags.Package)
			goGetNeeded = true
			if err := os.RemoveAll(cachedir); err != nil {
				return err
			}
		}
	} else {
		// Package cache directory not found, create needed
		goGetNeeded = true
	}

	if goGetNeeded {
		// Execute `go get` towards the cache directory
		// Create cachedir
		if err := os.MkdirAll(cachedir, 0777); err != nil {
			return err
		}
		if err := runGoGet(log, flags.Package, cachedir); err != nil {
			return err
		}
	}

	// Sync with local gopath
	if err := os.MkdirAll(gopathDir, 0777); err != nil {
		return err
	}
	if err := util.ExecPrintError(nil, "rsync", "-a", filepath.Join(cachedir, srcDir), gopathDir); err != nil {
		return err
	}

	return nil
}

func runGoGet(log *log.Logger, pkg, gopath string) error {
	// Save existing GOPATH
	oldGopath := os.Getenv("GOPATH")
	defer os.Setenv("GOPATH", oldGopath)
	// Set GOPATH
	if err := os.Setenv("GOPATH", gopath); err != nil {
		return err
	}
	//log.Info("GOPATH=%s", gopath)
	return util.ExecPrintError(log, "go", "get", pkg)
}
