package main

import (
	"os"
	"sync"

	"github.com/spf13/cobra"

	goget "git.pulcy.com/pulcy/pulcy/goget"
)

var (
	goCmd = &cobra.Command{
		Use:   "go",
		Short: "Execute `go get` with cache support",
		Run:   UsageFunc,
	}
	gogetCmd = &cobra.Command{
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
	if len(args) == 0 {
		CommandError(cmd, "Expected <package> argument\n")
	} else {
		wg := sync.WaitGroup{}
		errors := make(chan error, len(args))
		for _, pkg := range args {
			wg.Add(1)
			go func(pkg string) {
				defer wg.Done()
				gogetFlags := &goget.Flags{Package: pkg}
				if err := goget.Get(log, gogetFlags); err != nil {
					errors <- err
				}
			}(pkg)
		}
		wg.Wait()
		close(errors)
		failed := false
		for err := range errors {
			Printf("Go get failed: %v\n", err)
			failed = true
		}
		if failed {
			os.Exit(1)
		}
	}
}
