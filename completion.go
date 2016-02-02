package main

import (
	"github.com/spf13/cobra"
)

var (
	completionCmd = &cobra.Command{
		Use:   "command-list",
		Short: "Gets all commands",
		Run:   runCompletion,
	}
)

func init() {
	mainCmd.AddCommand(completionCmd)
}

func runCompletion(cmd *cobra.Command, args []string) {
	for _, c := range mainCmd.Commands() {
		Printf("%s ", c.Name())
	}
	Printf("\n")
}
