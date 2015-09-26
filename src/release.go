package main

import (
	"github.com/spf13/cobra"

	"arvika.pulcy.com/pulcy/pulcy/release"
)

const (
	defaultDockerRegistry = "arvika-ssh:5000"
)

var (
	releaseFlags = &release.Flags{}
	releaseCmd   = &cobra.Command{
		Use:   "release",
		Short: "Release a project",
		Long:  "Release a project. Update version, create tags, push etc",
		Run:   UsageFunc,
	}
)

func init() {
	releaseCmd.Flags().StringVarP(&releaseFlags.DockerRegistry, "registry", "r", defaultDockerRegistry, "Specify docker registry")
	mainCmd.AddCommand(releaseCmd)
	releaseCmd.AddCommand(&cobra.Command{
		Use:   "major",
		Short: "Create a major update",
		Run:   runRelease,
	})
	releaseCmd.AddCommand(&cobra.Command{
		Use:   "minor",
		Short: "Create a minor update",
		Run:   runRelease,
	})
	releaseCmd.AddCommand(&cobra.Command{
		Use:   "patch",
		Short: "Create a patch",
		Run:   runRelease,
	})
}

func runRelease(cmd *cobra.Command, args []string) {
	switch len(args) {
	case 0:
		releaseFlags.ReleaseType = cmd.Name()
		if err := release.Release(log, releaseFlags); err != nil {
			Quitf("Release failed: %v\n", err)
		} else {
			Infof("Release completed\n")
		}
	default:
		CommandError(cmd, "Too many arguments\n")
	}
}
