package main

import (
	"github.com/blang/semver"
	"github.com/spf13/cobra"

	"arvika.subliminl.com/developers/devtool/release"
)

var (
	dockerTagCmd = &cobra.Command{
		Use:   "docker-tag",
		Short: "Get the docker tag for the current project",
		Long:  "Returns the image:tag for the current project",
		Run:   runDockerTag,
	}
)

func init() {
	mainCmd.AddCommand(dockerTagCmd)
}

func runDockerTag(cmd *cobra.Command, args []string) {
	info, err := release.GetProjectInfo()
	if err != nil {
		Quitf("%s\n", err)
	}
	version, err := semver.New(info.Version)
	if err != nil {
		Quitf("%s\n", err)
	}
	tag := version.String()
	if len(version.Build) > 0 {
		tag = "latest"
	}
	Printf("%s:%s", info.Name, tag)
}
