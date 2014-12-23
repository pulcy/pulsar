package git

import (
	"fmt"
	"strings"

	log "github.com/op/go-logging"

	"arvika.subliminl.com/developers/subliminl/util"
)

const (
	cmdName = "git"
)

// Execute a `git add`
func Add(log *log.Logger, files ...string) error {
	args := []string{"add"}
	args = append(args, files...)
	return util.ExecPrintError(log, cmdName, args...)
	return nil
}

// Execute a `git commit`
func Commit(log *log.Logger, message string) error {
	if msg, err := util.Exec(log, cmdName, "commit", "-m", message); err != nil {
		fmt.Printf("%s\n", msg)
		return err
	}
	return nil
}

// Execute a `git status`
func Status(log *log.Logger, porcelain bool) (string, error) {
	args := []string{"status"}
	if porcelain {
		args = append(args, "--porcelain")
	}
	if msg, err := util.Exec(log, cmdName, args...); err != nil {
		if log != nil {
			log.Error(msg)
		} else {
			fmt.Printf("%s\n", msg)
		}
		return "", err
	} else {
		return strings.TrimSpace(msg), nil
	}
}

// Execute a `git push`
func Push(log *log.Logger, remote string, tags bool) error {
	args := []string{
		"push",
	}
	if tags {
		args = append(args, "--tags")
	}
	if remote != "" {
		args = append(args, remote)
	}
	return util.ExecPrintError(log, cmdName, args...)
}

// Execute a `git tag <tag>`
func Tag(log *log.Logger, tag string) error {
	args := []string{
		"tag",
		tag,
	}
	return util.ExecPrintError(log, cmdName, args...)
}
