// Copyright (c) 2016 Pulcy.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/pulcy/pulsar/git"
	"github.com/pulcy/pulsar/settings"
)

var (
	projectCmd = &cobra.Command{
		Use:   "project",
		Short: "Project helpers",
		Run:   UsageFunc,
	}
	projectCommitCmd = &cobra.Command{
		Use:   "commit",
		Short: "Output project git commit",
		Run:   runProjectCommit,
	}
	projectVersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Output project version",
		Run:   runProjectVersion,
	}
)

func init() {
	mainCmd.AddCommand(projectCmd)
	projectCmd.AddCommand(projectCommitCmd)
	projectCmd.AddCommand(projectVersionCmd)
}

func runProjectCommit(cmd *cobra.Command, args []string) {
	commit, err := git.GetLatestLocalCommit(nil, ".", "", true)
	if err != nil {
		Quitf("%s\n", err)
	}
	fmt.Println(commit)
}

func runProjectVersion(cmd *cobra.Command, args []string) {
	version, err := settings.ReadVersion()
	if err != nil {
		Quitf("%s\n", err)
	}
	fmt.Println(version)
}
