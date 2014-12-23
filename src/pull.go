package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"arvika.subliminl.com/developers/subliminl/docker"
)

var (
	pullDockerRegistry string
	pullCmd            = &cobra.Command{
		Use:   "pull",
		Short: "Pull an image from the arvika-ssh registry",
		Long:  "Pull an image from the arvika-ssh registry",
		Run:   runPull,
	}
)

func init() {
	pullCmd.Flags().StringVarP(&pullDockerRegistry, "registry", "r", defaultDockerRegistry, "Specify docker registry")
	mainCmd.AddCommand(pullCmd)
}

func runPull(cmd *cobra.Command, args []string) {
	switch len(args) {
	case 0:
		pullUsageError("Too few arguments\n")
	case 1:
		err := docker.Pull(log, args[0], pushDockerRegistry)
		if err != nil {
			Quitf("%s\n", err)
		}
	default:
		pullUsageError("Too many arguments\n", args[0])
	}
}

func pullUsageError(prefix string, args ...interface{}) {
	prefix = fmt.Sprintf(prefix, args...)
	Quitf("%sUsage: %s pull image\n", prefix, projectName)
}
