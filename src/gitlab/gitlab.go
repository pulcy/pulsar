package gitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bugagazavr/go-gitlab-client"
	"github.com/juju/errgo"

	"arvika.subliminl.com/developers/devtool/git"
)

const (
	configFile = ".subliminl/gitlab"
)

var (
	Mask = errgo.Mask
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
		if _, err := os.Stat(p.Name); err == nil {
			// Folder already exists, don't clone
			continue
		}

		fmt.Printf("Cloning %s\n", p.Name)
		git.Clone(nil, p.SshRepoUrl, p.Name)
	}
	return nil
}
