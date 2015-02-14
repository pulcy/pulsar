package main

import (
	"github.com/coreos/go-semver/semver"
	"github.com/spf13/cobra"

	"arvika.subliminl.com/developers/devtool/release"
)

var (
	dockerTagImage string
	dockerTagCmd   = &cobra.Command{
		Use:   "docker-tag",
		Short: "Get the docker tag for the current project",
		Long:  "Returns the image:tag for the current project",
		Run:   runDockerTag,
	}
)

func init() {
	dockerTagCmd.Flags().StringVarP(&dockerTagImage, "image", "i", "", "Docker image name")
	mainCmd.AddCommand(dockerTagCmd)
}

func runDockerTag(cmd *cobra.Command, args []string) {
	info, err := release.GetProjectInfo()
	if err != nil {
		Quitf("%s\n", err)
	}
	version, err := semver.NewVersion(info.Version)
	if err != nil {
		Quitf("%s\n", err)
	}
	tag := version.String()
	if version.Metadata != "" {
		tag = "latest"
	}
	if dockerTagImage == "" {
		dockerTagImage = info.Name
	}
	Printf("%s:%s", dockerTagImage, tag)
}
