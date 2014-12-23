package util

import (
	"bytes"
	"fmt"
	"os/exec"

	log "github.com/op/go-logging"
)

// Execute a given command.
// Return stderr on error, stdout on no error
func Exec(log *log.Logger, cmdName string, arguments ...string) (string, error) {
	if log != nil {
		log.Debug("Running %s %v", cmdName, arguments)
	}
	cmd := exec.Command(cmdName, arguments...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stderr.String(), err
	} else {
		return stdout.String(), nil
	}
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
	if data, err := Exec(log, cmdName, arguments...); err != nil {
		if log != nil {
			log.Error(data)
		} else {
			fmt.Printf("%s\n", data)
		}
		return err
	}
	return nil
}
