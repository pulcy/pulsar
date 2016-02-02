package golang

import (
	"os"
	"path/filepath"
	"time"

	log "github.com/op/go-logging"

	"git.pulcy.com/pulcy/pulcy/cache"
	"git.pulcy.com/pulcy/pulcy/util"
)

type VendorFlags struct {
	Package   string
	VendorDir string
}

// Get executes a `go get` with a cache support.
func Vendor(log *log.Logger, flags *VendorFlags) error {
	// Get cache dir
	cachedir, _, err := cache.Dir(flags.Package, time.Millisecond)
	if err != nil {
		return maskAny(err)
	}

	// Cache has become invalid
	log.Info(updating("Fetching %s"), flags.Package)
	// Execute `go get` towards the cache directory
	if err := runGoGet(log, flags.Package, cachedir); err != nil {
		return maskAny(err)
	}

	// Sync with vendor dir
	if err := os.MkdirAll(flags.VendorDir, 0777); err != nil {
		return maskAny(err)
	}
	if err := util.ExecPrintError(nil, "rsync", "--exclude", ".git", "-a", filepath.Join(cachedir, srcDir)+"/", flags.VendorDir); err != nil {
		return maskAny(err)
	}

	return nil
}
