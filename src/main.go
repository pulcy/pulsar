package main

import (
	"fmt"
	"os"

	logPkg "github.com/op/go-logging"
	"github.com/spf13/cobra"
)

var (
	projectVersion = "dev"
	projectName    = "subliminl"
	projectBuild   = "dev"
	log            = logPkg.MustGetLogger(projectName)

	mainCmd = &cobra.Command{
		Use:   "subliminl",
		Short: "Subliminl is a helper for development environments",
		Run: func(cmd *cobra.Command, args []string) {
			// Do Stuff Here
		},
	}
)

func main() {
	mainCmd.Execute()
}

func Printf(message string, args ...interface{}) {
	fmt.Printf(message, args...)
}

func Quitf(message string, args ...interface{}) {
	Printf(message, args...)
	os.Exit(1)
}

// Print if quiet flag has not been set
func Infof(message string, args ...interface{}) {
	fmt.Printf(message, args...)
}
