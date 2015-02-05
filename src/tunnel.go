package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"arvika.subliminl.com/developers/devtool/tunnel"
)

var (
	tunnelCmd = &cobra.Command{
		Use:   "tunnel [open|close]",
		Short: "Open/close an SSH tunnel to arvika-ssh",
		Long:  "Open/close an SSH tunnel that listens on arvika-ssh:<various ports>",
		Run:   runTunnel,
	}
)

func init() {
	mainCmd.AddCommand(tunnelCmd)
}

func runTunnel(cmd *cobra.Command, args []string) {
	switch len(args) {
	case 0:
		openTunnel()
	case 1:
		switch args[0] {
		case "open":
			openTunnel()
		case "close":
			closeTunnel()
		default:
			tunnelUsageError("Unknown argument %s\n", args[0])
		}
	default:
		tunnelUsageError("Too many arguments\n", args[0])
	}
}

func openTunnel() {
	err := tunnel.OpenTunnel(log)
	if err != nil {
		Quitf("%s\n", err)
	}
}

func closeTunnel() {
	err := tunnel.CloseTunnel(log)
	if err != nil {
		Quitf("%s\n", err)
	}
}

func tunnelUsageError(prefix string, args ...interface{}) {
	prefix = fmt.Sprintf(prefix, args...)
	Quitf("%sUsage: %s tunnel [open|close]\n", prefix, projectName)
}
