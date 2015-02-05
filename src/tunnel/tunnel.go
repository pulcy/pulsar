package tunnel

import (
	"os"
	"path/filepath"
	"time"

	log "github.com/op/go-logging"

	"arvika.subliminl.com/developers/devtool/util"
)

const (
	registryTunnelSocket = "${HOME}/.ssh-tunnel-subliminl-registry"
	registryHost         = "arvika.subliminl.com"
	mongoTunnelSocket    = "${HOME}/.ssh-tunnel-mongo1-admin"
	mongoHost            = "mongo1.subliminl.com"
)

// Open an SSH tunnel to arvika-ssh
func OpenTunnel(log *log.Logger) error {
	links := []string{
		"-L5000:localhost:5000",
		"-L16012:localhost:16012",
		"-L16073:localhost:16073",
	}
	return openTunnel(log, registryTunnelSocket, registryHost, links...)
}

// Close any existing SSH tunnel to arvika-ssh
func CloseTunnel(log *log.Logger) error {
	return closeTunnel(log, registryTunnelSocket, registryHost)
}

// Open an SSH tunnel to mongo1.subliminl.com
func OpenMongoTunnel(log *log.Logger) error {
	links := []string{
		"-L27018:localhost:27017",
	}
	return openTunnel(log, mongoTunnelSocket, mongoHost, links...)
}

// Close any existing SSH tunnel to mongo1.subliminl.com
func CloseMongoTunnel(log *log.Logger) error {
	return closeTunnel(log, mongoTunnelSocket, mongoHost)
}

// Open an SSH tunnel to given host
func openTunnel(log *log.Logger, tunnelSocket, host string, links ...string) error {
	socket, err := filepath.Abs(os.ExpandEnv(tunnelSocket))
	if err != nil {
		return err
	}
	if _, err := os.Stat(socket); os.IsNotExist(err) {
		args := append([]string{}, links...)
		args = append(args, []string{
			"-f",
			"-N",
			"-M",
			"-S",
			socket,
			"admin@" + host,
		}...)
		if err := util.ExecDetached(log, "ssh", args...); err != nil {
			return err
		}
		// Wait a while for the tunnel to settle
		time.Sleep(time.Millisecond * 750)
	}
	return nil
}

// Close any existing SSH tunnel to given host
func closeTunnel(log *log.Logger, tunnelSocket, host string) error {
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
			"admin@" + host,
		}
		return util.ExecPrintError(log, "ssh", args...)
	}
	return nil
}
