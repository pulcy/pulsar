package main

import (
	"github.com/spf13/cobra"

	"arvika.subliminl.com/developers/devtool/gitlab"
)

const (
	defaultGitlabHost    = "https://arvika.subliminl.com"
	defaultGitlabApiPath = "/api/v3"
)

var (
	gitlabFlags = &gitlab.Config{}
	gitlabCmd   = &cobra.Command{
		Use:   "gitlab",
		Short: "Gitlab utilities",
		Run:   UsageFunc,
	}
)

func init() {
	gitlabCmd.Flags().StringVarP(&gitlabFlags.Host, "host", "", defaultGitlabHost, "Specify gitlab host")
	gitlabCmd.Flags().StringVarP(&gitlabFlags.ApiPath, "api-path", "", defaultGitlabApiPath, "Specify gitlab API path")
	gitlabCmd.Flags().StringVarP(&gitlabFlags.Token, "token", "", "", "Specify gitlab token")
	mainCmd.AddCommand(gitlabCmd)
	gitlabCmd.AddCommand(&cobra.Command{
		Use:   "projects",
		Short: "List all projects",
		Run:   runListGitlabProjects,
	})
}

func mergeDefaultGitlabConfig() {
	if gitlabFlags.Host == "" || gitlabFlags.ApiPath == "" || gitlabFlags.Token == "" {
		defCfg, err := gitlab.GetDefaultConfig()
		if err != nil {
			Quitf("Cannot find gitlab config: %v\n", err)
		}
		if gitlabFlags.Host == "" {
			gitlabFlags.Host = defCfg.Host
		}
		if gitlabFlags.ApiPath == "" {
			gitlabFlags.ApiPath = defCfg.ApiPath
		}
		if gitlabFlags.Token == "" {
			gitlabFlags.Token = defCfg.Token
		}
	}
}

func runListGitlabProjects(cmd *cobra.Command, args []string) {
	mergeDefaultGitlabConfig()
	err := gitlab.ListProjects(gitlabFlags)
	if err != nil {
		Quitf("Cannot list projects: %v\n", err)
	}
}
