package git

import (
	"bufio"
	"fmt"
	"strings"

	log "github.com/op/go-logging"

	"arvika.subliminl.com/developers/subliminl/util"
)

const (
	cmdName   = "git"
	tagMarker = "refs/tags/"
)

type TagList []string

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

// Execute a `git status a b`
func Diff(log *log.Logger, a, b string) (string, error) {
	args := []string{"diff",
		a,
		b,
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

// Execute a `git pull`
func Pull(log *log.Logger, remote string) error {
	args := []string{
		"pull",
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

// Execute a `git fetch <remote>`
func Fetch(log *log.Logger, remote string) error {
	args := []string{
		"fetch",
		remote,
	}
	return util.ExecPrintError(log, cmdName, args...)
}

// Execute a `git clone <repo-url> <folder>`
func Clone(log *log.Logger, repoUrl, folder string) error {
	args := []string{
		"clone",
		repoUrl,
		folder,
	}
	return util.ExecPrintError(log, cmdName, args...)
}

// Gets the latest tag from the repo in given folder.
func GetLatestTag(log *log.Logger, folder string) (string, error) {
	args := []string{
		"describe",
		"--abbrev=0",
		"--tags",
	}
	cmd := util.PrepareCommand(log, cmdName, args...)
	cmd.SetDir(folder)
	output, err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// Execute a `git checkout <branch>`
func Checkout(log *log.Logger, branch string) error {
	args := []string{
		"checkout",
		branch,
	}
	return util.ExecPrintError(log, cmdName, args...)
}

// Gets the tags from the given remote git repo.
func GetRemoteTags(log *log.Logger, repoUrl string) (TagList, error) {
	args := []string{
		"ls-remote",
		"--tags",
		repoUrl,
	}
	output, err := util.Exec(log, cmdName, args...)
	if err != nil {
		return []string{}, err
	}
	tags := TagList{}
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		index := strings.Index(line, tagMarker)
		if index < 0 {
			continue
		}
		tag := line[index+len(tagMarker):]
		tags = append(tags, tag)
	}
	if err := scanner.Err(); err != nil {
		return tags, err
	}
	return tags, nil
}

// Gets the latest tags from the given remote git repo.
func GetLatestRemoteTag(log *log.Logger, repoUrl string) (string, error) {
	tags, err := GetRemoteTags(log, repoUrl)
	if err != nil {
		return "", err
	}
	if len(tags) > 0 {
		return tags[len(tags)-1], nil
	}
	return "", nil
}

// Gets the latest commit hash from the given local git folder.
func GetLatestLocalCommit(log *log.Logger, folder, branch string) (string, error) {
	if branch == "" {
		branch = "HEAD"
	}
	args := []string{
		"rev-parse",
		branch,
	}
	output, err := util.Exec(log, cmdName, args...)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// Gets the latest commit hash from the given remote git repo + optional branch.
func GetLatestRemoteCommit(log *log.Logger, repoUrl, branch string) (string, error) {
	args := []string{
		"ls-remote",
		repoUrl,
	}
	if branch != "" {
		args = append(args, branch)
	}
	output, err := util.Exec(log, cmdName, args...)
	if err != nil {
		return "", err
	}
	parts := strings.Split(output, "\t")
	return parts[0], nil
}
