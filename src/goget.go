package main

import (
	"github.com/spf13/cobra"

	goget "arvika.pulcy.com/developers/devtool/goget"
)

var (
	goCmd = &cobra.Command{
		Use:   "go",
		Short: "Execute `go get` with cache support",
		Run:   UsageFunc,
	}
	gogetFlags = &goget.Flags{}
	gogetCmd   = &cobra.Command{
		Use:   "get",
		Short: "Execute `go get` with cache support",
		Run:   runGoGet,
	}
)

func init() {
	mainCmd.AddCommand(goCmd)
	goCmd.AddCommand(gogetCmd)
}

func runGoGet(cmd *cobra.Command, args []string) {
	switch len(args) {
	case 1:
		gogetFlags.Package = args[0]
		if err := goget.Get(log, gogetFlags); err != nil {
			Quitf("Go get failed: %v\n", err)
		}
	default:
		CommandError(cmd, "Expected <package> argument\n")
	}
}
