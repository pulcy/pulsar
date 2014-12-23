package tunnel

import (
	"os"
	"path/filepath"

	log "github.com/op/go-logging"

	"arvika.subliminl.com/developers/subliminl/util"
)

const (
	tunnelSocket = "${HOME}/.ssh-tunnel-subliminl-registry"
	registryhost = "arvika.subliminl.com"
)

// Open an SSH tunnel to arvika-ssh
func OpenTunnel(log *log.Logger) error {
	socket, err := filepath.Abs(os.ExpandEnv(tunnelSocket))
	if err != nil {
		return err
	}
	if _, err := os.Stat(socket); os.IsNotExist(err) {
		args := []string{
			"-f",
			"-L5000:localhost:5000",
			"-L16012:localhost:16012",
			"-L16073:localhost:16073",
			"-N",
			"-M",
			"-S",
			socket,
			"admin@" + registryhost,
		}
		return util.ExecDetached(log, "ssh", args...)
	}
	return nil
}

// Close any existing SSH tunnel to arvika-ssh
func CloseTunnel(log *log.Logger) error {
	socket, err := filepath.Abs(os.ExpandEnv(tunnelSocket))
	if err != nil {
		return err
	}
	if _, err := os.Stat(socket); err == nil {
		args := []string{
			"-S",
			socket,
			"-O",
			"exit",
			"admin@" + registryhost,
		}
		return util.ExecPrintError(log, "ssh", args...)
	}
	return nil
}
