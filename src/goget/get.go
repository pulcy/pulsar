package get

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/juju/errgo"
	"github.com/mgutz/ansi"
	log "github.com/op/go-logging"

	"git.pulcy.com/pulcy/pulcy/cache"
	"git.pulcy.com/pulcy/pulcy/util"
)

const (
	srcDir          = "src"
	cacheValidHours = 12
)

var (
	maskAny   = errgo.MaskFunc(errgo.Any)
	allGood   = ansi.ColorFunc("")
	updating  = ansi.ColorFunc("cyan")
	attention = ansi.ColorFunc("yellow")
	envMutex  sync.Mutex
	gopath    string
)

type Flags struct {
	Package string
}

func init() {
	gopath = os.Getenv("GOPATH")
}

// Get executes a `go get` with a cache support.
func Get(log *log.Logger, flags *Flags) error {
	// Check GOPATH
	if gopath == "" {
		return maskAny(errors.New("Specify GOPATH"))
	}
	gopathDir := strings.Split(gopath, string(os.PathListSeparator))[0]

	// Get cache dir
	cachedir, cacheIsValid, err := cache.Dir(flags.Package, cacheValidHours)
	if err != nil {
		return maskAny(err)
	}

	if !cacheIsValid {
		// Cache has become invalid
		log.Info(updating("Refreshing cache of %s"), flags.Package)
		// Execute `go get` towards the cache directory
		if err := runGoGet(log, flags.Package, cachedir); err != nil {
			return maskAny(err)
		}
	}

	// Sync with local gopath
	if err := os.MkdirAll(gopathDir, 0777); err != nil {
		return maskAny(err)
	}
	if err := util.ExecPrintError(nil, "rsync", "-a", filepath.Join(cachedir, srcDir), gopathDir); err != nil {
		return maskAny(err)
	}

	return nil
}

func runGoGet(log *log.Logger, pkg, gopath string) error {
	envMutex.Lock()
	defer envMutex.Unlock()

	return func() error {
		// Restore GOPATH on exit
		defer os.Setenv("GOPATH", gopath)
		// Set GOPATH
		if err := os.Setenv("GOPATH", gopath); err != nil {
			return maskAny(err)
		}
		//log.Info("GOPATH=%s", gopath)
		return maskAny(util.ExecPrintError(log, "go", "get", pkg))
	}()

}
