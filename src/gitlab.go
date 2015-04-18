package main

import (
	"github.com/spf13/cobra"

	"arvika.pulcy.com/developers/devtool/gitlab"
)

const (
	defaultGitlabHost           = "https://arvika.pulcy.com"
	defaultGitlabApiPath        = "/api/v3"
	defaultGitlabPrTargetBranch = "master"
)

var (
	gitlabFlags = &gitlab.Config{}
	gitlabCmd   = &cobra.Command{
		Use:   "gitlab",
		Short: "Gitlab utilities",
		Run:   UsageFunc,
	}

	gitlabPrTargetBranch string
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
	gitlabCmd.AddCommand(&cobra.Command{
		Use:   "clone-all",
		Short: "Clone all projects",
		Run:   runCloneGitlabProjects,
	})
	prCmd := &cobra.Command{
		Use:   "pr",
		Short: "Create pull request",
		Run:   runGitlabCreatePullRequest,
	}
	prCmd.Flags().StringVarP(&gitlabPrTargetBranch, "target", "", defaultGitlabPrTargetBranch, "Specify target branch")
	gitlabCmd.AddCommand(prCmd)
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
		Quitf("Failed to list projects: %v\n", err)
	}
}

func runCloneGitlabProjects(cmd *cobra.Command, args []string) {
	mergeDefaultGitlabConfig()
	err := gitlab.CloneProjects(gitlabFlags)
	if err != nil {
		Quitf("Failed to clone projects: %v\n", err)
	}
}

func runGitlabCreatePullRequest(cmd *cobra.Command, args []string) {
	mergeDefaultGitlabConfig()

	if len(args) == 0 {
		Quitf("Please provide a title\n")
	}
	title := args[0]
	err := gitlab.AddPullRequest(gitlabFlags, gitlabPrTargetBranch, title)
	if err != nil {
		Quitf("Failed to add PR: %v\n", err)
	}
}
