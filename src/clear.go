package main

import (
	"github.com/spf13/cobra"

	"git.pulcy.com/pulcy/pulcy/cache"
)

var (
	clearCmd = &cobra.Command{
		Use: "clear",
		Run: UsageFunc,
	}
	clearCacheCmd = &cobra.Command{
		Use:   "cache",
		Short: "Clear a cache folder",
		Long:  "Clear a cache folder",
		Run:   runClearCache,
	}
)

func init() {
	clearCmd.AddCommand(clearCacheCmd)
	mainCmd.AddCommand(clearCmd)
}

func runClearCache(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		if err := cache.ClearAll(); err != nil {
			Quitf("Clear cache failed: %v\n", err)
		}
	} else {
		for _, key := range args {
			if err := cache.Clear(key); err != nil {
				Quitf("Clear cache failed: %v\n", err)
			}
		}
	}
}
