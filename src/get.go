package main

import (
	"github.com/spf13/cobra"

	"arvika.pulcy.com/developers/devtool/get"
)

var (
	getFlags = &get.Flags{}
	getCmd   = &cobra.Command{
		Use:   "get",
		Short: "Clone a repo into a folder",
		Long:  "Clone a repo into a folder, checking it out to a specific version",
		Run:   runGet,
	}
)

func init() {
	getCmd.Flags().StringVarP(&getFlags.Version, "version", "b", "", "Specify checkout version")
	mainCmd.AddCommand(getCmd)
}

func runGet(cmd *cobra.Command, args []string) {
	switch len(args) {
	case 2:
		getFlags.RepoUrl = args[0]
		getFlags.Folder = args[1]
		if err := get.Get(log, getFlags); err != nil {
			Quitf("Get failed: %v\n", err)
		}
	default:
		CommandError(cmd, "Expected <repo-url> <folder> arguments\n")
	}
}
