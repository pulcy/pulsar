package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"arvika.subliminl.com/developers/subliminl/docker"
)

var (
	pushDockerRegistry string
	pushCmd            = &cobra.Command{
		Use:   "push",
		Short: "Push an image to the arvika-ssh registry",
		Long:  "Push an image to the arvika-ssh registry",
		Run:   runPush,
	}
)

func init() {
	pushCmd.Flags().StringVarP(&pushDockerRegistry, "registry", "r", defaultDockerRegistry, "Specify docker registry")
	mainCmd.AddCommand(pushCmd)
}

func runPush(cmd *cobra.Command, args []string) {
	switch len(args) {
	case 0:
		pushUsageError("Too few arguments\n")
	case 1:
		err := docker.Push(log, args[0], pushDockerRegistry)
		if err != nil {
			Quitf("%s\n", err)
		}
	default:
		pushUsageError("Too many arguments\n", args[0])
	}
}

func pushUsageError(prefix string, args ...interface{}) {
	prefix = fmt.Sprintf(prefix, args...)
	Quitf("%sUsage: %s push image\n", prefix, projectName)
}
