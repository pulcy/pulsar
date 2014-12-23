package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"arvika.subliminl.com/developers/subliminl/release"
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
		Run:   runRelease,
	}
)

func init() {
	releaseCmd.Flags().StringVarP(&releaseFlags.DockerRegistry, "registry", "r", defaultDockerRegistry, "Specify docker registry")
	mainCmd.AddCommand(releaseCmd)
}

func runRelease(cmd *cobra.Command, args []string) {
	switch len(args) {
	case 0:
		runUsageError("Too few arguments.\n")
	case 1:
		switch args[0] {
		case "major":
			fallthrough
		case "minor":
			fallthrough
		case "patch":
			releaseFlags.ReleaseType = args[0]
			if err := release.Release(log, releaseFlags); err != nil {
				Quitf("Release failed: %v\n", err)
			} else {
				Infof("Release completed\n")
			}
		default:
			runUsageError("Invalid release type %s.\n", args[0])
		}
	default:
		runUsageError("Too many arguments.\n")
	}
}

func runUsageError(prefix string, args ...interface{}) {
	prefix = fmt.Sprintf(prefix, args...)
	Quitf("%sUsage: %s release major|minor|patch\n", prefix, projectName)
}
