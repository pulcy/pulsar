package main

import (
	"strings"

	"github.com/spf13/cobra"
)

const (
	hdr = "   ___       __            " + "\n" +
		"  / _ \\__ __/ /__ ___ _____" + "\n" +
		" / ___/ // / (_-</ _ `/ __/" + "\n" +
		"/_/   \\_,_/_/___/\\_,_/_/   " + "\n"
)

var (
	cmdVersion = &cobra.Command{
		Use: "version",
		Run: showVersion,
	}
)

func init() {
	mainCmd.AddCommand(cmdVersion)
}

func showVersion(cmd *cobra.Command, args []string) {
	for _, line := range strings.Split(hdr, "\n") {
		log.Info(line)
	}
	log.Info("%s %s, build %s\n", mainCmd.Use, projectVersion, projectBuild)
}
