package gitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/juju/errgo"
	"github.com/subliminl/go-gitlab-client"

	"arvika.subliminl.com/developers/devtool/git"
)

const (
	configFile = ".subliminl/gitlab"
)

var (
	Mask               = errgo.Mask
	ErrProjectNotFound = errgo.New("Project now found")
)

type Config struct {
	Host    string `json:"host"`
	ApiPath string `json:"api_path"`
	Token   string `json:"token"`
}

func GetDefaultConfig() (*Config, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return nil, Mask(errgo.New("Cannot find HOME"))
	}

	file, err := ioutil.ReadFile(filepath.Join(home, configFile))
	if err != nil {
		return nil, Mask(err)
	}

	config := &Config{}
	if err := json.Unmarshal(file, config); err != nil {
		return nil, Mask(err)
	}

	return config, nil
}

// Show a list of all projects
func ListProjects(config *Config) error {
	gitlab := gogitlab.NewGitlab(config.Host, config.ApiPath, config.Token)
	projects, err := gitlab.AllProjects()
	if err != nil {
		return Mask(err)
	}
	for _, p := range projects {
		if p.Archived {
			continue
		}
		fmt.Printf("%s\n", p.Name)
	}
	return nil
}

// Clone all projects in the current folder
func CloneProjects(config *Config) error {
	gitlab := gogitlab.NewGitlab(config.Host, config.ApiPath, config.Token)
	projects, err := gitlab.AllProjects()
	if err != nil {
		return Mask(err)
	}
	for _, p := range projects {
		if p.Archived {
			continue
		}
		if _, err := os.Stat(p.Name); err == nil {
			// Folder already exists, don't clone
			continue
		}

		fmt.Printf("Cloning %s\n", p.Name)
		git.Clone(nil, p.SshRepoUrl, p.Name)
	}
	return nil
}

// AddPullRequest creates a new pull request for the current branch.
func AddPullRequest(config *Config) error {
	gitlab := gogitlab.NewGitlab(config.Host, config.ApiPath, config.Token)
	id, err := getProjectId(gitlab)
	if err != nil {
		return Mask(err)
	}
	fmt.Println(id)

	targetBranch, err := git.GetLocalBranchName(nil)
	if err != nil {
		return Mask(err)
	}
	fmt.Println(targetBranch)
	//gitlab.AddMergeRequest(id, sourceBranch, targetBranch, title)
	return nil
}

// getProjectId looks up the gitlab project id of the current project
func getProjectId(gitlab *gogitlab.Gitlab) (string, error) {
	url, err := git.GetRemoteOriginUrl(nil)
	if err != nil {
		return "", Mask(err)
	}
	projects, err := gitlab.AllProjects()
	if err != nil {
		return "", Mask(err)
	}
	for _, p := range projects {
		if p.SshRepoUrl == url {
			return strconv.Itoa(p.Id), nil
		}
	}
	return "", ErrProjectNotFound
}
