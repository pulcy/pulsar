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
	}
)

func init() {
	mainCmd.Run = runMain
}

func main() {
	mainCmd.Execute()
}

func runMain(cmd *cobra.Command, args []string) {
	mainCmd.Usage()
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

func CommandError(c *cobra.Command, prefix string, args ...interface{}) {
	prefix = fmt.Sprintf(prefix, args...)
	Quitf("%sUsage: %s\n", prefix, c.CommandPath())
}
