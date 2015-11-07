package util

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"github.com/juju/errgo"
	log "github.com/op/go-logging"
)

var (
	maskAny = errgo.MaskFunc(errgo.Any)
)

type Command struct {
	log       *log.Logger
	stdout    bytes.Buffer
	stderr    bytes.Buffer
	cmd       *exec.Cmd
	cmdName   string
	arguments []string
}

func PrepareCommand(log *log.Logger, cmdName string, arguments ...string) *Command {
	c := &Command{
		log:       log,
		cmdName:   cmdName,
		arguments: arguments,
		cmd:       exec.Command(cmdName, arguments...),
	}
	c.cmd.Stdout = &c.stdout
	c.cmd.Stderr = &c.stderr
	return c
}

func (c *Command) Announce() {
	if c.log != nil {
		c.log.Debug("Running %s %v", c.cmdName, c.arguments)
	}
}

// Connect stdin to the given reader
func (c *Command) SetStdin(r io.Reader) {
	c.cmd.Stdin = r
}

// Connect stdout to the given writer
func (c *Command) SetStdout(w io.Writer) {
	c.cmd.Stdout = w
}

// Connect stderr to the given writer
func (c *Command) SetStderr(w io.Writer) {
	c.cmd.Stderr = w
}

// Sets the directory in which to execute the commands
func (c *Command) SetDir(dir string) {
	c.cmd.Dir = dir
}

// Execute a given command.
// Return stderr on error, stdout on no error
func (c *Command) Run() (string, error) {
	err := c.cmd.Run()
	if err != nil {
		return c.stderr.String(), maskAny(err)
	} else {
		return c.stdout.String(), nil
	}
}

func IsExit(err error) bool {
	_, ok := errgo.Cause(err).(*exec.ExitError)
	return ok
}

// Execute a given command, printing stderr in case of an error
func (c *Command) RunPrintError() error {
	if data, err := c.Run(); err != nil {
		if c.log != nil {
			c.log.Error(data)
		} else {
			fmt.Printf("%s\n", data)
		}
		return err
	}
	return nil
}

// Execute a given command.
// Return stderr on error, stdout on no error
func Exec(log *log.Logger, cmdName string, arguments ...string) (string, error) {
	cmd := PrepareCommand(log, cmdName, arguments...)
	cmd.Announce()
	return cmd.Run()
}

// Execute a given command without waiting for its result.
func ExecDetached(log *log.Logger, cmdName string, arguments ...string) error {
	if log != nil {
		log.Debug("Running %s %v", cmdName, arguments)
	}
	cmd := exec.Command(cmdName, arguments...)
	return cmd.Start()
}

// Execute a given command, printing stderr in case of an error
func ExecPrintError(log *log.Logger, cmdName string, arguments ...string) error {
	cmd := PrepareCommand(log, cmdName, arguments...)
	cmd.Announce()
	return cmd.RunPrintError()
}
